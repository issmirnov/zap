#!/bin/bash

# sanity check for build env
if [[ ! -e zap ]]; then
    echo "no zap binary present"
    exit 1
fi

# start zap, fork.
./zap --port 16000 1>/dev/null 2>/dev/null &
ZAP_PID=$!

# pull results.
RESP="$(curl -I -L -H 'Host: g' localhost:16000/z 2>/dev/null | head -n 2)"

# Check header
if [[ $RESP != *"HTTP/1.1 302 Found"* ]]; then
    echo "302 status not found"
    echo "got: $RESP"
    kill -9 $ZAP_PID
    exit 2
fi

# check location
if [[ $RESP != *"Location: https://github.com/issmirnov/zap"* ]]; then
    echo "Location is wrong"
    echo "Got: $RESP"
    kill -9 $ZAP_PID
    exit 3
fi

kill -9 $ZAP_PID
echo "End to end test passed."
