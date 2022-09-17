#!/bin/bash

# Welcome
echo "Linux: Welcome to the Stable Diffusion tasktype!"

# Get the first script parameter
SD_PATH=$1
THE_PROMPT=$2
OUTPUT_DIR=$3
OUT_FILE="$OUTPUT_DIR/output.txt"
ERROR_FILE="$OUTPUT_DIR/error.txt"

echo "Diffusing: $THE_PROMPT to $OUTPUT_DIR using $SD_PATH"


# Go to stable diffusion directory
cd $SD_PATH

# Activate ldm conda environment
source activate ldm

echo "Running stable diffusion..."

# Run stable diffusion
# result=$(python scripts/txt2img.py --prompt "$THE_PROMPT" --plms --ckpt /home/mark/stable-diffusion/sd-v1-4.ckpt --skip_grid --n_samples 1 --outdir "$OUTPUT_DIR")
python scripts/txt2img.py --prompt "$THE_PROMPT" --plms --ckpt /home/mark/stable-diffusion/sd-v1-4.ckpt --skip_grid --n_samples 1 --outdir "$OUTPUT_DIR" 2>$ERROR_FILE 1>$OUT_FILE

echo "done"
# echo $result
exit 0
