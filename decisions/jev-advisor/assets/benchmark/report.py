"""Prepare blinded quality evidence and export native assessment records."""

import argparse
from pathlib import Path

from evidence import digest, encoded, read, save, validate
from integrity import verified_run


def reviews(frozen, runs, output):
    frozen, runs, output = Path(frozen), Path(runs), Path(output)
    manifest = validate(frozen)
    plan, trials = verified_run(runs)
    if plan['freeze_sha256'] != digest(encoded(manifest)):
        raise ValueError('run and review source differ')
    if output.exists():
        raise ValueError('review destination already exists')
    tasks = {t['id']: t for t in read(frozen / 'corpus.json')}
    index = {}
    packets = {}
    for trial_path, trial in trials:
        if 'task' not in trial['attempts']:
            continue
        task = tasks[trial['task']]
        artifacts = {}
        for name in task['writes']:
            path = trial_path.parent / 'workspace' / name
            artifacts[name] = path.read_text() if path.is_file() and not path.is_symlink() else None
        evidence = {'task': task['brief'], 'criteria': task['criteria'],
                    'artifacts': artifacts, 'source_revision': manifest['revision'],
                    'validity': {'completed': trial['status'] == 'completed', 'changed': trial.get('changed'),
                                 'outside_write_fence': trial.get('outside_write_fence')},
                    'assurance_claim_counts': {key: sum(isinstance(claim, str) and bool(claim.strip())
                                                        for claim in (trial.get(key) or [])
                                                        ) if isinstance(trial.get(key), list) else 0
                                               for key in ('self_review_claims', 'verification_claims')}}
        identity = digest(encoded(evidence))
        key = 'artifact-' + identity
        if key not in index:
            packets[key] = evidence
        index.setdefault(key, []).append(str(trial_path.relative_to(runs)))
    for key, evidence in packets.items():
        save(output / 'packets' / (key + '.json'), evidence)
    save(output / 'private-index.json', index)
    return {'unique_review_packets': len(index), 'quality': 'awaiting independent criterion verdicts'}


def export(frozen, runs, output, repo_key):
    frozen, runs, output = Path(frozen), Path(runs), Path(output)
    manifest = validate(frozen)
    plan, trials = verified_run(runs)
    if plan['freeze_sha256'] != digest(encoded(manifest)):
        raise ValueError('run and export source differ')
    if output.exists():
        raise ValueError('export destination already exists')
    tasks = {t['id']: t for t in read(frozen / 'corpus.json')}
    records = []
    for trial_path, trial in trials:
        identity = digest(encoded(plan))[:16] + '-' + trial_path.parent.name
        reference = {'producer': 'jev-benchmark', 'native': str(trial_path.resolve())}
        record = {'version': 1, 'run_id': identity, 'repo_key': repo_key,
                  'source': manifest['revision'], 'condition': trial['condition'],
                  'task_id': trial['task'], 'holdout': tasks[trial['task']]['split'] == 'holdout',
                  'state': 'incomplete', 'attempts': [], 'evidence': [reference], 'quality': {}}
        for role in trial['attempts']:
            path = trial_path.parent / role / 'result.json'
            result = read(path)
            response = result['response'] if isinstance(result['response'], dict) else {}
            ref = {'producer': 'jev-benchmark-adapter', 'native': str(path.resolve())}
            session = response.get('session_id', identity + '-' + role)
            usage = response.get('assessment_usage')
            if role == 'jev' and isinstance(response.get('usage'), dict):
                usage = [{'event_id': 'response', 'session_id': session, 'epoch': 0,
                          'sequence': 0, 'mode': 'delta', 'counter': 'tokens', 'reference': ref,
                          'usage': {'input_uncached': response['usage'].get('input_tokens'),
                                    'input_cached': 0, 'output': response['usage'].get('output_tokens')}}]
            if usage is None:
                usage = [{'event_id': 'unknown', 'session_id': session, 'epoch': 0,
                          'sequence': 0, 'mode': 'delta', 'counter': 'tokens',
                          'reference': ref, 'usage': {'unknown': ['No normalized native usage supplied.']}}]
            attempt = {'attempt_id': identity + '-' + role, 'chunk_id': trial['task'],
                       'role': 'implementation' if role == 'task' else 'orchestration',
                       'session_id': session, 'model': response.get('model', 'unknown'),
                       'effort': response.get('effort', 'unknown'),
                       'state': ('succeeded' if result['status'] == 'completed' else
                                 'cancelled' if result['status'] == 'interrupted' else 'failed'),
                       'usage': usage, 'cost': {'estimated': plan['config'].get('rates', {}).get(
                           'jev' if role == 'jev' else 'native'),
                           'actual': response.get('actual_charges', [])},
                       'measures': {'elapsed_ns': {'value': result['elapsed_ns'], 'reference': ref}},
                       'evidence': [ref]}
            if role != 'jev':
                attempt['cost']['other'] = [{'kind': 'nested-native-usage-coverage-unverified',
                                             'currency': 'USD', 'reference': ref}]
            if role == 'task' and isinstance(trial.get('rescued_skills'), list):
                attempt['measures']['rescued_skills'] = {'value': len(trial['rescued_skills']), 'reference': reference}
            record['attempts'].append(attempt)
        save(output / (identity + '.json'), record)
        records.append(identity)
    return {'assessment_records': records,
            'note': 'Import with bench assessment record --input. Missing costs and quality stay unknown.'}


def summary(runs):
    runs = Path(runs)
    plan, trials = verified_run(runs)
    captured = {path: trial for path, trial in trials}
    rows = []
    for index, scheduled in enumerate(plan['schedule']):
        path = runs / ('trial-%03d' % index) / 'trial.json'
        if path not in captured:
            rows.append(dict(scheduled, status='missing', elapsed_ns=None))
            continue
        trial = captured[path]
        row = {k: trial[k] for k in ('task', 'condition', 'repeat', 'status', 'elapsed_ns')}
        row.update(rescued_skills=trial.get('rescued_skills'), self_review_claims=trial.get('self_review_claims'),
                   verification_claims=trial.get('verification_claims'))
        row['roles'] = {role: {k: read(path.parent / role / 'result.json')[k]
                              for k in ('status', 'elapsed_ns')} for role in trial['attempts']}
        rows.append(row)
    return {'trials': rows, 'quality': 'ungraded', 'cost_comparison': 'use native assessment records',
            'adoption_evidence': False,
            'limits': ['Completion is not a quality verdict.', 'Elapsed time uses monotonic nanoseconds.',
                       'Provider caches and model randomness remain uncontrolled.',
                       'The supplied corpus is exposed development evidence.',
                       'Native costs cover observed parent turns; nested usage inclusion remains unknown.']}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=['reviews', 'export', 'summary'])
    parser.add_argument('--runs', required=True)
    parser.add_argument('--frozen')
    parser.add_argument('--out')
    parser.add_argument('--repo-key')
    args = parser.parse_args()
    if args.action == 'summary':
        result = summary(args.runs)
    else:
        if not args.frozen or not args.out:
            parser.error('--frozen and --out are required')
        if args.action == 'reviews':
            result = reviews(args.frozen, args.runs, args.out)
        else:
            if not args.repo_key:
                parser.error('--repo-key is required')
            result = export(args.frozen, args.runs, args.out, args.repo_key)
    print(encoded(result).decode(), end='')


if __name__ == '__main__':
    main()
