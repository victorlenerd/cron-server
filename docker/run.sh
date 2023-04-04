#!/usr/bin/env bash

#docker build -t scheduler0:latest --no-cache .
docker run \
    -p 9090:9090 \
    -p 7070:7070 \
    --name scheduler0 \
    -e SCHEDULER0_SECRET_KEY=AB551DED82B93DC8035D624A625920E2121367C7538C02277D2D4DB3C0BFFE94 \
    -e SCHEDULER0_AUTH_PASSWORD=admin \
    -e SCHEDULER0_AUTH_USERNAME=admin \
    -e SCHEDULER0_PROTOCOL=http \
    -e SCHEDULER0_HOST=127.0.0.1 \
    -e SCHEDULER0_PORT=9090 \
    -e SCHEDULER0_REPLICAS="" \
    -e SCHEDULER0_PEER_AUTH_REQUEST_TIMEOUT_MS=2 \
    -e SCHEDULER0_PEER_CONNECT_RETRY_MAX=2 \
    -e SCHEDULER0_PEER_CONNECT_RETRY_DELAY_SECONDS=2 \
    -e SCHEDULER0_BOOTSTRAP=true \
    -e SCHEDULER0_NODE_ID=1 \
    -e SCHEDULER0_RAFT_ADDRESS=127.0.0.1:7070 \
    -e SCHEDULER0_RAFT_TRANSPORT_MAX_POOL=1 \
    -e SCHEDULER0_JOB_EXECUTION_RETRY_DELAY=1 \
    -e SCHEDULER0_JOB_EXECUTION_RETRY_MAX=5 \
    -e SCHEDULER0_MAX_WORKERS=5 \
    -e SCHEDULER0_JOB_QUEUE_DEBOUNCE_DELAY=2 \
    -e SCHEDULER0_MAX_MEMORY=5000 \
    -e SCHEDULER0_EXECUTION_LOG_FETCH_FAN_IN=2 \
    -e SCHEDULER0_EXECUTION_LOG_FETCH_INTERVAL_SECONDS=2 \
    -e SCHEDULER0_JOB_INVOCATION_DEBOUNCE_DELAY=2 \
    -e SCHEDULER0_HTTP_EXECUTOR_PAYLOAD_MAX_SIZE_MB=2 \
    -d scheduler0:latest