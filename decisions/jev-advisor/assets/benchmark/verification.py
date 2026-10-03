"""Run the mandatory gate in a separate copy of the exact task inputs."""

import argparse
import json
import os
import shutil
import subprocess
import tempfile
import time
from pathlib import Path

from evidence import digest, encoded, file_inventory, inventory_changes, read, save, validate_links
from execution import execute

BUILD_COMMAND = ['bash', 'scripts/go-build.sh', '--manifest-dir', 'dist', '.', 'dist/bench']
GATE_COMMAND = ['bash', 'bin/bench.sh', 'gate']
GATE_TIMEOUT = 600


def copy_workspace(source, destination):
    validate_links(source)
    shutil.copytree(source, destination, symlinks=True, ignore=shutil.ignore_patterns('.git'))
    environment = dict(os.environ, GIT_CONFIG_NOSYSTEM='1', GIT_CONFIG_GLOBAL='/dev/null')
    for args in (['init', '-q'], ['add', '-f', '.'],
                 ['-c', 'user.name=Benchmark', '-c', 'user.email=benchmark@example.invalid',
                  '-c', 'commit.gpgsign=false', '-c', 'core.hooksPath=/dev/null',
                  'commit', '-qm', 'Frozen benchmark source']):
        subprocess.run(['git', *args], cwd=destination, env=environment, check=True,
                       stdout=subprocess.DEVNULL, stderr=subprocess.PIPE)


def verify(source, output, timeout=GATE_TIMEOUT):
    deadline = time.monotonic() + timeout
    def remaining():
        return max(0.001, deadline - time.monotonic())
    source, output = Path(source).resolve(), Path(output).resolve()
    output.mkdir(parents=True, exist_ok=False)
    before = file_inventory(source)
    work = output / 'workspace'
    copy_workspace(source, work)
    if file_inventory(work) != before:
        raise ValueError('verification copy differs from task inputs')
    environment = {k: v for k, v in os.environ.items()
                   if not k.startswith(('BENCH_', 'JEV_', 'GIT_')) and k != 'TYPESAFE_API_KEY'}
    environment.update(BENCH_KIT=str(work), BENCH_OFFLINE='1', GOPROXY='off')
    go_available = shutil.which('go') is not None
    if go_available:
        go_paths = subprocess.check_output(['go', 'env', '-json', 'GOPATH', 'GOMODCACHE', 'GOFLAGS'],
                                           cwd=work, env=environment, text=True, timeout=remaining())
        environment.update(json.loads(go_paths))
    # Bench derives its cache from the child home, even when GOCACHE is supplied.
    for key, name in (('HOME', 'home'), ('BENCH_HOME', 'bench'), ('GOCACHE', 'go-build'),
                      ('XDG_CACHE_HOME', 'cache'), ('TMPDIR', 'tmp')):
        directory = output / 'runtime' / name
        directory.mkdir(parents=True)
        environment[key] = str(directory)
    if go_available:
        # The gate's environment filter retains HOME but drops Go path overrides.
        environment.pop('GOENV', None)
        environment.pop('XDG_CONFIG_HOME', None)
        subprocess.run(['go', 'telemetry', 'off'], cwd=work,
                       env=dict(environment, GOTOOLCHAIN='local'), check=True,
                       capture_output=True, timeout=remaining())
        settings = {key: environment[key] for key in ('GOPATH', 'GOMODCACHE', 'GOFLAGS', 'GOPROXY')}
        subprocess.run(['go', 'env', '-w', *(key + '=' + value for key, value in settings.items())],
                       cwd=work, env=dict(environment, GOTOOLCHAIN='local'), check=True,
                       capture_output=True, timeout=remaining())
    bootstrap = output / 'bootstrap'
    bootstrap.mkdir()
    built = execute(BUILD_COMMAND, work, bootstrap, remaining(), environment)
    save(bootstrap / 'result.json', built)
    if built['status'] == 'completed':
        result = execute(GATE_COMMAND, work, output, remaining(), environment)
        result['phase'] = 'gate'
    else:
        result = dict(built, command=GATE_COMMAND, phase='bootstrap')
        for name in ('stdout', 'stderr'):
            shutil.copyfile(bootstrap / name, output / name)
    result['bootstrap'] = built
    after = file_inventory(work)
    changed = inventory_changes(before, after)
    result.update(input_sha256=digest(encoded(before)), input_inventory=before,
                  workspace_sha256=digest(encoded(after)), generated_or_changed=changed)
    if any(after.get(name) != value for name, value in before.items()):
        result.update(status='invalid_verification', error='gate changed its source inputs')
    if file_inventory(source) != before:
        result.update(status='invalid_verification', error='task inputs changed during verification')
    save(output / 'result.json', result)
    return result


def verified_result(directory):
    result = read(directory / 'result.json')
    if (result['command'] != GATE_COMMAND or result['bootstrap']['command'] != BUILD_COMMAND
            or result['clock'] != 'monotonic_ns'
            or result['workspace_sha256'] != digest(encoded(file_inventory(directory / 'workspace')))
            or result['input_sha256'] != digest(encoded(result['input_inventory']))):
        raise ValueError('verification evidence changed')
    return result


def ensure_verified(source, output, force=False, timeout=GATE_TIMEOUT):
    output.mkdir(parents=True, exist_ok=True)
    wanted = digest(encoded(file_inventory(source)))
    records = []
    for directory in sorted(output.iterdir()):
        if directory.is_symlink() or not directory.is_dir():
            raise ValueError('invalid verification inventory')
        # Codex reserves these directories at writable roots. They are not receipts.
        if directory.name in ('.git', '.agents', '.codex'):
            continue
        result = verified_result(directory)
        records.append((directory.name, result))
    current = [(name, result) for name, result in records if result['input_sha256'] == wanted]
    if current and not force:
        name, result = current[-1]
    else:
        destination = Path(tempfile.mkdtemp(prefix='gate-', dir=output))
        destination.rmdir()
        result = verify(source, destination, timeout=timeout)
        name = destination.name
    return {'record': name, 'status': result['status'], 'input_sha256': result['input_sha256'],
            'elapsed_ns': result['elapsed_ns']}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--source', required=True)
    parser.add_argument('--out', required=True)
    args = parser.parse_args()
    result = ensure_verified(Path(args.source), Path(args.out))
    print(json.dumps(result))
    record = Path(args.out) / result['record']
    print((record / 'stdout').read_text(), end='')
    print((record / 'stderr').read_text(), end='')
    raise SystemExit(0 if result['status'] == 'completed' else 1)


if __name__ == '__main__':
    main()
