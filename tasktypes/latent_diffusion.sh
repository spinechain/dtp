#!/bin/bash

# Welcome
echo "Linux: Welcome to the Latent Diffusion tasktype!"

# Get the first script parameter
THE_PROMPT=$1

LD_PATH="/home/mark/stable-diffusion"
OUTPUT_DIR="/home/mark/spinechain.dtp/tasktypes/output"

# Go to latent diffusion directory
cd $LD_PATH

# Activate ldm conda environment
source activate ldm

# Run latent diffusion and get return value
result=$(python scripts/txt2img.py --prompt "$THE_PROMPT" --plms --ckpt sd-v1-4.ckpt --skip_grid --n_samples 1 --outdir $OUTPUT_DIR)
echo $result
