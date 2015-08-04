#!/bin/sh

cd `dirname $0` && go build && docker build -t raftis/dashboard .
docker push raftis/dashboard
echo "Published!"
