"""Post the exact input bytes to the hosted Jev evaluation endpoint."""

import os
import json
import sys
import urllib.error
import urllib.request

from evidence import require_hosted_opt_in


def main():
    if os.environ.get('BENCH_OFFLINE'):
        raise SystemExit('BENCH_OFFLINE suppresses hosted Jev')
    root = os.environ.get('JEV_BENCHMARK_FROZEN_ROOT')
    if not root:
        raise SystemExit('Validated benchmark source is required')
    require_hosted_opt_in(root)
    key = os.environ.get('TYPESAFE_API_KEY')
    if not key:
        raise SystemExit('TYPESAFE_API_KEY is absent')
    body = sys.stdin.buffer.read()
    request = urllib.request.Request('https://api.typesafe.ai/v1/systemone', data=body,
                                     headers={'Content-Type': 'application/json',
                                              'Authorization': 'Bearer ' + key}, method='POST')
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            sys.stdout.buffer.write(response.read())
    except urllib.error.HTTPError as error:
        sys.stderr.write('Jev HTTP status: %d\n' % error.code)
        sys.stdout.buffer.write(json.dumps({'http_status': error.code,
                                           'provider_body': error.read().decode(errors='replace')}).encode())
        raise SystemExit(1)


if __name__ == '__main__':
    main()
