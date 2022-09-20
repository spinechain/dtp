@echo off
 
echo "Win32: Welcome to the Latent Diffusion tasktype!"

set SD_PATH=%1
set THE_PROMPT=%2
set OUTPUT_DIR=%3
set OUT_FILE="%OUTPUT_DIR%/samples/output.txt"

echo "Diffusing: %THE_PROMPT% to %OUTPUT_DIR% using %SD_PATH%"

if not exist "%OUTPUT_DIR%/samples/" mkdir "%OUTPUT_DIR%/samples/"

ping google.com -n 4 >> %OUT_FILE%
