#!/usr/bin/env bash

if [[ -z ${CONSUL_PORT_8500_TCP_ADDR} ]] || [[ -z ${CONSUL_PORT_8500_TCP_PORT} ]]
then
    echo "You need to specify a link to the consul container"
fi

RUN_OPTS="-dataDir /data"
RUN_OPTS="$RUN_OPTS -host ${CONSUL_PORT_8500_TCP_ADDR}"
RUN_OPTS="$RUN_OPTS -port ${CONSUL_PORT_8500_TCP_PORT}"
RUN_OPTS="$RUN_OPTS -repo https://github.com/flyhard/testData.git"

/gitProperties2Consul ${RUN_OPTS}