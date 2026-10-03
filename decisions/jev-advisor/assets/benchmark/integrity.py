"""Bind captured run evidence to an append-only set of trial digests."""

import os
from pathlib import Path

from evidence import digest, encoded, file_inventory, read, save


def start_ledger(root):
    ledger = {'version': 2, 'plan_sha256': digest((root / 'plan.json').read_bytes()),
              'preflight_sha256': digest(encoded(file_inventory(root / 'preflight'))), 'trials': {}}
    save(root / 'ledger.json', ledger)
    return ledger


def capture_trial(root, name, ledger):
    if name in ledger['trials']:
        raise ValueError('trial already captured')
    ledger['trials'][name] = digest(encoded(file_inventory(root / name)))
    pending = root / 'ledger.pending'
    save(pending, ledger)
    os.replace(pending, root / 'ledger.json')


def verified_run(root):
    root = Path(root)
    for name in ('ledger.json', 'plan.json'):
        if (root / name).is_symlink() or not (root / name).is_file():
            raise ValueError('missing or invalid run integrity evidence: ' + name)
    ledger = read(root / 'ledger.json')
    if ledger.get('version') != 2 or ledger.get('plan_sha256') != digest((root / 'plan.json').read_bytes()):
        raise ValueError('plan changed after capture')
    if (not (root / 'preflight').is_dir() or (root / 'preflight').is_symlink()
            or ledger.get('preflight_sha256') != digest(encoded(file_inventory(root / 'preflight')))):
        raise ValueError('baseline preflight changed after capture')
    plan = read(root / 'plan.json')
    if plan.get('config_sha256') != digest(encoded(plan['config'])):
        raise ValueError('plan configuration fingerprint mismatch')
    captured = ledger['trials']
    observed = {path.name for path in root.glob('trial-*')}
    if observed != set(captured):
        raise ValueError('trial inventory changed or contains unsealed interrupted evidence')
    schedule = {'trial-%03d' % index: row for index, row in enumerate(plan['schedule'])}
    if not set(captured).issubset(schedule):
        raise ValueError('captured trial is outside the schedule')
    trials = []
    for name in sorted(captured):
        directory = root / name
        if (directory.is_symlink() or not directory.is_dir()
                or captured[name] != digest(encoded(file_inventory(directory)))):
            raise ValueError('trial evidence changed after capture: ' + name)
        path = directory / 'trial.json'
        trial = read(path)
        if any(trial.get(key) != value for key, value in schedule[name].items()):
            raise ValueError('trial identity differs from frozen schedule')
        trials.append((path, trial))
    return plan, trials
