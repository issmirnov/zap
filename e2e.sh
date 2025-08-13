#!/bin/bash
set -euo pipefail
PORT=16000
# sanity check for build env
if [[ ! -e zap ]]; then
    echo "no zap binary present"
    exit 1
fi

cleanup() {
	if [[ -n "${ZAP_PID:-}" ]] && kill -0 "$ZAP_PID" 2>/dev/null; then
		kill -9 "$ZAP_PID" || true
	fi
}
trap cleanup EXIT

# start zap, fork.
./zap --port $PORT 1>/dev/null 2>/dev/null &
ZAP_PID=$!

# wait for port to open (with timeout)
ATTEMPTS=0
MAX_ATTEMPTS=30
while ! nc -z localhost $PORT </dev/null 2>/dev/null; do
	sleep 1
	ATTEMPTS=$((ATTEMPTS+1))
	if [[ $ATTEMPTS -ge $MAX_ATTEMPTS ]]; then
		echo "zap did not start listening on port $PORT within timeout"
		exit 3
	fi
done

# pull results.
RESP="$(curl -s -o /dev/null -w '%{http_code} %{redirect_url}\n' -H 'Host: g' localhost:$PORT/z)"

# Check response
expected="302 https://github.com/issmirnov/zap"
if [[ $RESP != $expected ]];then
    echo "Status code or location don't match expectations"
    echo "expected: $expected"
    echo "got: $RESP"
    exit 2
fi

echo "End to end test passed."
