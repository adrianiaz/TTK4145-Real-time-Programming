#!/bin/bash

while true; do
    echo "go run main.go -name dune2 -port 15657"
    go run main.go -name dune2 -port 15657

    echo "Go program exited. Restarting..."
    sleep 0.5
done