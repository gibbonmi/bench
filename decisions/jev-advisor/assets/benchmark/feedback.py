"""Request fixed-command gate feedback from the task controller."""

import argparse
import os
import re
import stat
import subprocess
import time
import uuid
from pathlib import Path

from evidence import digest, encoded, file_inventory, inventory_changes, read
from verification import GATE_TIMEOUT, ensure_verified


def publish(path, value):
    temporary = path.with_name(uuid.uuid4().hex + '.tmp')
    with temporary.open('xb') as stream:
        stream.write(encoded(value))
    temporary.replace(path)


class Feedback:
    def __init__(self, source, baseline, allowed, mailbox, receipts):
        self.source, self.baseline = source, file_inventory(baseline)
        self.allowed = set(allowed)
        self.mailbox, self.receipts = mailbox, receipts
        mailbox.mkdir()
        self.seen = set()

    def __call__(self, remaining):
        for path in sorted(self.mailbox.glob('*.request')):
            if path.name in self.seen:
                continue
            if not re.fullmatch(r'[0-9a-f]{32}\.request', path.name):
                raise ValueError('invalid feedback request name')
            self.seen.add(path.name)
            with os.fdopen(os.open(path, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK), 'rb') as stream:
                if not stat.S_ISREG(os.fstat(stream.fileno()).st_mode):
                    raise ValueError('feedback request is not a regular file')
                body = stream.read(1025)
            inventory = file_inventory(self.source)
            expected = encoded({'input_sha256': digest(encoded(inventory))})
            changed = set(inventory_changes(self.baseline, inventory))
            if body != expected:
                response = {'status': 'refused', 'error': 'request differs from current task inputs'}
            elif changed - self.allowed:
                response = {'status': 'refused', 'error': 'task writes exceed the frozen fence'}
            elif remaining <= 0:
                response = {'status': 'timeout'}
            else:
                try:
                    response = ensure_verified(self.source, self.receipts, timeout=min(remaining, GATE_TIMEOUT))
                    record = self.receipts / response['record']
                    response.update(stdout=(record / 'stdout').read_text(errors='replace'),
                                    stderr=(record / 'stderr').read_text(errors='replace'))
                except (OSError, ValueError, subprocess.SubprocessError) as error:
                    response = {'status': 'failed', 'error': str(error)}
            publish(path.with_suffix('.response'), response)
            if response['status'] == 'interrupted':
                raise KeyboardInterrupt
            return


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--source', type=Path, required=True)
    parser.add_argument('--mailbox', type=Path, required=True)
    parser.add_argument('--timeout', type=float, required=True)
    args = parser.parse_args()
    token = uuid.uuid4().hex
    request = args.mailbox / (token + '.request')
    publish(request, {'input_sha256': digest(encoded(file_inventory(args.source)))})
    response_path = request.with_suffix('.response')
    deadline = time.monotonic() + args.timeout
    while not response_path.exists():
        if time.monotonic() >= deadline:
            raise SystemExit('controller feedback timed out')
        time.sleep(0.05)
    response = read(response_path)
    print(encoded({k: v for k, v in response.items() if k not in ('stdout', 'stderr')}).decode())
    print(response.get('stdout', ''), end='')
    print(response.get('stderr', ''), end='')
    raise SystemExit(0 if response['status'] == 'completed' else 1)


if __name__ == '__main__':
    main()
