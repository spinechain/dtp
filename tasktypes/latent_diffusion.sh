#!/bin/bash

function getcpid() {
    cpids=`pgrep -P $1|xargs`
#    echo "cpids=$cpids"
    for cpid in $cpids;
    do
        echo "$cpid"
        getcpid $cpid
    done
}

# Welcome
echo "Linux: Welcome to the Latent Diffusion tasktype!"

# Get the first script parameter
THE_PROMPT=$1
OUTPUT_DIR=$2

echo "Diffusing: $THE_PROMPT to $OUTPUT_DIR"

LD_PATH="/home/mark/stable-diffusion"


# Go to latent diffusion directory
cd $LD_PATH

# Activate ldm conda environment
source activate ldm

# Run latent diffusion and get return value
# result=$(python scripts/txt2img.py --prompt "$THE_PROMPT" --plms --ckpt sd-v1-4.ckpt --skip_grid --n_samples 1 --outdir "$OUTPUT_DIR")

python scripts/txt2img.py --prompt "$THE_PROMPT" --plms --ckpt /home/mark/stable-diffusion/sd-v1-4.ckpt --skip_grid --n_samples 1 --outdir "$OUTPUT_DIR" 2>/home/mark/error.txt 1>/home/mark/output.txt
# process_id=$!

# Print process ID
# echo "Waiting for Process ID: $process_id"

# wait $process_id

# echo "Process completed with exit code $?."

echo $result

