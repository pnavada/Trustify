#!/bin/bash

partition_network() {
    local node="$1"
    local network="$2"
    echo "Partitioning $node from $network"
    docker network disconnect "$network" "$node"
}

restore_network() {
    local node="$1"
    local network="$2"
    echo "Restoring $node to $network"
    docker network connect "$network" "$node"
}

simulate_partition() {
    local node="$1"
    local network="$2"
    local duration="${3:-30}"

    partition_network "$node" "$network"
    echo "Network partition applied for $node for $duration seconds"

    sleep "$duration"

    restore_network "$node" "$network"
    echo "Network connectivity restored for $node"
}

# Check arguments
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <node_name> <network1> <network2> [duration]"
    echo "Example: $0 node5 network1 network2 30"
    exit 1
fi

simulate_partition "$@"