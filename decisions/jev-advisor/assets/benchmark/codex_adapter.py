"""Run a pinned Codex task or fallback and retain native usage events."""

import json
import os
import shutil
import subprocess
import sys
from pathlib import Path

from evidence import digest, encoded, read, save


def native_identity(path=None):
    executable = Path(path or shutil.which('codex') or '').resolve()
    if not executable.is_file() or not os.access(executable, os.X_OK):
        raise ValueError('native executable unavailable')
    return {'path': str(executable), 'sha256': digest(executable.read_bytes())}


def usage_events(events, session, reference):
    result = []
    for sequence, event in enumerate(events):
        if event.get('type') != 'turn.completed':
            continue
        native = event.get('usage', {})
        usage = {'input_total': native.get('input_tokens'),
                 'input_cached': native.get('cached_input_tokens'),
                 'output': native.get('output_tokens'), 'total_semantics': 'inclusive'}
        result.append({'event_id': 'turn-%d' % sequence, 'session_id': session,
                       'epoch': 0, 'sequence': sequence, 'mode': 'delta', 'counter': 'tokens',
                       'reference': reference, 'usage': usage})
    return result or [{'event_id': 'missing', 'session_id': session, 'epoch': 0, 'sequence': 0,
                       'mode': 'delta', 'counter': 'tokens', 'reference': reference,
                       'usage': {'unknown': ['No complete native usage event.']}}]


def captured_native(output, execution, exit_code):
    events = []
    path = output / 'events.jsonl'
    for line in path.read_text().splitlines() if path.exists() else []:
        try:
            event = json.loads(line)
            if isinstance(event, dict):
                events.append(event)
        except ValueError:
            pass
    session = next((e['thread_id'] for e in events if e.get('type') == 'thread.started'), 'unknown')
    try:
        result = read(output / 'final.json')
        if not isinstance(result, dict):
            result = {}
    except (OSError, ValueError):
        result = {}
    reference = {'producer': execution['harness'], 'native': str(path)}
    result.update(model=execution['native_line']['model'], effort=execution['native_line']['effort'],
                  session_id=session, native_exit_code=exit_code,
                  usage_scope='Observed parent turn events; nested usage inclusion is unverified.',
                  assessment_usage=usage_events(events, session, reference))
    if exit_code != 0:
        result['assessment_usage'].append({
            'event_id': 'incomplete-call', 'session_id': session, 'epoch': 0,
            'sequence': len(events), 'mode': 'delta', 'counter': 'tokens', 'reference': reference,
            'usage': {'unknown': ['Native call did not complete; uncaptured usage may exist.']}})
    return result


def main():
    if os.environ.get('BENCH_OFFLINE'):
        raise SystemExit('BENCH_OFFLINE suppresses native model calls')
    payload = json.load(sys.stdin)
    execution = payload['execution']
    output = Path(os.environ['JEV_BENCHMARK_CALL_DIR'])
    identity = execution['native_executable']
    if native_identity(identity['path']) != identity:
        raise SystemExit('native executable identity differs from frozen config')
    executable = identity['path']
    version = subprocess.check_output([executable, '--version'], text=True).strip()
    if version != execution['harness']:
        raise SystemExit('Codex version differs from the frozen config')
    line = execution['native_line']
    command = [executable, 'exec', '--ephemeral', '--ignore-user-config', '--skip-git-repo-check',
               '--sandbox', 'workspace-write', '--cd', str(Path.cwd()),
               '--model', line['model'], '-c', 'model_reasoning_effort=' + json.dumps(line['effort']),
               '-c', 'service_tier=' + json.dumps(execution['service_tier']),
               '--json', '--color', 'never', '--output-last-message', str(output / 'final.json')]
    if payload.get('verification_output'):
        command.extend(['--add-dir', payload['verification_output']])
    fields = {'selected': {'type': 'array', 'items': {'type': 'string'}},
              'abstain': {'type': 'boolean'}} if payload['role'] == 'fallback' else {
                  'summary': {'type': 'string'},
                  'rescued_skills': {'type': 'array', 'items': {'type': 'string'}},
                  'self_review': {'type': 'array', 'minItems': 1, 'items': {'type': 'string', 'minLength': 1}},
                  'verification': {'type': 'array', 'minItems': 1, 'items': {'type': 'string', 'minLength': 1}}}
    schema = {'type': 'object', 'properties': fields, 'required': list(fields), 'additionalProperties': False}
    save(output / 'schema.json', schema)
    command.extend(['--output-schema', str(output / 'schema.json'), '-'])
    save(output / 'native-command.json', {'argv': command, 'version': version})
    with (output / 'events.jsonl').open('wb') as stdout:
        child = subprocess.run(command, input=encoded(payload), stdout=stdout, stderr=sys.stderr.buffer)
    result = captured_native(output, execution, child.returncode)
    sys.stdout.buffer.write(encoded(result))
    raise SystemExit(child.returncode)


if __name__ == '__main__':
    main()
