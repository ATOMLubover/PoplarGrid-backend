#!/bin/bash

# Process blank path.
ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
cd "$ROOT_DIR" || exit 1

BIN_DIR="./bin"
mkdir -p "$BIN_DIR"

GO_OUT="$BIN_DIR/gateway"
GO_SRC="./cmd/gateway"

# Resolve the action.
ACTION="all"
if (($# > 0)); then
    case $1 in
        build) ACTION="build" ;;
        run)   ACTION="run"   ;;
        *)     
            echo "Usage: $0 [build|run]"
            exit 1
        ;;
    esac
fi

# Check Go environment.
if ! command -v go &> /dev/null; then
    echo "ERROR: 'go' command not found. Please install Go first."
    exit 1
fi

case $ACTION in
    build)
        echo "Building binary..."
        if go build -o "$GO_OUT" "$GO_SRC"; then
            echo "Build successful. Binary: $GO_OUT"
        else
            echo "Build failed"
            exit 1
        fi
        ;;

    run)
        if [[ ! -f "$GO_OUT" ]]; then
            echo "Binary not found. Building first..."
            go build -o "$GO_OUT" "$GO_SRC" || exit 1
        fi
        echo "Running $GO_OUT"
        "$GO_OUT"
        ;;

    all)
        echo "Building and running..."
        go build -o "$GO_OUT" "$GO_SRC" || exit 1
        "$GO_OUT"
        ;;
esac