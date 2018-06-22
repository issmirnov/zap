#!/bin/bash
PORT=16000
# sanity check for build env
if [[ ! -e zap ]]; then
    echo "no zap binary present"
    exit 1
fi

# start zap, fork.
./zap --port $PORT 1>/dev/null 2>/dev/null &
ZAP_PID=$!

# wait for port to open
while ! nc -z localhost $PORT </dev/null 2>/dev/null; do sleep 1; done

# pull results.
RESP="$(curl -s -o /dev/null -w '%{http_code} %{redirect_url}\n' -H 'Host: g' localhost:$PORT/z)"

# Check response
expected="302 https://github.com/issmirnov/zap"
if [[ $RESP != $expected ]];then
    echo "Status code or location don't match expectations"
    echo "expected: $expected"
    echo "got: $RESP"
    kill -9 $ZAP_PID
    exit 2
fi

# cleanup
kill -9 $ZAP_PID
echo "End to end test passed."
