"""Build and validate source-backed skill-selection evidence."""

import hashlib
import json
import math
import os
import stat
import re
import subprocess
from pathlib import Path, PurePosixPath


RULES = ('AGENTS.md', '.bench/BENCH.md',
         'projects/benchkit.md', 'DATA_HANDLING.md',
         '.agents/skills/bench-craft-spec/references/ste-prose.md')
LOOKUP_DOCUMENTS = ('.bench/BENCH-reference.md', 'CONTEXT.md')
CONTEXT_LIMITS = {'state_and_longest_question_bytes': 120000, 'request_bytes': 240000}
POLICY = {'required_negative_max': 0.2, 'required_positive_min': 0.8,
          'useful_positive_min': 0.8}
SCOPE = {
    'unit': 'Complete bounded task in an isolated replay of the repository.',
    'included': ['source inspection', 'implementation when requested',
                 'focused verification', 'self-review', 'in-session repairs'],
    'excluded': ['production ticket or spec authoring', 'production commit', 'production landing', 'publication'],
    'verification': 'The runner performs the mandatory final bench gate in an exact copy of your task files. '
                    'Use focused checks and self-review during the task. If you need earlier full-gate feedback, '
                    'use the supplied verification_command. Do not run bench gate directly in the authored '
                    'workspace. All gate outputs are captured separately and final completion requires a green gate.',
    'override': 'This replay permits the bounded named artifact without a production ticket, spec, or phase close. '
                'All other project rules apply. Self-review is in scope for both arms.',
}


def encoded(value):
    return (json.dumps(value, sort_keys=True, ensure_ascii=False, allow_nan=False,
                       indent=2) + '\n').encode()


def digest(data):
    return hashlib.sha256(data).hexdigest()


def save(path, value):
    path = Path(path)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open('xb') as stream:
        stream.write(encoded(value))


def read(path):
    return json.loads(Path(path).read_text())


def relative(name):
    path = PurePosixPath(name)
    if not name or path.is_absolute() or '..' in path.parts or str(path) != name:
        raise ValueError('unsafe relative path: ' + name)
    return name


def git(root, *args):
    return subprocess.check_output(['git', *args], cwd=root)


def file_inventory(root):
    inventory = {}
    for path in root.rglob('*'):
        if '.git' in path.relative_to(root).parts:
            continue
        if path.is_symlink():
            value = {'symlink': os.readlink(path)}
        elif path.is_file():
            value = {'sha256': digest(path.read_bytes()), 'mode': stat.S_IMODE(path.stat().st_mode)}
        elif not path.is_dir():
            value = {'special': path.lstat().st_mode}
        else:
            continue
        inventory[str(path.relative_to(root))] = value
    return inventory


def validate_links(source):
    for path in source.rglob('*'):
        if not path.is_symlink():
            continue
        try:
            target = os.readlink(path)
            resolved = path.resolve(strict=True)
            relative_target = resolved.relative_to(source.resolve())
            if Path(target).is_absolute() or '.git' in relative_target.parts:
                raise ValueError('absolute or Git target')
        except (OSError, RuntimeError, ValueError) as error:
            raise ValueError('unsafe snapshot link: ' + str(path.relative_to(source))) from error


def snapshot(root):
    revision = git(root, 'rev-parse', 'HEAD').decode().strip()
    files = {}
    for row in git(root, 'ls-tree', '-rz', '--full-tree', revision).split(b'\0'):
        if not row:
            continue
        metadata, name = row.split(b'\t', 1)
        mode, kind, oid = metadata.decode().split()
        name = relative(name.decode())
        if kind != 'blob' or mode not in ('100644', '100755', '120000'):
            raise ValueError('unsupported snapshot entry: ' + name)
        files[name] = {'mode': mode, 'oid': oid}
    return revision, files


def skill_entry(path, content):
    body = content(path)
    fields = dict(line.split(':', 1) for line in body.split('---', 2)[1].splitlines()
                  if ':' in line and not line.startswith(' '))
    entry = {key: fields[key].strip() for key in ('name', 'description')}
    if 'index' in fields:
        entry['index'] = fields['index'].strip()
    entry.update(path=path, sha256=digest(body.encode()))
    return entry


def catalog(files, content):
    reference = content('.bench/BENCH-reference.md')
    block = reference.split('<!-- bench:skills-index:start -->')[1].split(
        '<!-- bench:skills-index:end -->')[0]
    paths = re.findall(r'`(\.agents/skills/[^`]+/SKILL.md)`', block)
    entries = []
    for path in paths:
        if path not in files:
            raise ValueError('missing indexed skill: ' + path)
        entry = skill_entry(path, content)
        if 'index' not in entry:
            raise ValueError('missing skill index: ' + path)
        entries.append(entry)
    if not entries or len({e['name'] for e in entries}) != len(entries):
        raise ValueError('empty or duplicate skill catalog')
    return entries


def additional_guidance(files, content, entries):
    indexed = {entry['path'] for entry in entries}
    additional = [skill_entry(path, content) for path in sorted(files)
                  if re.fullmatch(r'\.agents/skills/[^/]+/SKILL.md', path) and path not in indexed]
    all_entries = entries + additional
    if len({entry['name'] for entry in all_entries}) != len(all_entries):
        raise ValueError('duplicate guidance name')
    return [entry['name'] for entry in additional]


def validate_packet(packet, source):
    if packet['version'] != 1 or packet['scope'] != SCOPE:
        raise ValueError('ambiguous or unsupported action scope')
    documents = packet['documents']
    if any(not (source / path).is_file() for path in LOOKUP_DOCUMENTS):
        raise ValueError('missing on-disk project lookup document')
    if not set(RULES).issubset(documents):
        raise ValueError('missing governing documents')
    if not packet['task']['source_paths']:
        raise ValueError('missing task source evidence')
    if not set(packet['task']['source_paths']).issubset(documents):
        raise ValueError('missing task source evidence')
    for path, document in documents.items():
        raw = (source / relative(path)).read_bytes()
        if 'excerpts' in document:
            lines = raw.decode().splitlines(keepends=True)
            valid = path not in RULES and document['sha256'] == digest(raw) and bool(document['excerpts'])
            last = 0
            for excerpt in document['excerpts']:
                start, end = excerpt['start_line'], excerpt['end_line']
                valid = valid and last < start <= end <= len(lines)
                valid = valid and excerpt['text'] == ''.join(lines[start - 1:end])
                last = end
        else:
            valid = document == {'text': raw.decode(), 'sha256': digest(raw)}
        if not valid:
            raise ValueError('source evidence mismatch: ' + path)
    entries = packet['catalog']
    files = {str(p.relative_to(source)): {} for p in source.rglob('SKILL.md')}
    content = lambda name: (source / name).read_text()
    actual = catalog(files, content)
    if entries != actual:
        raise ValueError('catalog differs from canonical index and skill sources')
    if packet.get('additional_guidance') != additional_guidance(files, content, actual):
        raise ValueError('additional guidance differs from frozen skill sources')
    available = {e['name'] for e in entries}
    if not set(packet['task']['protected']).issubset(available):
        raise ValueError('unavailable protected skill')
    if not set(packet['task']['optional']).issubset(available):
        raise ValueError('unavailable optional skill')
    if set(packet['task']['protected']) & set(packet['task']['optional']):
        raise ValueError('skill cannot be protected and optional')
    protected_paths = {entry['path'] for entry in entries if entry['name'] in packet['task']['protected']}
    if not protected_paths.issubset(documents):
        raise ValueError('missing protected skill body')
    if not packet['task']['brief'].strip():
        raise ValueError('missing task action')


def request(packet, model):
    questions = {}
    for entry in packet['catalog']:
        for dimension, question in (
            ('required', 'Do the supplied rules or explicit task instructions require this skill '
             'for the current action, including self-review? Optional benefit alone means no.'),
            ('useful', 'Assume this skill is not mandatory. Would its guidance materially help '
             'the current action? Judge optional benefit only.')):
            questions[entry['name'] + ':' + dimension] = {
                'type': 'noul', 'instructions': {
                    'candidate': entry, 'dimension': dimension, 'question': question}}
    return {'model': model, 'state': packet, 'questions': questions}


def context_check(packet, model):
    body = request(packet, model)
    state_bytes = len(encoded(packet))
    longest = max((len(encoded(q)) for q in body['questions'].values()), default=0)
    observed = {'state_and_longest_question_bytes': state_bytes + longest,
                'request_bytes': len(encoded(body))}
    for key, limit in CONTEXT_LIMITS.items():
        if observed[key] > limit:
            raise ValueError('context byte budget exceeded: %s=%d > %d' % (key, observed[key], limit))
    return {'observed': observed, 'limits': CONTEXT_LIMITS,
            'token_fit': 'not certified; conservative byte refusal based on English-text guidance',
            'source': 'https://docs.typesafe.ai/models.md',
            'provider_token_limits': {'request': 64000, 'state_and_longest_question': 32000}}


def document(path, task, content):
    text = content(path)
    excerpts = task.get('source_excerpts', {}).get(path)
    if not excerpts or path in RULES:
        return {'text': text, 'sha256': digest(text.encode())}
    lines = text.splitlines(keepends=True)
    result = []
    for start, end, anchor in excerpts:
        excerpt = ''.join(lines[start - 1:end])
        if not 1 <= start <= end <= len(lines) or anchor not in excerpt:
            raise ValueError('stale source excerpt: ' + path)
        result.append({'start_line': start, 'end_line': end, 'text': excerpt})
    return {'sha256': digest(text.encode()), 'excerpts': result}


def route(packet, response, model, policy=POLICY):
    protected = set(packet['task']['protected'])
    optional = set(packet['task']['optional'])
    result = {'selected': sorted(protected), 'route': 'fallback', 'reason': 'provider_failure'}
    try:
        expected = request(packet, model)['questions']
        answers = response['answers']
        if response['model'] != model or set(answers) != set(expected):
            return result
        for answer in answers.values():
            p = answer['noul']
            if answer['type'] != 'noul' or type(p) not in (int, float) or not math.isfinite(p) or not 0 <= p <= 1:
                return result
        selected = set(protected)
        uncertain = []
        for entry in packet['catalog']:
            name = entry['name']
            required = answers[name + ':required']['noul']
            useful = answers[name + ':useful']['noul']
            if name in protected:
                continue
            if name not in optional:
                if policy['required_negative_max'] < required < policy['required_positive_min']:
                    uncertain.append(name)
                if required >= policy['required_positive_min']:
                    selected.add(name)
            if useful >= policy['useful_positive_min']:
                selected.add(name)
        return {'selected': sorted(selected), 'route': 'fallback' if uncertain else 'selected',
                'reason': 'required_uncertainty' if uncertain else 'complete', 'uncertain': uncertain}
    except (KeyError, TypeError):
        return result


def freeze(root, corpus_path, output):
    root, output = Path(root), Path(output)
    if output.exists():
        raise ValueError('freeze destination already exists')
    revision, files = snapshot(root)
    corpus = read(corpus_path)
    seen, families = set(), {}
    for task in corpus:
        if not re.fullmatch(r'[a-z0-9-]+', task['id']) or task['id'] in seen:
            raise ValueError('invalid or duplicate task id')
        seen.add(task['id'])
        split = task['split']
        if split not in ('development', 'holdout'):
            raise ValueError('unknown task partition')
        if families.setdefault(task['family'], split) != split:
            raise ValueError('related task family crosses development and holdout')
        if not task['criteria'] or not task['origin'] or not task['protected_evidence']:
            raise ValueError('missing task criteria, origin, or requirement evidence')
        for path in task['source_paths'] + task['writes']:
            relative(path)
    if not seen:
        raise ValueError('empty corpus')
    source = output / 'source'
    source.mkdir(parents=True)
    for name, metadata in files.items():
        raw = git(root, 'cat-file', 'blob', metadata['oid'])
        path = source / name
        path.parent.mkdir(parents=True, exist_ok=True)
        if metadata['mode'] == '120000':
            path.symlink_to(raw.decode())
        else:
            path.write_bytes(raw)
            path.chmod(0o755 if metadata['mode'] == '100755' else 0o644)
    validate_links(source)
    content = lambda name: (source / name).read_text()
    entries = catalog(files, content)
    for task in corpus:
        paths = set(RULES) | set(task['source_paths'])
        paths.update(entry['path'] for entry in entries if entry['name'] in task['protected'])
        packet = {'version': 1, 'revision': revision, 'scope': SCOPE,
                  'task': {k: task[k] for k in ('id', 'brief', 'source_paths', 'writes',
                                              'protected', 'optional', 'protected_evidence')},
                  'catalog': entries, 'additional_guidance': additional_guidance(files, content, entries),
                  'documents': {p: document(p, task, content) for p in sorted(paths)}}
        validate_packet(packet, source)
        context_check(packet, 'jev-1.13.0')
        save(output / 'packets' / (task['id'] + '.json'), packet)
    save(output / 'corpus.json', corpus)
    manifest = file_inventory(output)
    save(output / 'freeze.json', {'version': 4, 'revision': revision, 'files': manifest})
    return {'revision': revision, 'tasks': len(corpus), 'files': len(manifest)}


def validate(root):
    root = Path(root)
    if (root / 'freeze.json').is_symlink():
        raise ValueError('invalid freeze manifest')
    manifest = read(root / 'freeze.json')
    if manifest['version'] != 4:
        raise ValueError('unsupported freeze version; prepare a new snapshot')
    actual = file_inventory(root)
    actual.pop('freeze.json', None)
    if any('special' in value for value in actual.values()):
        raise ValueError('frozen inventory contains a special file')
    if set(actual) != set(manifest['files']):
        raise ValueError('frozen file inventory changed')
    for name, expected in manifest['files'].items():
        relative(name)
        if actual[name] != expected:
            raise ValueError('frozen input changed: ' + name)
        if 'symlink' in expected and not name.startswith('source/'):
            raise ValueError('frozen evidence contains a symbolic link')
    validate_links(root / 'source')
    for task in read(root / 'corpus.json'):
        packet = read(root / 'packets' / (task['id'] + '.json'))
        validate_packet(packet, root / 'source')
        context_check(packet, 'jev-1.13.0')
    return manifest


def require_hosted_opt_in(root):
    root = Path(root)
    validate(root)
    path = root / 'source/.bench/jev-benchmark-opt-in.json'
    if not path.is_file() or read(path) != {'hosted_jev_benchmark': True}:
        raise ValueError('live run requires committed hosted Jev benchmark opt-in in the frozen source')
