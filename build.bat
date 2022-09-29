echo "Building SpineChain..."

set VERSION=0.0.1
set GOOS=windows
go build -o package/windows/spinechain-%VERSION%.exe -ldflags "-X main.version=%VERSION%" .

gh release delete v%VERSION%
gh release create v%VERSION% --title "v%VERSION%" --generate-notes
gh release upload v%VERSION% .\package\windows\spinechain-%VERSION%.exe