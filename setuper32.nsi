OutFile "imgconverts_install32.exe"

InstallDir "c:\tts\imgconverts"

Section

#CreateDirectory $INSTDIR\img
#CreateDirectory $INSTDIR\platforms
#CreateDirectory $INSTDIR\Qt\labs\folderlistmodel
#CreateDirectory $INSTDIR\QtQuick\Controls
#CreateDirectory $INSTDIR\QtQuick\Dialogs
#CreateDirectory $INSTDIR\QtQuick\Layouts
#CreateDirectory $INSTDIR\QtQuick\Window.2
#CreateDirectory $INSTDIR\QtQuick.2

SetOutPath $INSTDIR
File imgconverts.exe
File $%PKG_CONFIG_PATH%\..\..\bin\libgcc_s_dw2-1.dll
File $%PKG_CONFIG_PATH%\..\..\bin\libstdc++-6.dll
File $%PKG_CONFIG_PATH%\..\..\bin\libwinpthread-1.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Core.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Gui.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Network.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Qml.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Quick.dll
File $%PKG_CONFIG_PATH%\..\..\bin\Qt5Widgets.dll


SetOutPath $INSTDIR\img
File img\out.png

SetOutPath $INSTDIR\platforms
File $%PKG_CONFIG_PATH%\..\..\plugins\platforms\qwindows.dll

SetOutPath $INSTDIR\Qt\labs\folderlistmodel
File $%PKG_CONFIG_PATH%\..\..\qml\Qt\labs\folderlistmodel\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\Qt\labs\folderlistmodel\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\Qt\labs\folderlistmodel\qmlfolderlistmodelplugin.dll


SetOutPath $INSTDIR\QtQuick\Controls
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Controls\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Controls\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Controls\qtquickcontrolsplugin.dll


SetOutPath $INSTDIR\QtQuick\Dialogs
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Dialogs\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Dialogs\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Dialogs\dialogplugin.dll


SetOutPath $INSTDIR\QtQuick\Layouts
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Layouts\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Layouts\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Layouts\qquicklayoutsplugin.dll

SetOutPath $INSTDIR\QtQuick\Window.2
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Window.2\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Window.2\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick\Window.2\windowplugin.dll

SetOutPath $INSTDIR\QtQuick.2
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick.2\plugins.qmltypes
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick.2\qmldir
File $%PKG_CONFIG_PATH%\..\..\qml\QtQuick.2\qtquick2plugin.dll

WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTSImgConverts" "DisplayName" "TTS Watermark Converter"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTSImgConverts" "UninstallString" "$\"$INSTDIR\uninstaller.exe$\""
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTSImgConverts" "QuietUninstallString" "$\"$INSTDIR\uninstaller.exe$\" /S"

CreateShortCut "$DESKTOP\ImgConverters.lnk" "$INSTDIR\imgconverts.exe" ""

CreateDirectory "$SMPROGRAMS\TTS"
CreateDirectory "$SMPROGRAMS\TTS\ImgConverts"
CreateShortCut "$SMPROGRAMS\TTS\ImgConverts\Uninstall.lnk" "$INSTDIR\uninstaller.exe" "" "$INSTDIR\uninstaller.exe" 0
CreateShortCut "$SMPROGRAMS\TTS\ImgConverts\ImgConverters.lnk" "$INSTDIR\imgconverts.exe" "" "$INSTDIR\imgconverts.exe" 0

WriteUninstaller $INSTDIR\uninstaller.exe

SectionEnd

Section "Uninstall"

DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TTSImgConverts"

Delete "$DESKTOP\ImgConverters.lnk"
Delete "$SMPROGRAMS\TTS\ImgConverts\*.*"
RmDir  "$SMPROGRAMS\TTS\ImgConverts\"

Delete $INSTDIR\*.exe
Delete $INSTDIR\*.dll


RMDir /r $INSTDIR\img
RMDir /r $INSTDIR\platforms

RMDir /r $INSTDIR\Qt
RMDir /r $INSTDIR\QtQuick
RMDir /r $INSTDIR\QtQuick.2

SectionEnd