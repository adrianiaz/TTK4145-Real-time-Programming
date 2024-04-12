#!/bin/bash

while true; do
    echo "go run main.go -name dune2 -port 15658"
    go run main.go -name dune2 -port 15658

    echo "Go program exited. Restarting..."
    sleep 0.5
done