
!define PRODUCT_VERSION "0.0.3"
!define PRODUCT_NAME "SpineChain"
!define PRODUCT_PUBLISHER "Mark Essien"
!define PRODUCT_WEB_SITE "https://spinecha.in"
!define INSTALLER_OUTPUT_FILE "spinechain-installer-${PRODUCT_VERSION}.exe"

;!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\AppMainExe.exe"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"

!define REGISTRY_APP_PATHS "Software\Microsoft\Windows\CurrentVersion\App Paths"


; AddToPath and friends should work with all users
!define ALL_USERS

!include "FileFunc.nsh"  ; GetOptions
!include "x64.nsh"
!include "LogicLib.nsh"


; --------------- General Settings


; this is needed for proper start menu item manipulation (for all users) in vista
RequestExecutionLevel admin

; This compressor gives us the best results
SetCompressor /SOLID lzma

; Do a CRC check before installing
CRCCheck On

; This is used in titles
Name "${PRODUCT_NAME}"  ;  ${PRODUCT_VERSION}

; Output File Name
OutFile "${INSTALLER_OUTPUT_FILE}"

Unicode true

; The Default Installation Directory
InstallDir "$PROGRAMFILES64\${PRODUCT_NAME}"
;InstallDir "$WINDIR"
; Detect the old installation
InstallDirRegKey HKLM "SOFTWARE\${PRODUCT_NAME}" ""
;InstallDirRegKey HKLM "${PRODUCT_DIR_REGKEY}" ""

ShowInstDetails show
ShowUnInstDetails show





; --------------------- MUI INTERFACE

; MUI 2.0 compatible install
!include "MUI2.nsh"
!include "InstallOptions.nsh"

Icon "logo.ico"
UninstallIcon "logo.ico"

; MUI Settings
!define MUI_ABORTWARNING
!define MUI_ICON "logo.ico"
;!define MUI_UNICON "nsi_uninstall.ico"
!define MUI_UNICON "logo.ico"


; Things that need to be extracted on first (keep these lines before any File command!)
; Only useful for BZIP2 compression
; ReserveFile "nsi_pathpage.ini"
; ReserveFile "nsi_configpage.ini"
; ReserveFile "${NSISDIR}\Plugins\x86-unicode\InstallOptions.dll"


; Pages to show during installation
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\LICENSE"
!insertmacro MUI_PAGE_COMPONENTS
!define MUI_PAGE_CUSTOMFUNCTION_LEAVE DirectoryPageExit
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
UninstPage custom un.DeleteConfig  ;Custom page
!insertmacro MUI_UNPAGE_INSTFILES





; Note: These options are unsupported unless the installer is launched in silent mode (/S).
; e.g. /setpath=no /dllpath=root /sideeffects=no
var install_option_setpath  ; set PATH: yes (default), no
var install_option_dllpath  ; bin (default), lib, root
var install_option_sideeffects  ; yes (default), no. no = don't write to registry, PATH or start menu.
var install_option_translations  ; install translations: yes, no (default)
var install_option_removeold  ; uninstall the old version first (if present): yes (default), no.

var LIB_INSTDIR
var DLL_DIR_NAME
var DLL_TMP





; Set $LIB_INSTDIR to <instpath>\bin, <instpath>\lib or <instpath>\ .
; Must be after the directory selection page.
Function DirectoryPageExit
	StrCpy $LIB_INSTDIR "$INSTDIR"
	StrCmp $DLL_DIR_NAME "" no_dll_append
		StrCpy $LIB_INSTDIR "$INSTDIR\$DLL_DIR_NAME"
	no_dll_append:
FunctionEnd



Function un.DeleteConfig
	; !insertmacro MUI_HEADER_TEXT "$(TEXT_IO_TITLE)" "$(TEXT_IO_SUBTITLE)"
	!insertmacro INSTALLOPTIONS_DISPLAY "nsi_configpage.ini"
FunctionEnd



; Language files
!insertmacro MUI_LANGUAGE "English"


; --------------- END MUI



;Description
; LangString DESC_SecCopyUI ${LANG_ENGLISH} "GTK3 Runtime 64-bit"
; LangString TEXT_IO_TITLE ${LANG_ENGLISH} "GTK3 Runtime 64-bit"
; LangString TEXT_IO_SUBTITLE ${LANG_ENGLISH} "Additional options"


;License page Introduction
;LicenseText "You must agree to this license before installing."
;License text
;LicenseData /LANG=${LANG_ENGLISH} "license.txt"





; ----------------- INSTALLATION TYPES

InstType "Recommended"  ; 1
InstType "Full"  ; 2


var SEC_TRANSLATIONS_INSTALLED


Section "SpineChain Core (required)" SecGTK
SectionIn 1 2 RO
	SetShellVarContext all  ; use all user variables as opposed to current user
	SetOverwrite On

    SetOutPath "$INSTDIR\bin"

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
	SetOutPath "$INSTDIR\bin"
	file "..\package\windows\spinechain-${PRODUCT_VERSION}.exe"
	file "logo.ico"
	# Add any other files for the install directory (license files, app data, etc) here
	
	SetOutPath "$INSTDIR\"
	file "logo.ico"
	# Uninstaller - See function un.onInit and section "uninstall" for configuration
	


	; this script updates some config files, but it's unsafe
	; (gtk or pango may not work afterwards), so don't call it.
	; Push $INSTDIR\gtk3-runtime\gtk-postinstall.bat
	; Call WritePostInstall
	; update pango.modules, not working for now
	; Exec '$INSTDIR\gtk2-runtime\gtk-postinstall.bat'

SectionEnd ; end of gtk section




; Section descriptions
;!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
;	!insertmacro MUI_DESCRIPTION_TEXT ${SecGTK} "GTK+ 64-bit Runtime Libraries"
;	!insertmacro MUI_DESCRIPTION_TEXT ${SecTranslations} "Additional translations (some are incomplete)"
; !insertmacro MUI_FUNCTION_DESCRIPTION_END




; Executed on installer run
Function .onInit
	SetShellVarContext all  ; use all user variables as opposed to current user
	${IfNot} ${RunningX64} 
		MessageBox MB_OK|MB_ICONEXCLAMATION "This installation requires a 64-bit Windows system" /SD IDOK
		Abort
	${EndIf}

	SetRegView 64
		
	

	!insertmacro INSTALLOPTIONS_EXTRACT "nsi_pathpage.ini"

	StrCpy $SEC_TRANSLATIONS_INSTALLED "0"  ; set to 1 in appropriate section

	${GetOptions} "$CMDLINE" "/setpath=" $install_option_setpath
	${GetOptions} "$CMDLINE" "/dllpath=" $install_option_dllpath
	${GetOptions} "$CMDLINE" "/sideeffects=" $install_option_sideeffects
	${GetOptions} "$CMDLINE" "/translations=" $install_option_translations
	${GetOptions} "$CMDLINE" "/removeold=" $install_option_removeold

	; Debug stuff
	; MessageBox MB_ICONINFORMATION|MB_OK "/setpath=$install_option_setpath \
	;	/dllpath=$install_option_dllpath /sideeffects=$install_option_sideeffects \
	;	INSTDIR: $INSTDIR" /SD IDOK

	; if we're using /sideeffects=no, set /setpath=no, because we can't
	; revert it during uninstall (there's no dllpath in registry).
	StrCmp $install_option_sideeffects "no" "" init_sideeffects
		StrCpy $install_option_setpath "no"  ; set /setpath=no
		goto init_sideeffects_exit
	init_sideeffects:
		Call PreventMultipleInstances  ; in no-sideeffects mode this has no purpose
		Call DetectPrevInstallation  ; we don't want local installations to interfere with global ones.
	init_sideeffects_exit:


	; Page callbacks are not called if in silent mode, so call these manually
	IfSilent "" +3
		Call DirectoryPageExit

FunctionEnd



; ------------------ POST INSTALL


var ADD_TO_PATH


Section -post
	SetShellVarContext all  ; use all user variables as opposed to current user

	; write out uninstaller
	WriteUninstaller "$INSTDIR\spinechain_uninst.exe"

	StrCmp $install_option_sideeffects "no" no_sideeffects
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "InstallationDirectory" "$INSTDIR"
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "DllPath" "$LIB_INSTDIR"
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "Vendor" "${PRODUCT_PUBLISHER}"
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "PackageVersion" "${PRODUCT_VERSION}"
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "Version" "${GTK_VERSION}"
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "BinVersion" "${GTK_BIN_VERSION}"

		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "DllDirName" "$DLL_DIR_NAME"  ; lib, bin, or ""
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "UsingSystemPath" $ADD_TO_PATH
		WriteRegStr HKLM "SOFTWARE\${PRODUCT_NAME}" "TranslationsInstalled" $SEC_TRANSLATIONS_INSTALLED

		; compat with installer from http://gimp-win.sourceforge.net/
		; WriteRegStr HKLM "SOFTWARE\GTK\3.0" "Path" "$INSTDIR"
		; WriteRegStr HKLM "SOFTWARE\GTK\3.0" "Version" "${GTK_VERSION}"
		; compat with Dropline's GTK
		; WriteRegStr HKLM "SOFTWARE\GTK\3.0" "DllPath" "$LIB_INSTDIR"

		; Information for the uninstall component in "add/remove programs"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayName" "${PRODUCT_NAME}"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\spinechain_uninst.exe"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "InstallLocation" "$INSTDIR"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\gtk3-runtime\gtk.ico"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
		WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
		WriteRegDWORD HKLM "${PRODUCT_UNINST_KEY}" "NoModify" 1
		WriteRegDWORD HKLM "${PRODUCT_UNINST_KEY}" "NoRepair" 1

		; uninstall shortcut
		CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
		CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\bin\spinechain-${PRODUCT_VERSION}.exe" "" "$INSTDIR\logo.ico"
		CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\Uninstall ${PRODUCT_NAME}.lnk" "$INSTDIR\spinechain_uninst.exe" "" ""
		WriteIniStr "$SMPROGRAMS\${PRODUCT_NAME}\Go to the website.url" "InternetShortcut" "URL" "${PRODUCT_WEB_SITE}"


		; Write $INSTDIR\gtk2-runtime\gtk2r-env.bat
		; This script sets the GTK environment variables
		; DetailPrint "Generating $INSTDIR\gtk3-runtime\gtk3r-env.bat"
		; Push $INSTDIR\gtk3-runtime\gtk3r-env.bat
		; Call WriteEnvBat
		DetailPrint "Done"

	no_sideeffects:

SectionEnd ; post





; ---------------- UNINSTALL


; Note: These options are unsupported unless the uninstaller is launched in silent mode (/S).
var uninstall_option_remove_config  ; yes, no (default).
var uninstall_option_sideeffects  ; yes (default), no. Use if it was installed with this option.
; These are used only if /sideffects=no :
var uninstall_option_dllpath  ; uninstall dlls from: bin (default), lib, root.
var uninstall_option_translations  ; uninstall translations: yes, no (default)


Function un.onInit
	SetRegView 64
	${GetOptions} "$CMDLINE" "/remove_config=" $uninstall_option_remove_config
	${GetOptions} "$CMDLINE" "/sideeffects=" $uninstall_option_sideeffects
	${GetOptions} "$CMDLINE" "/dllpath=" $uninstall_option_dllpath
	${GetOptions} "$CMDLINE" "/translations=" $uninstall_option_translations
FunctionEnd



Function un.onUninstSuccess
	HideWindow
	MessageBox MB_ICONINFORMATION|MB_OK "$(^Name) was successfully removed from your computer." /SD IDOK
FunctionEnd




var leave_config  ; don't remove global gtk config (uninstall page option)
;var INST_DIR_REG  ; installation dir from registry (uninstall stage)


; This is a separate function to allow multiple calls to it
Function un.DeleteDlls
	SetShellVarContext all  ; use all user variables as opposed to current user

	StrCpy $LIB_INSTDIR "$INSTDIR"
	StrCmp $DLL_DIR_NAME "" un_no_dll_append
		StrCpy $LIB_INSTDIR "$INSTDIR\$DLL_DIR_NAME"
	un_no_dll_append:

	; bin stuff (they are in the same directory)
	Delete $LIB_INSTDIR\gdbus.exe
	Delete $LIB_INSTDIR\fc-cache.exe
	Delete $LIB_INSTDIR\fc-cat.exe
	Delete $LIB_INSTDIR\fc-list.exe
	Delete $LIB_INSTDIR\fc-match.exe
	Delete $LIB_INSTDIR\fc-pattern.exe
	Delete $LIB_INSTDIR\fc-query.exe
	Delete $LIB_INSTDIR\fc-scan.exe
	Delete $LIB_INSTDIR\fc-validate.exe
	Delete $LIB_INSTDIR\gdk-pixbuf-query-loaders.exe
	Delete $LIB_INSTDIR\gspawn-win64-helper.exe
	Delete $LIB_INSTDIR\gspawn-win64-helper-console.exe
	Delete $LIB_INSTDIR\gtk-query-immodules-3.0.exe
	Delete $LIB_INSTDIR\gtk-update-icon-cache.exe

	; dlls
	Delete $LIB_INSTDIR\libatk-1.0-0.dll  ; atk
	Delete $LIB_INSTDIR\libatkmm-1.6-1.dll  ; atk
	Delete $LIB_INSTDIR\libcairo-2.dll  ; cairo, needed by gtk
	Delete $LIB_INSTDIR\libcairo-gobject-2.dll  ; cairo. Doesn't seem to be required, but since we're distributing cairo...
	Delete $LIB_INSTDIR\libcairo-script-interpreter-2.dll  ; cairo. Doesn't seem to be required, but since we're distributing cairo...
	Delete $LIB_INSTDIR\libcairomm-1.0-1.dll
	Delete $LIB_INSTDIR\libssp-0.dll
	Delete $LIB_INSTDIR\libepoxy-0.dll
	Delete $LIB_INSTDIR\libexslt-0.dll
	Delete $LIB_INSTDIR\libffi-7.dll  ; libffi is required by glib 
	Delete $LIB_INSTDIR\libfontconfig-1.dll  ; fontconfig is needed for ft2 pango backend
	Delete $LIB_INSTDIR\libfreetype-6.dll  ; freetype is needed for ft2 pango backend
	Delete $LIB_INSTDIR\libfribidi-0.dll
	Delete $LIB_INSTDIR\libgailutil-3-0.dll  ; from gtk
	Delete $LIB_INSTDIR\libgdk_pixbuf-2.0-0.dll  ; from gtk
	Delete $LIB_INSTDIR\liblzma-5.dll  ; from gtk
	Delete $LIB_INSTDIR\libcroco-0.6-3.dll  ; from gtk
	Delete $LIB_INSTDIR\libgdk-3-0.dll  ; from gtk
	Delete $LIB_INSTDIR\libgdkmm-3.0-1.dll
	Delete $LIB_INSTDIR\libgio-2.0-0.dll  ; from glib
	Delete $LIB_INSTDIR\libglib-2.0-0.dll  ; glib
	Delete $LIB_INSTDIR\libglibmm-2.4-1.dll  ; glib
	Delete $LIB_INSTDIR\libgiomm-2.4-1.dll  ; glib
	Delete $LIB_INSTDIR\libsigc-2.0-0.dll  ; glib
	Delete $LIB_INSTDIR\libglibmm_generate_extra_defs-2.4-1.dll  ; glib
	Delete $LIB_INSTDIR\libgmodule-2.0-0.dll  ; from glib
	Delete $LIB_INSTDIR\libgobject-2.0-0.dll  ; from glib
	Delete $LIB_INSTDIR\libgthread-2.0-0.dll  ; from glib
	Delete $LIB_INSTDIR\libgtk-3-0.dll  ; gtk
	Delete $LIB_INSTDIR\libgtksourceview-3.0-1.dll  ; gtk
	Delete $LIB_INSTDIR\libgtksourceview-4-0.dll  ; gtk
	Delete $LIB_INSTDIR\libgtksourceviewmm-3.0-0.dll  ; gtk
	Delete $LIB_INSTDIR\libgtkmm-3.0-1.dll
	Delete $LIB_INSTDIR\libharfbuzz-0.dll
	Delete $LIB_INSTDIR\libintl-8.dll  ; gettext, needed by all i18n libs
	Delete $LIB_INSTDIR\libiconv-2.dll
	Delete $LIB_INSTDIR\libjson-glib-1.0-0.dll 
	Delete $LIB_INSTDIR\libpango-1.0-0.dll  ; pango, needed by gtk
	Delete $LIB_INSTDIR\libpangocairo-1.0-0.dll  ; pango, needed by gtk
	Delete $LIB_INSTDIR\libpangowin32-1.0-0.dll  ; pango, needed by gtk
	Delete $LIB_INSTDIR\libpangoft2-1.0-0.dll  ; pango, needed by gtk
	Delete $LIB_INSTDIR\libpangomm-1.4-1.dll
	Delete $LIB_INSTDIR\libpixman-1-0.dll  ; libpixman, needed by cairo
	Delete $LIB_INSTDIR\libpng16-16.dll  ; for gdk_pixbuf loader.
	Delete $LIB_INSTDIR\libjasper-4.dll  ; for gdk_pixbuf loader.
	Delete $LIB_INSTDIR\libjpeg-8.dll  ; for gdk_pixbuf loader.
	Delete $LIB_INSTDIR\libxslt-1.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libpcre-1.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libthai-0.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libdatrie-1.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libxml2-2.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libxml++-2.6-2.dll  ; fontconfig needs this
	Delete $LIB_INSTDIR\libxml++-3.0-1.dll
	Delete $LIB_INSTDIR\zlib1.dll  ; png and many others need this
	Delete $LIB_INSTDIR\libexpat-1.dll
	Delete $LIB_INSTDIR\libbz2-1.dll
	Delete $LIB_INSTDIR\libgraphite2.dll
	Delete $LIB_INSTDIR\librsvg-2-2.dll
	Delete $LIB_INSTDIR\libtiff-5.dll
	Delete $LIB_INSTDIR\libstdc++-6.dll
	Delete $LIB_INSTDIR\libgcc_s_seh-1.dll
	Delete $LIB_INSTDIR\libwinpthread-1.dll
	Delete $LIB_INSTDIR\libsoup-2.4-1.dll      ; libsoup
	Delete $LIB_INSTDIR\libsoup-gnome-2.4-1.dll      ; libsoup
	Delete $LIB_INSTDIR\libsqlite3-0.dll       ; libsoup dependency
	Delete $LIB_INSTDIR\libpsl-5.dll       ; libsoup dependency
	Delete $LIB_INSTDIR\libbrotlidec.dll       ; libsoup dependency
	Delete $LIB_INSTDIR\libbrotlicommon.dll       ; libsoup dependency
	Delete $LIB_INSTDIR\libgnutls-30.dll       ; glib-networking dependency
	Delete $LIB_INSTDIR\libgmp-10.dll		; glib-networking dependency
	Delete $LIB_INSTDIR\libhogweed-6.dll       ; glib-networking dependency
	Delete $LIB_INSTDIR\libnettle-8.dll	; glib-networking dependency
	Delete $LIB_INSTDIR\libidn2-0.dll		; glib-networking dependency
	Delete $LIB_INSTDIR\libp11-kit-0.dll	; glib-networking dependency
	Delete $LIB_INSTDIR\libtasn1-6.dll		; glib-networking dependency
	Delete $LIB_INSTDIR\libunistring-2.dll	; glib-networking dependency
	Delete $LIB_INSTDIR\libproxy-1.dll	; glib-networking dependency
	Delete $LIB_INSTDIR\libpeas-1.0-0.dll	
	Delete $LIB_INSTDIR\libpeas-gtk-1.0-0.dll
	Delete $LIB_INSTDIR\libgirepository-1.0-1.dll


FunctionEnd



var found_dir
var find_handle_lang_dir


Section Uninstall
	SetShellVarContext all  ; use all user variables as opposed to current user
	SetAutoClose false

	; Note: Checking if there is a registry key present, and using it to determine
	; if this is a private installation will not work, as it will break if a parallel
	; shared installation is present.

	; ReadRegStr $INST_DIR_REG HKLM "SOFTWARE\${PRODUCT_NAME}" "InstallationDirectory"
	; StrCmp $INST_DIR_REG "" uninst_no_sideeffects
	StrCmp $uninstall_option_sideeffects "no" uninst_no_sideeffects
		; For PATH removal
		ReadRegStr $LIB_INSTDIR HKLM "SOFTWARE\${PRODUCT_NAME}" "DllPath"
		ReadRegStr $DLL_DIR_NAME HKLM "SOFTWARE\${PRODUCT_NAME}" "DllDirName"
		ReadRegStr $ADD_TO_PATH HKLM "SOFTWARE\${PRODUCT_NAME}" "UsingSystemPath"
		ReadRegStr $SEC_TRANSLATIONS_INSTALLED HKLM "SOFTWARE\${PRODUCT_NAME}" "TranslationsInstalled"

		; DeleteRegKey HKLM "SOFTWARE\GTK\3.0"  ; dropline, etc...
		; DeleteRegKey /ifempty HKLM "SOFTWARE\GTK"  ; don't damage other installations

		DeleteRegKey HKLM "SOFTWARE\${PRODUCT_NAME}"
		DeleteRegKey HKLM "${PRODUCT_UNINST_KEY}"

		; FIXME: Do we have this registry key?
		; DeleteRegKey HKCU "Software\${PRODUCT_NAME}"

		Delete "$SMPROGRAMS\${PRODUCT_NAME}\Uninstall ${PRODUCT_NAME}.lnk"
		Delete "$SMPROGRAMS\${PRODUCT_NAME}\Go to the website.url"
		RMDir "$SMPROGRAMS\${PRODUCT_NAME}"  ; only if empty, theme selector may still be there

		; Remove GTK from $PATH
		; StrCmp $ADD_TO_PATH "0" un_nopath  ; Setting $PATH was not requested during installation
		; Push $LIB_INSTDIR
		; Call un.RemoveFromPath
		; EnVar::SetHKLM
		; EnVar::DeleteValue "PATH" "$LIB_INSTDIR"
		; un_nopath:

		; $DLL_DIR_NAME is from the registry here
		Call un.DeleteDlls

		goto delete_dlls_exit


	uninst_no_sideeffects:


		Strcpy $SEC_TRANSLATIONS_INSTALLED "0"
		StrCmp $uninstall_option_translations "yes" "" nodelete_translations
			Strcpy $SEC_TRANSLATIONS_INSTALLED "1"
		nodelete_translations:


		; All dll files. We delete them before /bin and /lib, so that
		; the directories are empty afterwards.

		; Since we have no registry, we have to remove dlls from all possible locations
		; StrCpy $DLL_DIR_NAME ""
		; Call un.DeleteDlls
		; StrCpy $DLL_DIR_NAME "bin"
		; Call un.DeleteDlls
		; StrCpy $DLL_DIR_NAME "lib"
		; Call un.DeleteDlls
		
		; Force users to use the command-line argument instead:
		; if silent, use the /dllpath= option
		StrCmp $uninstall_option_dllpath "bin" goto_undll_bin
		StrCmp $uninstall_option_dllpath "lib" goto_undll_lib
		StrCmp $uninstall_option_dllpath "root" goto_undll_none goto_undll_bin  ; default to bin if not matched

		goto_undll_none:
			StrCpy $DLL_DIR_NAME ""
			goto goto_undll_exit
		goto_undll_lib:
			StrCpy $DLL_DIR_NAME "lib"
			goto goto_undll_exit
		goto_undll_bin:
			StrCpy $DLL_DIR_NAME "bin"
			goto goto_undll_exit

		goto_undll_exit:
		    Call un.DeleteDlls

	delete_dlls_exit:


	; Delete config file?
	IfSilent "" read_config_page
		StrCmp $uninstall_option_remove_config "yes" delete_config skip_config

	read_config_page:
		!insertmacro INSTALLOPTIONS_READ $leave_config "nsi_configpage.ini" "Field 1" "State"
		StrCmp $leave_config "1" skip_config

	delete_config:
		;Delete "$INSTDIR\etc\gtk-2.0\gtkrc"
	skip_config:


	Delete "$INSTDIR\etc\fonts\fonts.conf"
	RMDir /r "$INSTDIR\etc\fonts\conf.d"
	RMDir "$INSTDIR\etc\fonts"  ; only if empty
	Delete "$INSTDIR\etc\gtk-3.0\im-multipress.conf"
	Delete "$INSTDIR\etc\gtk-3.0\settings.ini"
	RMDir "$INSTDIR\etc\gtk-3.0" ; only if empty
	RMDir "$INSTDIR\etc" ; only if empty

	; some helper files here
	RMDir /r "$INSTDIR\gtk3-runtime"

	RMDir "$INSTDIR\bin"  ; only if empty

	; RMDir /r "$INSTDIR\lib"
	; pango modules are gone as of gtk 2.10
	;  RMDir /r "$INSTDIR\lib\pango"
	
	Delete "$INSTDIR\lib\gdk-pixbuf-2.0\2.10.0\loaders.cache"
	RMDir /r "$INSTDIR\lib\gdk-pixbuf-2.0\2.10.0\loaders"  ; not forced
	RMDir "$INSTDIR\lib\gdk-pixbuf-2.0\2.10.0"  ; not forced
	RMDir "$INSTDIR\lib\gdk-pixbuf-2.0"  ; not forced


	Delete "$INSTDIR\lib\gio\modules\libgiognutls.dll"
	Delete "$INSTDIR\lib\gio\modules\libgiognomeproxy.dll"
	Delete "$INSTDIR\lib\gio\modules\libgiolibproxy.dll"
	RMDir "$INSTDIR\lib\gio\modules"
	RMDir "$INSTDIR\lib\gio"

	Delete "$INSTDIR\lib\libpeas-1.0\loaders\libpython3loader.dll"
	RMDir "$INSTDIR\lib\libpeas-1.0\loaders"
	RMDir "$INSTDIR\lib\libpeas-1.0"

	RMDir /r "$INSTDIR\lib\girepository-1.0"

	Delete "$INSTDIR\ssl\certs\ca-bundle.crt"
	Delete "$INSTDIR\ssl\certs\ca-bundle.trust.crt"
	RMDir "$INSTDIR\ssl\certs"
	RMDir "$INSTDIR\ssl"

	; no longer in gtk as of 2.14.5
	;RMDir /r "$INSTDIR\lib\gtk-3.0\${GTK_BIN_VERSION}\immodules"
	; gone as of gtk 2.16.6-2
	;RMDir /r "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}\loaders"

	;Delete "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}\engines\libwimp*.dll"
	; there should be no problem deleting this
	;Delete "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}\engines\libpixmap*.dll"

	;RMDir "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}\engines"  ; not forced
	;RMDir "$INSTDIR\lib\gtk-2.0\${GTK_BIN_VERSION}"  ; not forced
	;RMDir "$INSTDIR\lib\gtk-2.0"  ; not forced

	RMDir "$INSTDIR\lib"  ; not forced

	Delete "$INSTDIR\share\locale\locale.alias"  ; gettext

	RMDir "$INSTDIR\share\locale"  ; only if empty, not to remove the other programs' translations


	RMDir /r "$INSTDIR\share\themes\Default"
	RMDir /r "$INSTDIR\share\themes\Emacs"
	RMDir /r "$INSTDIR\share\glib-2.0"
	RMDir /r "$INSTDIR\share\icons"
	RMDir /r "$INSTDIR\share\gtksourceview-3.0"
	RMDir /r "$INSTDIR\share\gtksourceview-4"

	RMDir "$INSTDIR\share\themes"  ; not forced
	RMDir "$INSTDIR\share"  ; not forced

	Delete "$INSTDIR\spinechain_uninst.exe"

	RMDir "$INSTDIR"  ; delete only if empty

SectionEnd ; end of uninstall section






; --------------- Helpers



; WritePostInstall
Function WritePostInstall
	SetShellVarContext all  ; use all user variables as opposed to current user
	Pop $R0 ; Output file
	Push $R9
		FileOpen $R9 $R0 w
		FileWrite $R9 "@echo off$\r$\n"
		FileWrite $R9 "rem $\"$INSTDIR\bin\gdk-pixbuf-query-loaders.exe$\" > $\"$INSTDIR\etc\gtk-3.0\gdk-pixbuf.loaders$\"$\r$\n"
		FileWrite $R9 "$\"$INSTDIR\bin\gtk-query-immodules-3.0.exe$\" > $\"$INSTDIR\etc\gtk-3.0\gtk.immodules$\"$\r$\n"
		FileWrite $R9 "rem $\"$INSTDIR\bin\gtk-update-icon-cache.exe$\"$\r$\n"
		FileClose $R9
	Pop $R9
FunctionEnd




; Detect previous installation
Function DetectPrevInstallation
	; if /removeold=no option is given, don't check anything.
	StrCmp $install_option_removeold "no" old_detect_done

	SetShellVarContext all  ; use all user variables as opposed to current user
	push $R0

	; detect previous installation
	ReadRegStr $R0 HKLM "${PRODUCT_UNINST_KEY}" "UninstallString"
	StrCmp $R0 "" old_detect_done

	MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
		"${PRODUCT_NAME} is already installed. $\n$\nClick `OK` to remove the \
		previous version or `Cancel` to continue anyway." \
		/SD IDOK IDOK old_uninst
		; Abort
		goto old_detect_done

	; Run the old uninstaller
	old_uninst:
		ClearErrors
		IfSilent old_silent_uninst old_nosilent_uninst

		old_nosilent_uninst:
			ExecWait '$R0'
			goto old_uninst_continue

		old_silent_uninst:
			; We assume it's an NSIS-generated uninstaller.
			ExecWait '$R0 /S _?=$INSTDIR'

		old_uninst_continue:

		IfErrors old_no_remove_uninstaller

		; You can either use Delete /REBOOTOK in the uninstaller or add some code
		; here to remove to remove the uninstaller. Use a registry key to check
		; whether the user has chosen to uninstall. If you are using an uninstaller
		; components page, make sure all sections are uninstalled.
		old_no_remove_uninstaller:

	old_detect_done: ; old installation not found, all ok

	pop $R0
FunctionEnd



; Prevent running multiple instances of the installer
Function PreventMultipleInstances
	Push $R0
	System::Call 'kernel32::CreateMutexA(i 0, i 0, t ${PRODUCT_NAME}) ?e'
	Pop $R0
	StrCmp $R0 0 +3
		MessageBox MB_OK|MB_ICONEXCLAMATION "The installer is already running." /SD IDOK
		Abort
	Pop $R0
FunctionEnd




; eof
