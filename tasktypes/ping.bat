@echo off
 
echo "Win32: Welcome to thePing tasktype!"

set TARGET=%2
set OUTPUT_DIR=%3
set OUT_FILE="%OUTPUT_DIR%/output.txt"

echo "Pinging: %TARGET% to %OUTPUT_DIR%"

ping google.com -n 4 >> %OUT_FILE%
