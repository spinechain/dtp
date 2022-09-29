echo "Building SpineChain..."

set /p VERSION=<VERSION

set GOOS=windows
go build -o package/windows/spinechain-%VERSION%.exe -ldflags "-X main.version=%VERSION%" .

:: gh release delete v%VERSION% -y
gh release create v%VERSION% --title "v%VERSION%" --generate-notes
gh release upload v%VERSION% .\package\windows\spinechain-%VERSION%.exe#"spinechain-%VERSION%.exe - Windows 64bit"