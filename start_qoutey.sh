#!/bin/bash

# Script to start qoutey in a tmux session
# Usage: ./start_qoutey.sh [test]

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo "tmux is not installed. Please install it with your package manager."
    echo "For example: sudo apt install tmux"
    exit 1
fi

# Set the directory where qoutey is located
QOUTEY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$QOUTEY_DIR"

# Check if qoutey exists
if [ ! -f "./qoutey" ]; then
    echo "qoutey executable not found. Building it now..."
    go build -o qoutey ./cmd/qoutey
    if [ $? -ne 0 ]; then
        echo "Failed to build qoutey. Please check errors above."
        exit 1
    fi
fi

# Check if the config file exists
if [ ! -f "./config.json" ]; then
    echo "config.json not found. Please make sure to configure it."
    echo "Running qoutey once to generate a default config..."
    ./qoutey test
    echo "Please edit config.json with your SMTP settings and quotes, then run this script again."
    exit 1
fi

# Define tmux session name
SESSION_NAME="qoutey"

# Check if the session already exists
tmux has-session -t $SESSION_NAME 2>/dev/null

if [ $? -eq 0 ]; then
    echo "Session $SESSION_NAME already exists."
    echo "To attach to it, run: tmux attach -t $SESSION_NAME"
    echo "To kill it first, run: tmux kill-session -t $SESSION_NAME"
    exit 0
fi

# If the first argument is "test", run in test mode
if [ "$1" = "test" ]; then
    echo "Running qoutey in test mode..."
    ./qoutey test
    exit 0
fi

# Create new tmux session and start qoutey
echo "Starting qoutey in tmux session named '$SESSION_NAME'..."
tmux new-session -d -s $SESSION_NAME "./qoutey"

echo "qoutey is now running in a tmux session."
echo "To attach to the session: tmux attach -t $SESSION_NAME"
echo "To detach from the session (when attached): CTRL+B then D"
echo "To stop the application: tmux kill-session -t $SESSION_NAME"