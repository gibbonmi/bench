import copy
import io
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch
from pathlib import Path

from evidence import RULES, SCOPE, catalog, digest, encoded, freeze, read, request, route, save, validate, validate_packet
from runner import adapter_identities, assurance_status, changed_files, implementation_identity, invoke, run, schedule, selected_fallback
from report import export, reviews, summary
from codex_adapter import native_identity, usage_events
import jev_http


class BenchmarkTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.source = self.root / 'repo'
        self.source.mkdir()
        for name in ('AGENTS.md', '.bench/BENCH.md', '.bench/BENCH-reference.md',
                     'projects/benchkit.md', 'CONTEXT.md', 'DATA_HANDLING.md',
                     '.agents/skills/bench-craft-spec/references/ste-prose.md'):
            p = self.source / name
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_text('Project rules.\n')
        (self.source / 'scripts').mkdir()
        (self.source / 'scripts/go-build.sh').write_text('#!/bin/sh\nexit 0\n')
        (self.source / 'bin').mkdir()
        (self.source / 'bin/bench.sh').write_text(
            '#!/bin/sh\nmkdir -p .logs dist\nprintf gate > .logs/gate-fixture.out\nprintf binary > dist/bench-preflight\n')
        self.skill = '.agents/skills/bench-craft-review/SKILL.md'
        path = self.source / self.skill
        path.parent.mkdir(parents=True)
        path.write_text('---\nname: craft-review\ndescription: Review a diff.\nindex: reviewing\n---\nUse for self-review.\n')
        (self.source / '.bench/BENCH-reference.md').write_text(
            '<!-- bench:skills-index:start -->\n- reviewing → `' + self.skill +
            '`\n<!-- bench:skills-index:end -->\n')
        (self.source / 'code.txt').write_text('Actual task source.\n')
        self.task = {'id': 'case', 'family': 'case-family', 'split': 'development',
                     'origin': 'Source-derived test fixture.', 'brief': 'Write answer.md from code.txt.',
                     'source_paths': ['code.txt'], 'writes': ['answer.md'],
                     'protected': [], 'optional': [], 'protected_evidence': 'No explicit skill request.',
                     'criteria': ['The answer cites code.txt.']}
        docs = {p: {'text': (self.source / p).read_text(),
                    'sha256': digest((self.source / p).read_bytes())} for p in (*RULES, 'code.txt')}
        self.packet = {'version': 1, 'revision': 'test', 'scope': SCOPE,
                       'documents': docs, 'task': self.task, 'additional_guidance': [],
                       'catalog': catalog({self.skill: {}}, lambda p: (self.source / p).read_text())}

    def response(self, required=0.01, useful=0.5):
        return {'model': 'jev-test', 'answers': {
            'craft-review:required': {'type': 'noul', 'noul': required},
            'craft-review:useful': {'type': 'noul', 'noul': useful}}}

    def test_committed_safe_links_survive_freeze_and_tampering_is_refused(self):
        link = self.source / '.claude/skills/bench-craft-review'
        link.parent.mkdir(parents=True)
        link.symlink_to('../../.agents/skills/bench-craft-review')
        overlay = self.source / 'tests/canary/case/files/specs/retired/input.md'
        overlay.parent.mkdir(parents=True)
        overlay.symlink_to('../../../../../../code.txt')
        frozen = self.frozen()
        copied = frozen / 'source' / link.relative_to(self.source)
        self.assertTrue(copied.is_symlink(), 'freeze lost the tracked harness skill link')
        self.assertEqual((frozen / 'source' / overlay.relative_to(self.source)).read_text(), 'Actual task source.\n')
        self.assertEqual((copied / 'SKILL.md').read_bytes(), (link / 'SKILL.md').read_bytes())
        validate(frozen)
        copied.unlink()
        copied.symlink_to('/tmp')
        with self.assertRaises(ValueError):
            validate(frozen)

    def test_unsafe_committed_link_refuses_freeze(self):
        (self.source / 'escape').symlink_to('/etc/passwd')
        with self.assertRaisesRegex(ValueError, 'unsafe snapshot link'):
            self.frozen()

    def test_missing_rules_and_source_refuse_before_a_request(self):
        validate_packet(self.packet, self.source)
        for name in (*RULES, 'code.txt'):
            packet = copy.deepcopy(self.packet)
            del packet['documents'][name]
            with self.subTest(name=name), self.assertRaisesRegex(ValueError, 'missing'):
                validate_packet(packet, self.source)

    def test_changed_source_and_description_refuse(self):
        packet = copy.deepcopy(self.packet)
        packet['documents']['code.txt']['text'] = 'Invented source.'
        with self.assertRaisesRegex(ValueError, 'source evidence mismatch'):
            validate_packet(packet, self.source)
        packet = copy.deepcopy(self.packet)
        packet['catalog'][0]['description'] = 'Invented trigger.'
        with self.assertRaisesRegex(ValueError, 'catalog differs'):
            validate_packet(packet, self.source)

    def test_self_review_scope_cannot_be_removed(self):
        packet = copy.deepcopy(self.packet)
        packet['scope']['included'].remove('self-review')
        with self.assertRaisesRegex(ValueError, 'scope'):
            validate_packet(packet, self.source)

    def test_required_and_useful_are_independent_questions(self):
        body = request(self.packet, 'jev-test')
        self.assertEqual(len(body['questions']), 2)
        self.assertNotEqual(body['questions']['craft-review:required']['instructions']['question'],
                            body['questions']['craft-review:useful']['instructions']['question'])
        self.assertEqual(body['state']['documents'], self.packet['documents'])

    def test_optional_uncertainty_does_not_fallback(self):
        self.assertEqual(route(self.packet, self.response(), 'jev-test')['route'], 'selected')
        self.assertEqual(route(self.packet, self.response(0.5), 'jev-test')['route'], 'fallback')
        self.assertEqual(route(self.packet, self.response(0.8), 'jev-test')['selected'], ['craft-review'])
        self.assertEqual(route(self.packet, self.response(0.2), 'jev-test')['route'], 'selected')
        self.packet['task']['optional'] = ['craft-review']
        self.assertEqual(route(self.packet, self.response(0.49), 'jev-test')['route'], 'selected')

    def test_bad_provider_answers_fallback_and_preserve_protection(self):
        self.packet['task']['protected'] = ['craft-review']
        for bad in (None, {}, {'model': 'wrong'}, self.response(float('nan')),
                    self.response(True), self.response(-0.1), self.response(1.1)):
            with self.subTest(response=bad):
                routed = route(self.packet, bad, 'jev-test')
                self.assertEqual(routed['route'], 'fallback')
                self.assertEqual(routed['selected'], ['craft-review'])

    def test_invalid_fallback_cannot_drop_a_protected_skill(self):
        self.packet['task']['protected'] = ['craft-review']
        with self.assertRaisesRegex(ValueError, 'invalid fallback'):
            selected_fallback({'status': 'completed', 'response': {'selected': [], 'abstain': False}}, self.packet)

    def test_schedule_is_paired_repeated_and_balances_order(self):
        rows = schedule([self.task], 4, 17)
        self.assertEqual(rows, schedule([self.task], 4, 17))
        self.assertEqual(len(rows), 8)
        self.assertEqual(sum(r['condition'] == 'baseline' for r in rows[::2]), 2)
        self.assertEqual({r['repeat'] for r in rows}, {0, 1, 2, 3})

    def test_wire_bytes_and_timing_have_one_owner(self):
        adapter = self.root / 'echo.py'
        adapter.write_text('import sys\nsys.stdout.buffer.write(sys.stdin.buffer.read())\n')
        result = invoke([sys.executable, str(adapter)], self.response(), self.root, self.root / 'call', 2)
        self.assertEqual((self.root / 'call/request.json').read_bytes(), (self.root / 'call/stdout').read_bytes())
        self.assertEqual(result['request_sha256'], digest((self.root / 'call/stdout').read_bytes()))
        self.assertEqual(result['clock'], 'monotonic_ns')
        self.assertGreater(result['elapsed_ns'], 0)

    def test_timeout_is_retained(self):
        result = invoke([sys.executable, '-c', 'import time; time.sleep(30)'], {},
                        self.root, self.root / 'timeout', 0.05)
        self.assertEqual(result['status'], 'timeout')
        self.assertTrue((self.root / 'timeout/result.json').exists())

    def test_empty_response_is_failure_and_failed_usage_is_retained(self):
        empty = invoke([sys.executable, '-c', 'pass'], {}, self.root, self.root / 'empty', 2)
        self.assertEqual(empty['status'], 'failed')
        failed = invoke([sys.executable, '-c', 'print(\'{"usage":{"input_tokens":10}}\'); raise SystemExit(1)'],
                        {}, self.root, self.root / 'failed', 2)
        self.assertEqual(failed['status'], 'failed')
        self.assertEqual(failed['response']['usage']['input_tokens'], 10)

    def test_all_native_turn_usage_is_preserved_without_cost_rederivation(self):
        native = [{'type': 'turn.completed', 'usage': {'input_tokens': n, 'cached_input_tokens': 4,
                                                       'output_tokens': 2}} for n in (10, 20)]
        events = usage_events(native, 'session', {'producer': 'fixture', 'native': 'events'})
        self.assertEqual([e['usage']['input_total'] for e in events], [10, 20])
        self.assertTrue(all(e['usage']['total_semantics'] == 'inclusive' for e in events))
        self.assertNotEqual(events[0]['event_id'], events[1]['event_id'])

    def test_nested_native_usage_is_never_reported_as_complete(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        export(frozen, output, self.root / 'export', 'fixture-123')
        for path in (self.root / 'export').glob('*.json'):
            for attempt in read(path)['attempts']:
                if attempt['role'] == 'implementation':
                    self.assertTrue(attempt['cost'].get('other'), 'unobserved native usage priced as complete')

    def test_http_adapter_posts_saved_bytes_and_offline_never_posts(self):
        save(self.source / '.bench/jev-benchmark-opt-in.json', {'hosted_jev_benchmark': True})
        frozen = self.frozen()
        body = encoded(request(self.packet, 'jev-test'))
        posted = []
        class Response:
            def __enter__(self):
                return self
            def __exit__(self, *args):
                return False
            def read(self):
                return b'{"answers": {}}'
        def send(req, timeout):
            posted.append(req.data)
            self.assertEqual(req.full_url, 'https://api.typesafe.ai/v1/systemone')
            return Response()
        with patch.dict(os.environ, {'TYPESAFE_API_KEY': 'dummy-fixture-key', 'JEV_BENCHMARK_FROZEN_ROOT': str(frozen)}, clear=True), \
                patch.object(sys, 'stdin', io.TextIOWrapper(io.BytesIO(body))), \
                patch.object(sys, 'stdout', io.TextIOWrapper(io.BytesIO())), \
                patch('urllib.request.urlopen', send):
            jev_http.main()
            os.environ['BENCH_OFFLINE'] = '1'
            with self.assertRaisesRegex(SystemExit, 'suppresses hosted Jev'):
                jev_http.main()
        self.assertEqual(posted, [body])

    def test_interrupted_native_capture_keeps_known_and_unknown_usage(self):
        from codex_adapter import captured_native
        output = self.root / 'native-partial'
        output.mkdir()
        events = [{'type': 'thread.started', 'thread_id': 'parent'},
                  {'type': 'turn.completed', 'usage': {'input_tokens': 100,
                   'cached_input_tokens': 40, 'output_tokens': 5}}]
        (output / 'events.jsonl').write_text(''.join(json.dumps(e) + '\n' for e in events) + '{partial')
        result = captured_native(output, {'harness': 'fixture', 'native_line': {'model': 'fixture', 'effort': 'high'}}, None)
        self.assertEqual(result['assessment_usage'][0]['usage']['input_total'], 100)
        self.assertTrue(result['assessment_usage'][-1]['usage']['unknown'])
        self.assertIn('unverified', result['usage_scope'])

    def test_codex_adapter_with_native_cli_fixture(self):
        binary = self.root / 'codex'
        binary.write_text('#!' + sys.executable + '''
import json, sys
from pathlib import Path
if '--version' in sys.argv:
 print('codex-cli fixture')
else:
 p=json.load(sys.stdin)
 Path(sys.argv[sys.argv.index('--output-last-message')+1]).write_text(json.dumps({'summary':'fixture','rescued_skills':[],'verification':[]}))
 print(json.dumps({'type':'thread.started','thread_id':'fixture-session'}))
 print(json.dumps({'type':'turn.completed','usage':{'input_tokens':100,'cached_input_tokens':40,'output_tokens':5}}))
''')
        binary.chmod(0o755)
        payload = {'role': 'task', 'execution': {'harness': 'codex-cli fixture',
                   'native_line': {'model': 'fixture', 'effort': 'high'}, 'service_tier': 'default',
                   'native_executable': native_identity(binary)}}
        adapter = str(Path(__file__).with_name('codex_adapter.py').resolve())
        with patch.dict(os.environ, {'PATH': str(self.root) + os.pathsep + os.environ['PATH']}):
            result = invoke([sys.executable, '-B', adapter], payload, self.source, self.root / 'native', 3)
        self.assertEqual(result['status'], 'completed')
        self.assertEqual(result['response']['session_id'], 'fixture-session')
        self.assertEqual(result['response']['assessment_usage'][0]['usage']['input_total'], 100)
        self.assertEqual(result['response']['assessment_usage'][0]['usage']['input_cached'], 40)
        binary.write_text(binary.read_text() + '\n# Same version, different bytes.\n')
        changed = invoke([sys.executable, '-B', adapter], payload, self.source, self.root / 'changed-native', 3)
        self.assertEqual(changed['status'], 'failed')
        self.assertIn('native executable identity', (self.root / 'changed-native/stderr').read_text())
        self.assertFalse((self.root / 'changed-native/events.jsonl').exists())

    def commit_source(self):
        subprocess.run(['git', 'init', '-q', str(self.source)], check=True)
        subprocess.run(['git', 'add', '.'], cwd=self.source, check=True)
        subprocess.run(['git', '-c', 'user.name=Fixture', '-c', 'user.email=fixture@example.invalid',
                        '-c', 'commit.gpgsign=false', 'commit', '-qm', 'fixture'], cwd=self.source, check=True)

    def frozen(self):
        self.commit_source()
        save(self.root / 'corpus.json', [self.task])
        freeze(self.source, self.root / 'corpus.json', self.root / 'frozen')
        return self.root / 'frozen'

    def fixture_config(self):
        frozen = self.frozen()
        adapter = self.root / 'adapter.py'
        adapter.write_text('''import json, os, sys
from pathlib import Path
p=json.load(sys.stdin)
Path(os.environ['JEV_BENCHMARK_CALL_DIR'], 'events.jsonl').write_text('native evidence\\n')
if 'questions' in p:
 print(json.dumps({'model':p['model'],'answers':{k:{'type':'noul','noul':0.01} for k in p['questions']}}))
elif p['role']=='fallback':
 print(json.dumps({'selected':[], 'abstain':False}))
else:
 Path('answer.md').write_text('Source: code.txt.\\n')
 print(json.dumps({'native_usage':None,'rescued_skills':[], 'self_review':['Checked answer against source.'], 'verification':['Read code.txt.']}))
''')
        command = [sys.executable, str(adapter)]
        config = {'adapters': {k: command for k in ('jev', 'fallback', 'task')},
                  'input_freeze_sha256': digest(encoded(read(frozen / 'freeze.json'))),
                  'jev_model': 'jev-test', 'native_line': 'fixture', 'harness': 'fixture',
                  'service_tier': 'fixture',
                  'repetitions': 2, 'seed': 17, 'timeout_seconds': 3,
                  'max_adapter_calls': 8, 'budget_reference': 'offline fixture; zero paid calls'}
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        config['implementation_sha256'] = implementation_identity()
        return frozen, adapter, config

    def test_verification_artifacts_are_separate_and_baseline_red_stops_calls(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        self.assertTrue((output / 'preflight/result.json').is_file(), 'no clean baseline gate preflight')
        for path in output.glob('trial-*/trial.json'):
            trial = read(path)
            self.assertEqual(trial['status'], 'completed')
            self.assertFalse((path.parent / 'workspace/.logs').exists())
            self.assertTrue(list((path.parent / 'verification').glob('*/workspace/.logs/gate-fixture.out')))
        script = frozen / 'source/bin/bench.sh'
        script.write_text('#!/bin/sh\nexit 1\n')
        manifest = read(frozen / 'freeze.json')
        manifest['files']['source/bin/bench.sh']['sha256'] = digest(script.read_bytes())
        (frozen / 'freeze.json').write_bytes(encoded(manifest))
        config['input_freeze_sha256'] = digest(encoded(manifest))
        stopped = run(frozen, self.root / 'baseline-red', config, None, offline=True)
        self.assertEqual(stopped['stopped'], 'invalid_baseline')
        self.assertFalse(list((self.root / 'baseline-red').glob('trial-*')))

    def test_interrupted_trial_is_sealed_and_stops_future_calls(self):
        frozen, adapter, config = self.fixture_config()
        adapter.write_text('import os, signal, time\nos.kill(os.getppid(), signal.SIGINT)\ntime.sleep(30)\n')
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        output = self.root / 'run'
        result = run(frozen, output, config, None, offline=True)
        self.assertEqual(result['stopped'], 'interrupted')
        self.assertEqual(result['trials'], 1)
        self.assertEqual(read(output / 'trial-000/trial.json')['attempts'], ['task'])
        self.assertEqual([r['status'] for r in summary(output)['trials']],
                         ['interrupted', 'missing', 'missing', 'missing'])
        export(frozen, output, self.root / 'export', 'fixture-123')
        record = read(next((self.root / 'export').glob('*.json')))
        self.assertEqual(record['attempts'][0]['state'], 'cancelled')
        self.assertTrue(record['attempts'][0]['usage'][0]['usage']['unknown'])

    def test_interruption_remains_authoritative_after_outside_write(self):
        frozen, adapter, config = self.fixture_config()
        source = adapter.read_text().replace('import json, os, sys', 'import json, os, sys, signal, time')
        source = source.replace(" Path('answer.md').write_text('Source: code.txt.\\n')",
                                " Path('.unapproved').write_text('partial'); os.kill(os.getppid(), signal.SIGINT); time.sleep(30)")
        adapter.write_text(source)
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        output = self.root / 'run'
        result = run(frozen, output, config, None, offline=True)
        self.assertEqual(result['trials'], 1, 'later calls ran after the task interruption')
        self.assertEqual(result['stopped'], 'interrupted')
        trial = read(output / 'trial-000/trial.json')
        self.assertEqual(trial['status'], 'interrupted')
        self.assertEqual(trial['outside_write_fence'], ['.unapproved'])
        self.assertEqual([row['status'] for row in summary(output)['trials']],
                         ['interrupted', 'missing', 'missing', 'missing'])
        self.assertEqual(set(read(output / 'ledger.json')['trials']), {'trial-000'})

    def test_every_report_rejects_changed_verification_artifacts(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        verification = next((output / 'trial-000/verification').iterdir())
        for name in ('result.json', 'stdout', 'bootstrap/result.json', 'workspace/.logs/gate-fixture.out'):
            path = verification / name
            original = path.read_bytes()
            path.write_bytes(original + b' ')
            for report in self.report_consumers(frozen, output):
                with self.assertRaisesRegex(ValueError, 'changed after capture'):
                    report()
            path.write_bytes(original)

    def test_verifier_cli_accepts_native_metadata_and_reuses_its_receipt(self):
        output = self.root / 'verification'
        for name in ('.git', '.agents', '.codex'):
            (output / name).mkdir(parents=True)
        command = [sys.executable, '-B', str(Path(__file__).with_name('verification.py')),
                   '--source', str(self.source), '--out', str(output)]
        first = subprocess.run(command, capture_output=True, text=True)
        self.assertEqual(first.returncode, 0, first.stderr)
        second = subprocess.run(command, capture_output=True, text=True)
        self.assertEqual(second.returncode, 0, second.stderr)
        self.assertEqual(json.loads(first.stdout.splitlines()[0])['record'],
                         json.loads(second.stdout.splitlines()[0])['record'])
        (output / 'unknown').mkdir()
        corrupt = subprocess.run(command, capture_output=True, text=True)
        self.assertNotEqual(corrupt.returncode, 0, 'unknown output escaped validation')

    def test_verifier_owns_cache_writes_outside_authored_source(self):
        from verification import ensure_verified
        probe = self.root / 'cache_probe.py'
        probe.write_text("import os, pathlib\n"
                         "for key in ('HOME', 'BENCH_HOME', 'GOCACHE', 'XDG_CACHE_HOME', 'TMPDIR'):\n"
                         " p = pathlib.Path(os.environ[key]); p.mkdir(parents=True, exist_ok=True)\n"
                         " assert pathlib.Path(os.environ['EXPECTED_OUTPUT']) in p.parents, key\n"
                         " (p / 'probe').write_text(key)\n")
        (self.source / 'scripts/go-build.sh').write_text(
            '#!/bin/sh\n' + sys.executable + ' ' + str(probe) + '\n')
        output = self.root / 'verification'
        with patch.dict(os.environ, {'EXPECTED_OUTPUT': str(output)}):
            result = ensure_verified(self.source, output)
        self.assertEqual(result['status'], 'completed',
                         (output / result['record'] / 'stderr').read_text())
        self.assertFalse((self.source / '.logs').exists())
        if shutil.which('go'):
            record = output / result['record']
            child = {key: os.environ[key] for key in ('PATH',) if key in os.environ}
            child['HOME'] = str(record / 'runtime/home')
            settings = json.loads(subprocess.check_output(
                ['go', 'env', '-json', 'GOPATH', 'GOMODCACHE', 'GOPROXY'], env=child, text=True))
            expected = json.loads(subprocess.check_output(
                ['go', 'env', '-json', 'GOPATH', 'GOMODCACHE'], text=True))
            self.assertEqual(settings, dict(expected, GOPROXY='off'))
            self.assertEqual(subprocess.check_output(['go', 'telemetry'], env=child, text=True).strip(), 'off')

    def test_verification_receipt_cannot_replace_independent_final_gate(self):
        from verification import ensure_verified
        output = self.root / 'verification'
        first = ensure_verified(self.source, output)
        self.assertEqual(first['status'], 'completed')
        second = ensure_verified(self.source, output, force=True)
        self.assertNotEqual(first['record'], second['record'])
        (self.source / 'bin/bench.sh').write_text('#!/bin/sh\nexit 1\n')
        final = ensure_verified(self.source, output, force=True)
        self.assertEqual(final['status'], 'failed')
        self.assertFalse((self.source / '.logs').exists())

    def test_gate_cannot_grade_rewritten_source(self):
        from verification import ensure_verified
        (self.source / 'bin/bench.sh').write_text('#!/bin/sh\nprintf changed > code.txt\n')
        result = ensure_verified(self.source, self.root / 'verification')
        self.assertEqual(result['status'], 'invalid_verification')
        self.assertEqual((self.source / 'code.txt').read_text(), 'Actual task source.\n')

    def test_source_and_unapproved_ignored_writes_remain_fenced(self):
        frozen, adapter, config = self.fixture_config()
        adapter.write_text(adapter.read_text().replace(
            "Path('answer.md').write_text('Source: code.txt.\\n')",
            "Path('answer.md').write_text('Source: code.txt.\\n'); Path('.unapproved').write_text('x')"))
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        run(frozen, self.root / 'run', config, None, offline=True)
        for path in (self.root / 'run').glob('trial-*/trial.json'):
            self.assertEqual(read(path)['status'], 'invalid_write_fence')
            self.assertEqual(read(path)['outside_write_fence'], ['.unapproved'])

    def test_interrupt_retains_attempt_result(self):
        with patch('subprocess.Popen.communicate', side_effect=KeyboardInterrupt):
            result = invoke([sys.executable, '-c', 'import time; time.sleep(30)'], {},
                            self.root, self.root / 'interrupted', 3)
        self.assertEqual(result['status'], 'interrupted')
        self.assertEqual(read(self.root / 'interrupted/result.json')['status'], 'interrupted')

    def test_complete_offline_pair_keeps_quality_unknown(self):
        frozen, adapter, config = self.fixture_config()
        result = run(frozen, self.root / 'run', config, None, offline=True)
        self.assertEqual(result['trials'], 4)
        for path in (self.root / 'run').glob('trial-*/trial.json'):
            trial = read(path)
            self.assertEqual(trial['status'], 'completed')
            self.assertIsNone(trial['quality'])
            self.assertEqual(trial['outside_write_fence'], [])
        with self.assertRaisesRegex(ValueError, 'authorization'):
            run(frozen, self.root / 'live', config, None)
        with self.assertRaisesRegex(ValueError, 'already exists'):
            run(frozen, self.root / 'run', config, None, offline=True)
        overview = summary(self.root / 'run')
        self.assertFalse(overview['adoption_evidence'])
        reviews(frozen, self.root / 'run', self.root / 'review')
        for path in (self.root / 'review/packets').glob('*.json'):
            self.assertNotIn('condition', read(path))
            self.assertNotIn('selected', read(path))
            self.assertEqual(read(path)['validity']['outside_write_fence'], [])
            self.assertNotIn('self_review_claims', read(path))
            self.assertNotIn('status', read(path)['validity'])
        exported = export(frozen, self.root / 'run', self.root / 'export', 'fixture-123')
        self.assertEqual(len(exported['assessment_records']), 4)
        for path in (self.root / 'export').glob('*.json'):
            record = read(path)
            self.assertEqual(record['state'], 'incomplete')
            self.assertEqual(record['quality'], {})
            self.assertTrue(all('elapsed_ns' in attempt['measures'] for attempt in record['attempts']))
        with self.assertRaisesRegex(ValueError, 'committed hosted Jev'):
            run(frozen, self.root / 'without-opt-in', config, digest(encoded(config)))
        adapter.write_text(adapter.read_text() + '\n# Changed adapter bytes.\n')
        with self.assertRaisesRegex(ValueError, 'adapter identity'):
            run(frozen, self.root / 'changed-adapter', config, None, offline=True)
        adapter.write_text(adapter.read_text().replace("['Checked answer against source.']", "['  ']"))
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        run(frozen, self.root / 'no-review-run', config, None, offline=True)
        for path in (self.root / 'no-review-run').glob('trial-*/trial.json'):
            self.assertEqual(read(path)['status'], 'missing_assurance_evidence')

    def report_consumers(self, frozen, run):
        return (lambda: reviews(frozen, run, self.root / 'review'),
                lambda: export(frozen, run, self.root / 'export', 'fixture-123'),
                lambda: summary(run))

    def test_every_report_refuses_changed_scored_evidence(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        for name in ('plan.json', 'preflight/result.json', 'preflight/stdout',
                     'preflight/workspace/.logs/gate-fixture.out', 'trial-000/trial.json', 'trial-000/task/result.json',
                     'trial-000/task/request.json', 'trial-000/task/events.jsonl',
                     'trial-000/task/stdout', 'trial-000/task/stderr',
                     'trial-001/selection.json', 'trial-001/jev/result.json',
                     'trial-000/workspace/answer.md'):
            path = output / name
            original = path.read_bytes()
            path.write_bytes(original + b' ')
            try:
                for report in self.report_consumers(frozen, output):
                    with self.subTest(path=name, report=report), self.assertRaisesRegex(ValueError, 'changed after capture'):
                        report()
            finally:
                path.write_bytes(original)
        mode_path = output / 'trial-000/workspace/code.txt'
        original_mode = mode_path.stat().st_mode
        mode_path.chmod(0o755)
        for report in self.report_consumers(frozen, output):
            with self.assertRaisesRegex(ValueError, 'changed after capture'):
                report()
        mode_path.chmod(original_mode)
        for name in ('trial-000/task/unexpected', 'trial-000/workspace/unexpected'):
            path = output / name
            path.write_text('unreviewed')
            for report in self.report_consumers(frozen, output):
                with self.assertRaisesRegex(ValueError, 'changed after capture'):
                    report()
            path.unlink()
        self.assertFalse((self.root / 'review').exists())
        self.assertFalse((self.root / 'export').exists())
        self.assertEqual(len(summary(output)['trials']), 4)

    def test_every_report_refuses_removed_or_unsealed_trials(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        path = output / 'trial-000'
        path.rename(self.root / 'removed-trial')
        for report in self.report_consumers(frozen, output):
            with self.assertRaisesRegex(ValueError, 'trial inventory changed'):
                report()
        (self.root / 'removed-trial').rename(path)
        (output / 'trial-004').mkdir()
        for report in self.report_consumers(frozen, output):
            with self.assertRaisesRegex(ValueError, 'unsealed interrupted evidence'):
                report()

    def test_scored_values_cannot_be_rewritten_after_capture(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        mutations = (
            ('trial-000/trial.json', lambda data: data.update(attempts=[], status='invalid_provider_input')),
            ('trial-000/trial.json', lambda data: data.update(rescued_skills=['invented'])),
            ('trial-000/task/result.json', lambda data: data.update(elapsed_ns=data['elapsed_ns'] + 1000000)),
            ('trial-000/task/result.json', lambda data: data['response'].update(
                assessment_usage=[{'usage': {'input_total': 1}}], actual_charges=[{'amount': 0.01}])),
            ('plan.json', lambda data: data['config'].update(rates={'native': {'input_uncached': 1}})),
        )
        for name, mutate in mutations:
            path = output / name
            original = path.read_bytes()
            data = read(path)
            mutate(data)
            path.write_bytes(encoded(data))
            try:
                for report in self.report_consumers(frozen, output):
                    with self.subTest(path=name), self.assertRaisesRegex(ValueError, 'changed after capture'):
                        report()
            finally:
                path.write_bytes(original)
    def test_plan_configuration_fingerprint_is_independently_checked(self):
        frozen, adapter, config = self.fixture_config()
        output = self.root / 'run'
        run(frozen, output, config, None, offline=True)
        path = output / 'plan.json'
        plan = read(path)
        plan['config']['rates'] = {'native': {'input_uncached': 1}}
        path.write_bytes(encoded(plan))
        ledger_path = output / 'ledger.json'
        ledger = read(ledger_path)
        ledger['plan_sha256'] = digest(path.read_bytes())
        ledger_path.write_bytes(encoded(ledger))
        for report in self.report_consumers(frozen, output):
            with self.assertRaisesRegex(ValueError, 'configuration fingerprint mismatch'):
                report()

    def test_provider_rejection_is_captured_and_withheld_from_quality_review(self):
        frozen, adapter, config = self.fixture_config()
        adapter.write_text(adapter.read_text().replace(
            "print(json.dumps({'model':p['model'],'answers':{k:{'type':'noul','noul':0.01} for k in p['questions']}}))",
            "print(json.dumps({'http_status':422,'provider_body':'invalid input'}))"))
        config['adapter_sha256'] = adapter_identities(config['adapters'])
        config['seed'] = 17
        output = self.root / 'run'
        result = run(frozen, output, config, None, offline=True)
        self.assertEqual(result['stopped'], 'invalid_provider_input')
        self.assertEqual(result['trials'], 2)
        reviewed = reviews(frozen, output, self.root / 'review')
        self.assertEqual(reviewed['unique_review_packets'], 1)
        rows = summary(output)['trials']
        self.assertEqual([r['status'] for r in rows], ['completed', 'invalid_provider_input', 'missing', 'missing'])
        self.assertEqual(len(export(frozen, output, self.root / 'export', 'fixture-123')['assessment_records']), 2)

    def test_frozen_inventory_rejects_extra_and_modified_files(self):
        frozen = self.frozen()
        validate(frozen)
        (frozen / 'unexpected').write_text('unreviewed state')
        with self.assertRaisesRegex(ValueError, 'inventory changed'):
            validate(frozen)

    def test_related_family_cannot_cross_holdout_partition(self):
        corpus = [self.task, dict(self.task, id='related', split='holdout')]
        self.commit_source()
        save(self.root / 'corpus.json', corpus)
        with self.assertRaisesRegex(ValueError, 'family crosses'):
            freeze(self.source, self.root / 'corpus.json', self.root / 'frozen')

    def test_required_documents_are_present_in_produced_packet(self):
        frozen = self.frozen()
        packet = read(frozen / 'packets/case.json')
        for name in ('AGENTS.md', '.bench/BENCH.md', 'projects/benchkit.md', 'DATA_HANDLING.md',
                     '.agents/skills/bench-craft-spec/references/ste-prose.md'):
            self.assertIn(name, packet['documents'])
            self.assertEqual(packet['documents'][name]['text'], (self.source / name).read_text())
        for name in ('.bench/BENCH-reference.md', 'CONTEXT.md'):
            self.assertEqual((frozen / 'source' / name).read_bytes(), (self.source / name).read_bytes())

    def test_frozen_modes_and_special_files_are_not_invisible(self):
        frozen = self.frozen()
        source = frozen / 'source/code.txt'
        source.chmod(0o755)
        with self.assertRaisesRegex(ValueError, 'frozen input changed'):
            validate(frozen)
        source.chmod(0o644)
        os.mkfifo(frozen / 'unexpected-fifo')
        with self.assertRaisesRegex(ValueError, 'special file'):
            validate(frozen)
        workspace = self.root / 'work'
        workspace.mkdir()
        os.mkfifo(workspace / 'outside-fifo')
        self.assertIn('outside-fifo', changed_files(self.source, workspace))

    def test_context_limit_refuses_without_provider_call(self):
        from evidence import context_check
        packet = copy.deepcopy(self.packet)
        packet['task']['brief'] = 'X' * 120001
        with self.assertRaisesRegex(ValueError, 'context byte budget exceeded'):
            context_check(packet, 'jev-test')

    def test_adapter_identity_changes_when_bytes_change(self):
        adapter = self.root / 'mutable.py'
        adapter.write_text('print("reviewed")')
        commands = {'task': [sys.executable, str(adapter)]}
        before = adapter_identities(commands)
        adapter.write_text('print("changed")')
        self.assertNotEqual(adapter_identities(commands), before)

    def test_frozen_phase_guidance_can_be_rescued_without_becoming_a_candidate(self):
        skill = self.source / '.agents/skills/bench-debug/SKILL.md'
        skill.parent.mkdir(parents=True)
        skill.write_text('---\nname: bench-debug\ndescription: Use when something breaks.\n---\nDiagnose before changing code.\n')
        frozen = self.frozen()
        packet = read(frozen / 'packets/case.json')
        response = {'self_review': ['Checked the fix.'], 'verification': ['Ran the reproducer.'],
                    'rescued_skills': ['bench-debug']}
        self.assertEqual(assurance_status(response, packet, []), 'completed')
        self.assertNotIn('bench-debug:required', request(packet, 'jev-test')['questions'])
        missing = copy.deepcopy(packet)
        missing['additional_guidance'] = []
        with self.assertRaisesRegex(ValueError, 'additional guidance differs'):
            validate_packet(missing, frozen / 'source')
        response['rescued_skills'] = ['invented']
        self.assertEqual(assurance_status(response, packet, []), 'invalid_rescue_evidence')

    def test_assurance_and_rescues_require_meaningful_valid_claims(self):
        valid = {'self_review': ['Checked the diff.'], 'verification': ['Ran the named test.'],
                 'rescued_skills': ['craft-review']}
        self.assertEqual(assurance_status(valid, self.packet, []), 'completed')
        for key in ('self_review', 'verification'):
            for value in ([], [''], ['  '], [None], [1], 'a claim'):
                with self.subTest(key=key, value=value):
                    self.assertEqual(assurance_status(dict(valid, **{key: value}), self.packet, []),
                                     'missing_assurance_evidence')
        for value in (None, 'craft-review', ['invented'], ['craft-review', 'craft-review'], [1]):
            with self.subTest(rescues=value):
                self.assertEqual(assurance_status(dict(valid, rescued_skills=value), self.packet, []),
                                 'invalid_rescue_evidence')
        self.assertEqual(assurance_status(valid, self.packet, ['craft-review']), 'invalid_rescue_evidence')

    def test_workspace_mode_changes_are_visible_to_write_fence(self):
        import shutil
        workspace = self.root / 'workspace'
        shutil.copytree(self.source, workspace)
        self.assertEqual(changed_files(self.source, workspace), [])
        (workspace / 'code.txt').chmod(0o755)
        self.assertEqual(changed_files(self.source, workspace), ['code.txt'])


if __name__ == '__main__':
    unittest.main()
