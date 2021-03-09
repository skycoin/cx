@setlocal enableextensions enabledelayedexpansion
@echo off



rem start banner
call :setColorVarsIfOnWin10OrLater
rem either insert vars like %COLOR_RED% into normal echo commands or do like:
rem '
rem call :echoWithColor cyan "text goes here"
rem '


echo %COLOR_CYAN%---------------------------------------------
echo ^|             CX SETUP ^& UPDATE             ^|
echo ---------------------------------------------%COLOR_DEFAULT%
echo GREEN text is because of SUCCESSful command execution...
echo ...but Windows 10 and up is required to see color
echo %COLOR_YELLOW%



echo Checking your environment variables


rem set GI_PATH [go installation path]
if defined GOPATH (
    set GI_PATH=%GOPATH%
    call :echoWithColor green "OK:     GOPATH:"
    call :echoWithColor green "             %GOPATH%"
) else (
    set GI_PATH=%USERPROFILE%\go
    call :echoWithColor red "ERROR:  GOPATH not set^^^!  FIX IT^^^!"
    rem GI_PATH doesn't seem to ACTUALLY be set, until we exit this code block
    call :echoWithColor red "     ....so we'll use %USERPROFILE%\go for Go's install path"
)

rem set CI_PATH [CX installation path]
if defined CXPATH (
    set CI_PATH=%CXPATH%
    call :echoWithColor green "OK:     CXPATH, CX's workspace:"
    call :echoWithColor green "             %CXPATH%"
) else (
    set CI_PATH=%USERPROFILE%\cx
    call :echoWithColor red "ERROR:  CX's workspace environment variable; CXPATH is not set^^^!  FIX IT^^^!"
    call :echoWithColor red "     ....so we'll use %USERPROFILE%\cx for CX's workspace"
)

rem check for %BIN_PATH% in %PATH%
set BIN_PATH=%GI_PATH%\bin
set GH_PATH=%GI_PATH%\src\github.com
set SKYCOIN_PATH=%GH_PATH%\skycoin
set CXGO_PATH=%SKYCOIN_PATH%\cx\cxgo
rem attempt to replace %BIN_PATH% text with nothing...
call set COMPARISON_PATH=%%PATH:%BIN_PATH%=%%
rem ...if found/replaced %BIN_PATH%, below vars aren't equal

if "%PATH%"=="%COMPARISON_PATH%" (
   call :echoWithColor red "ERROR:  Your PATH var needs %BIN_PATH%"
) else (
   call :echoWithColor green "OK:     Your PATH includes %BIN_PATH%"
)



call :checkVersions



call :echoWithColor yellow "Checking your hard disk"


rem ensure CX workspace exists
call :echoWithColor yellow "     Ensuring that we have CX workspace folders....."
call :setupCXWorkspaceDir %CI_PATH%
call :setupCXWorkspaceDir %CI_PATH%\src
call :setupCXWorkspaceDir %CI_PATH%\bin
call :setupCXWorkspaceDir %CI_PATH%\pkg

rem get repositories that don't change often.
rem user will need to manually rebuild these when they DO change
call :echoWithColor yellow "     Ensuring that we have local repositories....."
call :getRepo go-gl\gl\v2.1\gl
call :getRepo go-gl\glfw\v3.2\glfw
call :getRepo go-gl\gltext
call :getRepo cznic\goyacc
call :getRepo skycoin\skycoin



call :echoWithColor yellow "     Ensuring we're setup with the latest CX....."


rem always build latest cx.exe
call :cloneOrPullLatest skycoin CX
call :buildCX

rem show CX version
echo %COLOR_YELLOW%
%BIN_PATH%\cx.exe -v
echo %COLOR_RESET%
if %ERRORLEVEL% neq 0 (
   call :echoWithColor red "ERROR:  %BIN_PATH% needs to be in your PATH environment variable"
   call :echoWithColor red "AND cx.exe must be accessible"
   goto :errorReport
)



rem run CX tests
%BIN_PATH%\cx.exe %SKYCOIN_PATH%\cx\lib\args.cx %SKYCOIN_PATH%\cx\tests\main.cx ++wdir=tests\ ++disable-tests=gui,issue
call :showResults cx\tests "Tested" "ERROR while testing"

call :echoWithColor yellow "You can re-run CX tests with:"
call :echoWithColor yellow "     'cx.exe %%%%GOPATH%%%%\src\github.com\skycoin\cx\tests'"



rem final report
call :echoWithColor cyan "FINISHED CX SETUP^^^!"
call :echoWithColor cyan "     Make sure to set any NEEDED environment variables."
call :echoWithColor cyan "     Which MAY have been shown as ERRORs, at the start of CX Setup."
rem main program flow must end with 'exit /b' or 'goto :EOF', which means end of file
goto :EOF











rem SUBROUTINES/FUNCTIONS will return to the next line after their call [upon 'exit /b']

:setColorVarsIfOnWin10OrLater
  rem Windows versions prior to 10, had no easy way to do colors in batch files.
  rem YES, this will probably show color code gibberish if you are running Windows 1.x

  ver | find "Version 1" > nul
  if %ERRORLEVEL% neq 0 (echo Running Windows 8.x or earlier, won't use text colors)
  if %ERRORLEVEL% neq 0 (exit /b)

  set COLOR_DEFAULT=[0m
  set COLOR_RESET=[0m
  set COLOR_GRAY=[90m
  set COLOR_GREY=[90m
  set COLOR_RED=[91m
  set COLOR_GREEN=[92m
  set COLOR_YELLOW=[93m
  set COLOR_BLUE=[94m
  set COLOR_MAGENTA=[95m
  set COLOR_CYAN=[96m
  set COLOR_WHITE=[97m
exit /b


:checkVersions
  set UNSUPPORTED=go1.0 go1.1 go1.2 go1.3 go1.4 go1.5 go1.6 go1.7 go1.8 go1.9
  for /F "tokens=* USEBACKQ" %%F IN (`go version`) DO (set VERSION=%%F)
  set VERSION=%VERSION:go version =%
  set VERSION=%VERSION: windows/amd64=%
  call :echoWithColor yellow "Go version is: %VERSION%"

  rem check against old versions
  for %%v IN (%UNSUPPORTED%) DO (
      echo Checking for UNSUPPORTED version: %%v
      if %VERSION% equ %%v (
         echo %COLOR_RED% Yours is too old^! %COLOR_DEFAULT%
         echo %COLOR_RED% ERROR: CX requires Go version 1.10+.  https://golang.org/dl/ %COLOR_DEFAULT%
      )
  )
exit /b


:getRepo
  set GH=github.com\

  if "%1" == "skycoin\skycoin" (
    set PARAM=%1\...
  ) else (
    set PARAM=%1
  )

  if exist %GH_PATH%\%1 (
    call :echoWithColor white "  Already got %PARAM%"
  ) else (
    call :echoWithColor white "  Issuing command: 'go get %GH%%PARAM%'"
    go get %GH%%PARAM%
    call :showResults %PARAM% "Just got" "ERROR GETting"
  )
exit /b


:cloneOrPullLatest
  rem git clone https://github.com/skycoin/cx.git %GOPATH%\src\github.com\skycoin\cx
  set CLONE_CMD=git clone https://github.com/%1/%2.git %GH_PATH%\%1\%2

  if exist %GH_PATH%\%1\%2\ (
    rem -------- PULL
    call :echoWithColor white "  Already got %1\%2"
    call :echoWithColor yellow "            Pulling %2"

    rem cx path
    cd %SKYCOIN_PATH%\%2
    call :showResults %SKYCOIN_PATH%\%2 "Changed DIR to" "Couldn't change DIR to"

    rem pull
    git pull
    rem someone recommended 'git pull origin master'
    call :showResults %1\%2 "Pulled" "ERROR WHILE PULLING"

    call :echoWithColor yellow "            Re-building %2"
  ) else (
    rem -------- CLONE
    call :echoWithColor white "NEED Repository %1\%2"
    call :echoWithColor white "  Issuing command: '%CLONE_CMD%'"

    rem clone
    %CLONE_CMD%
    call :showResults %1\%2 "Cloned" "ERROR Cloning"

    call :echoWithColor yellow "            Building %2"
  )
exit /b


:buildCX
  %BIN_PATH%\goyacc -o %CXGO_PATH%\cxgo0\cxgo0.go %CXGO_PATH%\cxgo0\cxgo0.y
  call :showResults "goyacc cxgo0" "1st pass -" "ERROR in 1st pass -"

  %BIN_PATH%\goyacc -o %CXGO_PATH%\parser\cxgo.go %CXGO_PATH%\parser\cxgo.y
  call :showResults "goyacc cxgo" "2nd pass -" "ERROR in 2nd pass -"


  go build -tags="base cxfx" -i -o %BIN_PATH%\cx.exe github.com\skycoin\cx\cmd\cx
  call :showResults skycoin\CX\CMD\CX "            Built CX.EXE from:" "ERROR building CX.EXE from:"
exit /b


:setupCXWorkspaceDir
  if exist %1 (
    call :echoWithColor white "  Already made folder '%1'"
  ) else (
    mkdir %1
    call :showResults %1 "  Created folder '%1'" "  ERROR creating folder '%1'^^^!"
  )
exit /b



:echoWithColor
         if "%1"=="gray"    ( echo %COLOR_GRAY%%~2%COLOR_RESET%
  ) else if "%1"=="grey"    ( echo %COLOR_GREY%%~2%COLOR_RESET%
  ) else if "%1"=="red"     ( echo %COLOR_RED%%~2%COLOR_RESET%
  ) else if "%1"=="green"   ( echo %COLOR_GREEN%%~2%COLOR_RESET%
  ) else if "%1"=="yellow"  ( echo %COLOR_YELLOW%%~2%COLOR_RESET%
  ) else if "%1"=="blue"    ( echo %COLOR_BLUE%%~2%COLOR_RESET%
  ) else if "%1"=="magenta" ( echo %COLOR_MAGENTA%%~2%COLOR_RESET%
  ) else if "%1"=="cyan"    ( echo %COLOR_CYAN%%~2%COLOR_RESET%
  ) else if "%1"=="white"   ( echo %COLOR_WHITE%%~2%COLOR_RESET%
  ) else (
    rem invalid color string...
    echo %COLOR_RED% ERROR... :echoWithColor was passed an UNRECOGNIZED COLOR: %1
    echo %COLOR_DEFAULT% %~2
  )
exit /b


rem parameters
rem 1st - most relevant target of the command; usually a path
rem      put texts below in double quotes:
rem 2nd - successful operation text
rem 3rd - ERROR/FAILURE text
:showResults
  if %ERRORLEVEL% equ 0 (
     call :echoWithColor green "%~2 %~1"
  ) else (
     call :echoWithColor red "%~3 %~1"
     goto :errorReport
  )
exit /b


:errorReport
  call :echoWithColor red "********* ERROR CODE: %ERRORLEVEL% *********"
  rem call :echoWithColor red "********* ABORTING SETUP^^^!  Due to ERROR CODE: %ERRORLEVEL% *********"
  rem call :echoWithColor red "You NEED to run CXSETUP again after FIXing that"
exit /b %ERRORLEVEL%
