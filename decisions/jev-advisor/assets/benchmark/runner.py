"""Execute frozen paired trials through explicit JSON adapters."""

import argparse
import json
import os
import random
import shutil
import subprocess
import sys
import time
from pathlib import Path

from evidence import POLICY, context_check, digest, encoded, freeze, inventory_changes, read, request, require_hosted_opt_in, route, save, validate
from codex_adapter import captured_native, native_identity
from integrity import capture_trial, file_inventory, start_ledger
from execution import execute, now
from feedback import Feedback
from verification import copy_workspace, ensure_verified, verify


def adapter_identities(adapters):
    identities = {}
    for role, argv in adapters.items():
        executable = shutil.which(argv[0])
        if not executable:
            raise ValueError('adapter executable unavailable: ' + role)
        paths = [Path(executable).resolve()]
        paths.extend(Path(arg).resolve() for arg in argv[1:] if Path(arg).is_file())
        identities[role] = {str(path): digest(path.read_bytes()) for path in paths}
    return identities


def implementation_identity():
    return {p.name: digest(p.read_bytes()) for p in sorted(Path(__file__).parent.glob('*.py'))}


class ProviderInputError(ValueError):
    pass


def schedule(corpus, repetitions, seed):
    if type(repetitions) is not int or repetitions < 2:
        raise ValueError('paired benchmark requires at least two repetitions')
    rng = random.Random(seed)
    pairs = [(task['id'], repeat) for task in corpus for repeat in range(repetitions)]
    rng.shuffle(pairs)
    rows = []
    for index, (task, repeat) in enumerate(pairs):
        conditions = ['baseline', 'treatment']
        if index % 2:
            conditions.reverse()
        rows.extend({'task': task, 'repeat': repeat, 'condition': condition} for condition in conditions)
    return rows


def invoke(command, payload, cwd, output, timeout, frozen=None, service=None):
    output.mkdir(parents=True, exist_ok=False)
    body = encoded(payload)
    (output / 'request.json').write_bytes(body)
    started, tick = now(), time.monotonic_ns()
    result = {'started_at': started, 'request_sha256': digest(body), 'command': command,
              'status': 'failed', 'response': None}
    environment = dict(os.environ, JEV_BENCHMARK_CALL_DIR=str(output.resolve()))
    if frozen is not None:
        environment['JEV_BENCHMARK_FROZEN_ROOT'] = str(frozen)
    if payload.get('role') in ('task', 'fallback'):
        environment.pop('TYPESAFE_API_KEY', None)
    result.update(execute(command, cwd, output, timeout, environment, body, service=service))
    try:
        if (output / 'stdout').stat().st_size:
            result['response'] = read(output / 'stdout')
        if result['status'] == 'completed' and result['response'] is None:
            result['status'] = 'failed'
    except (OSError, ValueError) as error:
        result['error'] = str(error)
        if result['status'] == 'completed':
            result['status'] = 'failed'
    if (result['response'] is None and payload.get('role') in ('task', 'fallback')
            and (output / 'events.jsonl').exists() and isinstance(payload.get('execution', {}).get('native_line'), dict)):
        result['response'] = captured_native(output, payload['execution'], result.get('exit_code'))
    result.update(ended_at=now(), elapsed_ns=time.monotonic_ns() - tick,
                  clock='monotonic_ns', timing_scope='adapter launch through output capture and decode')
    save(output / 'result.json', result)
    return result


def native_payload(packet, selected, source, role):
    bodies = {entry['name']: (source / entry['path']).read_text()
              for entry in packet['catalog'] if entry['name'] in selected}
    return {'role': role, 'evidence': packet, 'initial_skills': bodies,
            'instruction': 'Complete the task and self-review under the supplied rules. '
            'The full skill catalog and bodies remain available on disk. '
            'Load any additional required skill you discover. Record that rescue in your result. '
            'Rescues may name catalog skills or frozen additional_guidance; initial selection uses the catalog.'}


def selected_fallback(result, packet):
    if result['status'] != 'completed':
        raise ValueError('fallback failed')
    response = result['response']
    selected = response['selected']
    available = {entry['name'] for entry in packet['catalog']}
    if (not isinstance(selected, list) or any(not isinstance(s, str) for s in selected)
            or not set(selected).issubset(available) or len(selected) != len(set(selected))
            or not set(packet['task']['protected']).issubset(selected)
            or response.get('abstain') is not False):
        raise ValueError('invalid fallback selection')
    return selected


def changed_files(source, workspace):
    before, after = file_inventory(source), file_inventory(workspace)
    return inventory_changes(before, after)


def assurance_status(response, packet, selected):
    for key in ('verification', 'self_review'):
        claims = response.get(key)
        if not isinstance(claims, list) or not claims or any(
                not isinstance(claim, str) or not claim.strip() for claim in claims):
            return 'missing_assurance_evidence'
    rescues = response.get('rescued_skills')
    available = {entry['name'] for entry in packet['catalog']}.union(packet['additional_guidance'])
    if (not isinstance(rescues, list) or any(not isinstance(s, str) for s in rescues)
            or len(rescues) != len(set(rescues)) or not set(rescues).issubset(available)
            or set(rescues).intersection(selected)):
        return 'invalid_rescue_evidence'
    return 'completed'


def trial(root, output, row, config):
    started, tick = now(), time.monotonic_ns()
    packet = read(root / 'packets' / (row['task'] + '.json'))
    source = root / 'source'
    work = output / 'workspace'
    copy_workspace(source, work)
    protected = packet['task']['protected']
    execution = {key: config[key] for key in ('native_line', 'harness', 'service_tier')}
    if 'native_executable' in config:
        execution['native_executable'] = config['native_executable']
    selected = protected
    record = dict(row, status='incomplete', quality=None, attempts=[], started_at=started)
    try:
        if row['condition'] == 'treatment':
            result = invoke(config['adapters']['jev'], request(packet, config['jev_model']),
                            output, output / 'jev', config['timeout_seconds'], frozen=root)
            record['attempts'].append('jev')
            if result['status'] == 'interrupted':
                record['status'] = 'interrupted'
                return record
            if isinstance(result['response'], dict) and result['response'].get('http_status') in (401, 422):
                raise ProviderInputError('provider rejected credentials or input; stop the comparison without fallback')
            selection = route(packet, result['response'], config['jev_model'])
            if result['status'] != 'completed':
                selection = route(packet, None, config['jev_model'])
            save(output / 'selection.json', selection)
            selected = selection['selected']
            if selection['route'] == 'fallback':
                fallback = invoke(config['adapters']['fallback'],
                                  {'role': 'fallback', 'evidence': packet,
                                   'execution': execution,
                                   'instruction': 'Return selected skill names and abstain=false. '
                                   'Preserve protected skills; resolve required applicability. '
                                   'Optional benefit does not create a requirement.'},
                                  output, output / 'fallback', config['timeout_seconds'])
                record['attempts'].append('fallback')
                if fallback['status'] == 'interrupted':
                    record['status'] = 'interrupted'
                    return record
                selected = selected_fallback(fallback, packet)
        payload = native_payload(packet, selected, source, 'task')
        payload['execution'] = execution
        verification = output / 'verification'
        verification.mkdir()
        mailbox = output / 'feedback'
        feedback = Feedback(work, source, packet['task']['writes'], mailbox, verification)
        payload['verification_command'] = [sys.executable, '-B', str(Path(__file__).with_name('feedback.py').resolve()),
                                           '--source', str(work), '--mailbox', str(mailbox),
                                           '--timeout', str(config['timeout_seconds'])]
        payload['verification_output'] = str(mailbox)
        payload['selection_instruction'] = (
            'Select the initial skills from the catalog before the task.' if row['condition'] == 'baseline'
            else 'Start with the supplied skill selection; inspect the catalog if the action needs more skills.')
        native = invoke(config['adapters']['task'], payload, work, output / 'task',
                        config['timeout_seconds'], service=feedback)
        record['attempts'].append('task')
        record['status'] = native['status']
        response = native['response'] if isinstance(native['response'], dict) else {}
        assurance = assurance_status(response, packet, selected)
        record['rescued_skills'] = response.get('rescued_skills') if assurance == 'completed' else None
        record['verification_claims'] = response.get('verification')
        record['self_review_claims'] = response.get('self_review')
        if native['status'] == 'completed':
            record['status'] = assurance
        changed = changed_files(source, work)
        record['changed'] = changed
        record['outside_write_fence'] = sorted(set(changed) - set(packet['task']['writes']))
        if record['outside_write_fence'] and record['status'] != 'interrupted':
            record['status'] = 'invalid_write_fence'
        elif record['status'] == 'completed':
            record['verification_result'] = ensure_verified(work, verification, force=True)
            if record['verification_result']['status'] != 'completed':
                record['status'] = ('interrupted' if record['verification_result']['status'] == 'interrupted'
                                    else 'verification_failed')
    except ProviderInputError as error:
        record.update(status='invalid_provider_input', error=str(error))
    except (ValueError, KeyError, TypeError, OSError) as error:
        record.update(status='abstained', error=str(error))
    finally:
        record.update(ended_at=now(), elapsed_ns=time.monotonic_ns() - tick,
                      clock='monotonic_ns', timing_scope='workspace setup through selection, task, and artifact capture')
        save(output / 'trial.json', record)
    return record


def run(root, output, config, authorization, offline=False):
    root, output = Path(root).resolve(), Path(output).resolve()
    manifest = validate(root)
    if config.get('input_freeze_sha256') != digest(encoded(manifest)):
        raise ValueError('config does not name this frozen evidence')
    if output.exists():
        raise ValueError('run destination already exists; no hidden retry')
    if not offline and (os.environ.get('BENCH_OFFLINE') or
                        authorization != digest(encoded(config))):
        raise ValueError('live run requires authorization for the exact config and online mode')
    if not offline:
        require_hosted_opt_in(root)
    if set(config['adapters']) != {'jev', 'fallback', 'task'}:
        raise ValueError('three explicit adapters are required')
    for command in config['adapters'].values():
        if not isinstance(command, list) or not command or any(not isinstance(s, str) or not s for s in command):
            raise ValueError('adapter must be a nonempty argv array')
    if config.get('adapter_sha256') != adapter_identities(config['adapters']):
        raise ValueError('adapter identity differs from reviewed config')
    if config.get('implementation_sha256') != implementation_identity():
        raise ValueError('benchmark implementation differs from reviewed config')
    if 'native_executable' in config and config['native_executable'] != native_identity(config['native_executable']['path']):
        raise ValueError('native executable identity differs from reviewed config')
    if not 0 < config['timeout_seconds'] <= 3600:
        raise ValueError('invalid adapter timeout')
    rows = schedule(read(root / 'corpus.json'), config['repetitions'], config['seed'])
    if config['max_adapter_calls'] < len(rows) * 2:
        raise ValueError('call budget cannot cover the worst-case schedule')
    if not config.get('budget_reference') or not config.get('harness') or not config.get('native_line'):
        raise ValueError('missing budget or native execution provenance')
    output.mkdir()
    save(output / 'plan.json', {'config': config, 'config_sha256': digest(encoded(config)),
                             'freeze_sha256': digest(encoded(manifest)), 'schedule': rows,
                             'offline': offline, 'policy': POLICY,
                             'runner_sha256': digest(Path(__file__).read_bytes()),
                             'evidence_sha256': digest(Path(__file__).with_name('evidence.py').read_bytes())})
    preflight = verify(root / 'source', output / 'preflight')
    ledger = start_ledger(output)
    if preflight['status'] != 'completed':
        return {'trials': 0, 'stopped': 'invalid_baseline', 'quality': 'ungraded',
                'ledger_sha256': digest(encoded(ledger))}
    for index, row in enumerate(rows):
        validate(root)
        if config['adapter_sha256'] != adapter_identities(config['adapters']):
            raise ValueError('adapter identity changed during run')
        if config['implementation_sha256'] != implementation_identity():
            raise ValueError('benchmark implementation changed during run')
        if 'native_executable' in config and config['native_executable'] != native_identity(config['native_executable']['path']):
            raise ValueError('native executable changed during run')
        destination = output / ('trial-%03d' % index)
        destination.mkdir()
        result = trial(root, destination, row, config)
        capture_trial(output, destination.name, ledger)
        print(json.dumps(dict(row, completed=index + 1)), flush=True)
        if result['status'] in ('invalid_provider_input', 'interrupted'):
            return {'trials': index + 1, 'stopped': result['status'], 'quality': 'ungraded',
                    'ledger_sha256': digest(encoded(ledger))}
    return {'trials': len(rows), 'quality': 'ungraded', 'cost': 'use bench assessment with adapter native evidence',
            'ledger_sha256': digest(encoded(ledger))}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    subs = parser.add_subparsers(dest='action', required=True)
    prepare = subs.add_parser('freeze')
    prepare.add_argument('--repo', required=True)
    prepare.add_argument('--corpus', default=str(Path(__file__).with_name('corpus.json')))
    prepare.add_argument('--out', required=True)
    check = subs.add_parser('validate')
    check.add_argument('root')
    plan = subs.add_parser('plan')
    plan.add_argument('root')
    plan.add_argument('--out', required=True)
    plan.add_argument('--model', required=True)
    plan.add_argument('--effort', required=True)
    plan.add_argument('--harness', required=True)
    plan.add_argument('--budget-reference', required=True)
    plan.add_argument('--repetitions', type=int, default=3)
    plan.add_argument('--seed', type=int, default=20260921)
    launch = subs.add_parser('run')
    launch.add_argument('root')
    launch.add_argument('--out', required=True)
    launch.add_argument('--config', required=True)
    launch.add_argument('--authorize-config-sha256')
    args = parser.parse_args()
    if args.action == 'freeze':
        result = freeze(args.repo, args.corpus, args.out)
    elif args.action == 'validate':
        manifest = validate(args.root)
        result = {'valid': True, 'revision': manifest['revision']}
    elif args.action == 'plan':
        manifest = validate(args.root)
        rows = schedule(read(Path(args.root) / 'corpus.json'), args.repetitions, args.seed)
        here = Path(__file__).resolve().parent
        native = [sys.executable, str(here / 'codex_adapter.py')]
        config = {'input_freeze_sha256': digest(encoded(manifest)),
                  'adapters': {'jev': [sys.executable, str(here / 'jev_http.py')],
                               'fallback': native, 'task': native},
                  'native_line': {'model': args.model, 'effort': args.effort},
                  'harness': args.harness, 'service_tier': 'default', 'jev_model': 'jev-1.13.0',
                  'repetitions': args.repetitions, 'seed': args.seed, 'timeout_seconds': 600,
                  'max_adapter_calls': len(rows) * 2, 'budget_reference': args.budget_reference,
                  'rates': {'jev': None, 'native': None}}
        config['native_executable'] = native_identity()
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        config['implementation_sha256'] = implementation_identity()
        config['context_checks'] = {t['id']: context_check(read(Path(args.root) / 'packets' / (t['id'] + '.json')),
                                                         config['jev_model']) for t in read(Path(args.root) / 'corpus.json')}
        save(args.out, config)
        result = {'config_sha256': digest(encoded(config)), 'trials': len(rows),
                  'max_adapter_calls': config['max_adapter_calls'], 'paid_calls': 0}
    else:
        result = run(args.root, args.out, read(args.config), args.authorize_config_sha256)
    print(json.dumps(result))


if __name__ == '__main__':
    main()
