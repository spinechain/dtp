#!/bin/bash

# Welcome
echo "Linux: Welcome to the Ping tasktype!"

# Get the first script parameter
TARGET=$1
OUTPUT_DIR=$2
OUT_FILE="$OUTPUT_DIR/output.txt"

echo "Pinging: $TARGET to $OUTPUT_DIR"

ping $TARGET -c 4 > $OUT_FILE

echo "done"