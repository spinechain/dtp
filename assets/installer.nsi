
!define APPNAME "SpineChain"
!define COMPANYNAME "SpineChain"
!define DESCRIPTION "A journey to the edge of knowledge"
!define PRODUCT_WEB_SITE "https://spinecha.in"


# These three must be integers
!define VERSIONMAJOR 0
!define VERSIONMINOR 0
!define VERSIONBUILD 2

!define INSTALLER_OUTPUT_FILE "spinechain-install-${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}.exe"

# These will be displayed by the "Click here for support information" link in "Add/Remove Programs"
# It is possible to use "mailto:" links in here to open the email client
!define HELPURL ${PRODUCT_WEB_SITE} # "Support Information" link
!define UPDATEURL ${PRODUCT_WEB_SITE} # "Product Updates" link
!define ABOUTURL ${PRODUCT_WEB_SITE} # "Publisher" link

# This is the size (in kB) of all the files copied into "Program Files"
!define INSTALLSIZE 7233
 
RequestExecutionLevel admin ;Require admin rights on NT6+ (When UAC is turned on)
 
InstallDir "$PROGRAMFILES\${COMPANYNAME}\${APPNAME}"
 
SetCompressor /SOLID lzma

# rtf or txt file - remember if it is txt, it must be in the DOS text format (\r\n)
LicenseData "..\LICENSE"

CRCCheck On

# This will be in the installer/uninstaller's title bar
Name "${COMPANYNAME} - ${APPNAME}"
Icon "logo.ico"
outFile "${INSTALLER_OUTPUT_FILE}"

Unicode true

!include LogicLib.nsh
 
# Just three pages - license agreement, install location, and installation
page license
page directory
Page instfiles
 
!macro VerifyUserIsAdmin
UserInfo::GetAccountType
pop $0
${If} $0 != "admin" ;Require admin rights on NT4+
        messageBox mb_iconstop "Administrator rights required!"
        setErrorLevel 740 ;ERROR_ELEVATION_REQUIRED
        quit
${EndIf}
!macroend


; MUI 2.0 compatible install
!include "MUI2.nsh"
!include "InstallOptions.nsh"


; Pages to show during installation
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\LICENSE"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES

;!define MUI_FINISHPAGE_RUN "$INSTDIR\gtk2-runtime\gtk2_prefs.exe"
;!define MUI_FINISHPAGE_SHOWREADME "$INSTDIR\Example.file"
;!define MUI_FINISHPAGE_RUN_NOTCHECKED
!define MUI_FINISHPAGE_NOAUTOCLOSE
;!define MUI_FINISHPAGE_NOREBOOTSUPPORT
!insertmacro MUI_PAGE_FINISH



; Uninstaller page
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

function .onInit
	setShellVarContext all
	!insertmacro VerifyUserIsAdmin
functionEnd
 
section "install"

    SetShellVarContext all  ; use all user variables as opposed to current user
	SetOverwrite On

	# Files for the install directory - to build the installer, these should be in the same directory as the install script (this file)
	setOutPath "$INSTDIR\bin"

    File ..\package\windows\bin\libatk-1.0-0.dll		; atk
	File ..\package\windows\bin\libatkmm-1.6-1.dll		; atk
	File ..\package\windows\bin\libssp-0.dll			; needed by cairo
	File ..\package\windows\bin\libcairo-2.dll			; cairo, needed by gtk
	File ..\package\windows\bin\libcairo-gobject-2.dll	; cairo. Doesn't seem to be required, but since we're distributing cairo...
	File ..\package\windows\bin\libcairo-script-interpreter-2.dll  ; cairo. Doesn't seem to be required, but since we're distributing cairo...
	File ..\package\windows\bin\libcairomm-1.0-1.dll
	File ..\package\windows\bin\libepoxy-0.dll
	File ..\package\windows\bin\libexslt-0.dll
	File ..\package\windows\bin\libffi-7.dll  			; libffi is required by glib2 
	File ..\package\windows\bin\libfontconfig-1.dll	; fontconfig is needed for ft2 pango backend
	File ..\package\windows\bin\libfreetype-6.dll		; freetype is needed for ft2 pango backend
	File ..\package\windows\bin\libfribidi-0.dll  ; fribidi is needed for pango 
	File ..\package\windows\bin\libgailutil-3-0.dll	; from gtk
	File ..\package\windows\bin\libgdk_pixbuf-2.0-0.dll  ; from gtk
	File ..\package\windows\bin\liblzma-5.dll  		; from gtk
	File ..\package\windows\bin\libcroco-0.6-3.dll		; from gtk
	File ..\package\windows\bin\libgdk-3-0.dll  		; from gtk
	File ..\package\windows\bin\libgdkmm-3.0-1.dll
	File ..\package\windows\bin\libgio-2.0-0.dll  		; from glib
	File ..\package\windows\bin\libglib-2.0-0.dll  	; glib
	File ..\package\windows\bin\libglibmm-2.4-1.dll  	; glib
	File ..\package\windows\bin\libgiomm-2.4-1.dll  	; glib
	File ..\package\windows\bin\libsigc-2.0-0.dll
	File ..\package\windows\bin\libglibmm_generate_extra_defs-2.4-1.dll  ; glib
	File ..\package\windows\bin\libgmodule-2.0-0.dll	; from glib
	File ..\package\windows\bin\libgobject-2.0-0.dll	; from glib
	File ..\package\windows\bin\libgthread-2.0-0.dll	; from glib
	File ..\package\windows\bin\libgtk-3-0.dll  ; gtk
	File ..\package\windows\bin\libgtksourceview-3.0-1.dll
	File ..\package\windows\bin\libgtksourceview-4-0.dll
	File ..\package\windows\bin\libgtksourceviewmm-3.0-0.dll
	File ..\package\windows\bin\libgtkmm-3.0-1.dll
	File ..\package\windows\bin\libharfbuzz-0.dll 		; required by pango
	File ..\package\windows\bin\libintl-8.dll  		; gettext, needed by all i18n libs
	File ..\package\windows\bin\libiconv-2.dll			; required by fontconfig
	File ..\package\windows\bin\libjson-glib-1.0-0.dll	; gettext, needed by all i18n libs
	File ..\package\windows\bin\libpango-1.0-0.dll  	; pango, needed by gtk
	File ..\package\windows\bin\libpangocairo-1.0-0.dll  ; pango, needed by gtk
	File ..\package\windows\bin\libpangowin32-1.0-0.dll  ; pango, needed by gtk
	File ..\package\windows\bin\libpangoft2-1.0-0.dll	; pango, needed by gtk
	File ..\package\windows\bin\libpangomm-1.4-1.dll
	File ..\package\windows\bin\libpixman-1-0.dll  	; libpixman, needed by cairo
	File ..\package\windows\bin\libpng16-16.dll  		; required by gdk-pixbuf2
	File ..\package\windows\bin\libjpeg-8.dll  		; required by gdk-pixbuf2
	File ..\package\windows\bin\libjasper-4.dll  		; required by gdk-pixbuf2
	File ..\package\windows\bin\libxml++-2.6-2.dll  ; fontconfig needs this
	File ..\package\windows\bin\libxml++-3.0-1.dll
	File ..\package\windows\bin\libxml2-2.dll			; fontconfig needs this
	File ..\package\windows\bin\libxslt-1.dll			; fontconfig needs this
	File ..\package\windows\bin\libpcre-1.dll			; fontconfig needs this
	File ..\package\windows\bin\libthai-0.dll			; fontconfig needs this
	File ..\package\windows\bin\libdatrie-1.dll			; fontconfig needs this
	File ..\package\windows\bin\zlib1.dll				; png and many others need this
	File ..\package\windows\bin\libexpat-1.dll			; required by fontconfig
	File ..\package\windows\bin\libbz2-1.dll			; required by fontconfig
	File ..\package\windows\bin\libgraphite2.dll		; required by harfbuzz
	File ..\package\windows\bin\librsvg-2-2.dll		; required by adwaita-icon-theme
	File ..\package\windows\bin\libtiff-5.dll			; required by gdk-pixbuf2
	File ..\package\windows\bin\libstdc++-6.dll		; standard MSYS2 library
	File ..\package\windows\bin\libgcc_s_seh-1.dll		; standard MSYS2 library
	File ..\package\windows\bin\libwinpthread-1.dll	; standard MSYS2 library
	File ..\package\windows\bin\libsoup-2.4-1.dll      ; libsoup
	File ..\package\windows\bin\libsoup-gnome-2.4-1.dll      ; libsoup
	File ..\package\windows\bin\libsqlite3-0.dll       ; libsoup dependency
	File ..\package\windows\bin\libpsl-5.dll       ; libsoup dependency
	File ..\package\windows\bin\libbrotlidec.dll       ; libsoup dependency
	File ..\package\windows\bin\libbrotlicommon.dll       ; libsoup dependency
	File ..\package\windows\bin\libgnutls-30.dll       ; glib-networking dependency
	File ..\package\windows\bin\libgmp-10.dll		; glib-networking dependency
	File ..\package\windows\bin\libhogweed-6.dll       ; glib-networking dependency
	File ..\package\windows\bin\libnettle-8.dll	; glib-networking dependency
	File ..\package\windows\bin\libidn2-0.dll		; glib-networking dependency
	File ..\package\windows\bin\libp11-kit-0.dll	; glib-networking dependency
	File ..\package\windows\bin\libtasn1-6.dll		; glib-networking dependency
	File ..\package\windows\bin\libunistring-2.dll	; glib-networking dependency
	File ..\package\windows\bin\libproxy-1.dll	; glib-networking dependency
	File ..\package\windows\bin\libpeas-1.0-0.dll	; libpeas
	File ..\package\windows\bin\libpeas-gtk-1.0-0.dll	; libpeas
	File ..\package\windows\bin\libgirepository-1.0-1.dll	; gobject-introspection

	; We install this into the same place as the DLLs to avoid any PATH manipulation.
	SetOutPath "$INSTDIR\bin"
	File ..\package\windows\bin\gdbus.exe
	File ..\package\windows\bin\fc-cache.exe
	File ..\package\windows\bin\fc-cat.exe
	File ..\package\windows\bin\fc-list.exe
	File ..\package\windows\bin\fc-match.exe
	File ..\package\windows\bin\fc-pattern.exe
	File ..\package\windows\bin\fc-query.exe
	File ..\package\windows\bin\fc-scan.exe
	File ..\package\windows\bin\fc-validate.exe
	File ..\package\windows\bin\gdk-pixbuf-query-loaders.exe  ; from gdk_pixbuf
	File ..\package\windows\bin\gspawn-win64-helper.exe
	File ..\package\windows\bin\gspawn-win64-helper-console.exe
	File ..\package\windows\bin\gtk-query-immodules-3.0.exe
	File ..\package\windows\bin\gtk-update-icon-cache.exe


	SetOutPath "$INSTDIR\etc"
	SetOverwrite off
	File /r ..\package\windows\etc\gtk-3.0 
	SetOverwrite On
	File /r ..\package\windows\etc\fonts

	SetOutPath "$INSTDIR\lib\gdk-pixbuf-2.0\2.10.0"
	File ..\package\windows\lib\gdk-pixbuf-2.0\2.10.0\loaders.cache

	SetOutPath "$INSTDIR\lib\gdk-pixbuf-2.0\2.10.0\"
	File /r ..\package\windows\lib\gdk-pixbuf-2.0\2.10.0\loaders

	SetOutPath "$INSTDIR\lib\gio\modules"
	File ..\package\windows\lib\gio\modules\libgiognutls.dll
	File ..\package\windows\lib\gio\modules\libgiognomeproxy.dll
	File ..\package\windows\lib\gio\modules\libgiolibproxy.dll

	SetOutPath "$INSTDIR\lib\libpeas-1.0\loaders"
	File ..\package\windows\lib\libpeas-1.0\loaders\libpython3loader.dll

	SetOutPath "$INSTDIR\lib"
	File /r ..\package\windows\lib\girepository-1.0

	SetOutPath "$INSTDIR\ssl\certs"
	File ..\package\windows\ssl\certs\ca-bundle.crt
	File ..\package\windows\ssl\certs\ca-bundle.trust.crt

	;SetOutPath "$INSTDIR\lib\gtk-3.0\${GTK_BIN_VERSION}"
	; no longer in gtk as of 2.14.5.
	; File /r lib\gtk-2.0\${GTK_BIN_VERSION}\immodules
	; gone as of gtk 2.16.6-2.
	; File /r lib\gtk-2.0\${GTK_BIN_VERSION}\loaders

	; wimp
	; SetOutPath "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}\engines"
	; File lib\gtk-2.0\${GTK_BIN_VERSION}\engines\libwimp*.dll
	; We install this, but other installers may not have it.
	; File lib\gtk-2.0\${GTK_BIN_VERSION}\engines\libpixmap*.dll

	SetOutPath "$INSTDIR\share\locale"
	File ..\package\windows\share\locale\locale.alias  ; from gettext

	SetOutPath "$INSTDIR\share\themes\Emacs"
	File /r ..\package\windows\share\themes\Emacs\gtk-3.0
	SetOutPath "$INSTDIR\share\themes\Default"
	File /r ..\package\windows\share\themes\Default\gtk-3.0

	SetOutPath "$INSTDIR\share\glib-2.0"
	File /r ..\package\windows\share\glib-2.0\schemas

	SetOutPath "$INSTDIR\share"
	File /r ..\package\windows\share\icons

	SetOutPath "$INSTDIR\share"
	File /r ..\package\windows\share\gtksourceview-3.0

	SetOutPath "$INSTDIR\share"
	File /r ..\package\windows\share\gtksourceview-4


	# Files added here should be removed by the uninstaller (see section "uninstall")
	file "..\package\windows\spinechain-0.0.2.exe"
	file "logo.ico"
	# Add any other files for the install directory (license files, app data, etc) here
 
	# Uninstaller - See function un.onInit and section "uninstall" for configuration
	writeUninstaller "$INSTDIR\uninstall.exe"
 
	# Start Menu
	createDirectory "$SMPROGRAMS\${COMPANYNAME}"
	createShortCut "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk" "$INSTDIR\spinechain.exe" "" "$INSTDIR\logo.ico"
 
	# Registry information for add/remove programs
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "DisplayName" "${COMPANYNAME} - ${APPNAME} - ${DESCRIPTION}"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "QuietUninstallString" "$\"$INSTDIR\uninstall.exe$\" /S"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "InstallLocation" "$\"$INSTDIR$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "DisplayIcon" "$\"$INSTDIR\logo.ico$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "Publisher" "$\"${COMPANYNAME}$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "HelpLink" "$\"${HELPURL}$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "URLUpdateInfo" "$\"${UPDATEURL}$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "URLInfoAbout" "$\"${ABOUTURL}$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "DisplayVersion" "$\"${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}$\""
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "VersionMajor" ${VERSIONMAJOR}
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "VersionMinor" ${VERSIONMINOR}
	# There is no option for modifying or repairing the install
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "NoModify" 1
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "NoRepair" 1
	# Set the INSTALLSIZE constant (!defined at the top of this script) so Add/Remove Programs can accurately report the size
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}" "EstimatedSize" ${INSTALLSIZE}
sectionEnd
 
# Uninstaller
 
function un.onInit
	SetShellVarContext all
 
	#Verify the uninstaller - last chance to back out
	MessageBox MB_OKCANCEL "Permanantly remove ${APPNAME}?" IDOK next
		Abort
	next:
	!insertmacro VerifyUserIsAdmin
functionEnd
 
section "uninstall"
 
	# Remove Start Menu launcher
	delete "$SMPROGRAMS\${COMPANYNAME}\${APPNAME}.lnk"
	# Try to remove the Start Menu folder - this will only happen if it is empty
	rmDir "$SMPROGRAMS\${COMPANYNAME}"
 
	# Remove files
	delete $INSTDIR\spinechain.exe
	delete $INSTDIR\logo.ico
 
	# Always delete uninstaller as the last action
	delete $INSTDIR\uninstall.exe
 
	# Try to remove the install directory - this will only happen if it is empty
	rmDir $INSTDIR
 
	# Remove uninstaller information from the registry
	DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${COMPANYNAME} ${APPNAME}"
sectionEnd