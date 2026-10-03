"""Capture bounded local processes, including interrupted attempts."""

import os
from contextlib import suppress
import signal
import subprocess
import time
from datetime import datetime, timezone


def now():
    return datetime.now(timezone.utc).isoformat()


def execute(command, cwd, output, timeout, environment, body=None):
    started, tick = now(), time.monotonic_ns()
    result = {'started_at': started, 'command': command, 'status': 'failed'}
    child = None
    previous_term = signal.getsignal(signal.SIGTERM)
    def interrupted(signum, frame):
        raise KeyboardInterrupt
    signal.signal(signal.SIGTERM, interrupted)
    with (output / 'stdout').open('wb') as stdout, (output / 'stderr').open('wb') as stderr:
        try:
            child = subprocess.Popen(command, cwd=cwd, stdin=subprocess.PIPE,
                                     stdout=stdout, stderr=stderr, start_new_session=True, env=environment)
            child.communicate(body, timeout=timeout)
            result['exit_code'] = child.returncode
            if child.returncode == 0:
                result['status'] = 'completed'
        except subprocess.TimeoutExpired:
            result['status'] = 'timeout'
        except KeyboardInterrupt:
            result['status'] = 'interrupted'
        except OSError as error:
            result['error'] = str(error)
        finally:
            if child is not None and child.poll() is None:
                os.killpg(child.pid, signal.SIGTERM)
                try:
                    child.wait(timeout=2)
                except subprocess.TimeoutExpired:
                    os.killpg(child.pid, signal.SIGKILL)
                    child.wait()
            signal.signal(signal.SIGTERM, previous_term)
            if child is not None and child.stdin is not None:
                with suppress(OSError):
                    child.stdin.close()
    result.update(ended_at=now(), elapsed_ns=time.monotonic_ns() - tick,
                  clock='monotonic_ns', timing_scope='process launch through exit and output capture')
    return result
