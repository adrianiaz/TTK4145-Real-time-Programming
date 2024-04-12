#!/bin/bash

while true; do
    echo "go run main.go -name dune3 -port 15659"
    go run main.go -name dune3 -port 15659

    echo "Go program exited. Restarting..."
    sleep 0.5
done