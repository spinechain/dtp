VERSION=$(<VERSION)

gh release upload v%VERSION% ..\package\windows\spinechain-installer-%VERSION%.exe#"spinechain-%VERSION%.exe - Windows Installer 64bit"