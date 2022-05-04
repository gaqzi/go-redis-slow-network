#!/bin/bash

# Each test will restart docker-compose and the running server
docker-compose up --force-recreate -d slow_redis
# Add 100ms of delay to outgoing
docker-compose exec slow_redis tc qdisc add dev eth0 root netem delay 100ms

docker-compose up --force-recreate client

docker-compose exec slow_redis redis-cli INFO
docker-compose down

echo "Look at total_connections_received:!"
echo "DONE"
