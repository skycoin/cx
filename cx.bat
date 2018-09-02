@setlocal enableextensions enabledelayedexpansion
@echo off

SET UNSUPPORTED=go1.0 go1.1 go1.2 go1.3 go1.4 go1.5 go1.6 go1.7
FOR /F "tokens=* USEBACKQ" %%F IN (`go version`) DO (SET VERSION=%%F)
FOR /F "tokens=* USEBACKQ" %%F IN (`go env GOPATH`) DO (SET GOPATH=%%F)
SET VERSION=%VERSION:go version=%
SET VERSION=%VERSION:windows/amd64=%
SET VERSION=%VERSION: =%
ECHO NOTE: go version is: %VERSION% 

FOR %%v IN (%UNSUPPORTED%) DO (
    ::SET CHECK="!VERSION:%%v=!"
    ECHO NOTE: %%v
    IF %VERSION% EQU %%v (
       ECHO FAIL: CX requires Go version 1.8+. Please update your Go installation: https://golang.org/dl/
       EXIT /B 2
    )
)

rem echo %GOPATH%

IF NOT EXIST %GOPATH% (
    SET INSTALLATION_PATH=%USERPROFILE%\go
    echo NOTE:   \%GOPATH\% not set; using: %INSTALLATION_PATH%
) ELSE (
    SET INSTALLATION_PATH=%GOPATH%
)

echo NOTE:   Make sure to add CX executable's directory at %INSTALLATION_PATH%\go\bin to your \%PATH\% environment variable

IF NOT EXIST %GOPATH%\src\github.com\skycoin\skycoin\ (
   echo NOTE:   Repository github.com\skycoin\skycoin is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/skycoin/skycoin/...'

   go get github.com/skycoin/skycoin/...

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/skycoin/skycoin
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/skycoin/skycoin was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\go-gl\gl\v2.1\gl\ (
   echo NOTE:   Repository github.com\go-gl\gl\v2.1\gl is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/go-gl/gl/v2.1/gl'

   go get github.com/go-gl/gl/v2.1/gl

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/go-gl/gl/v2.1/gl
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/go-gl/gl/v2.1/gl was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\go-gl\glfw\v3.2\glfw\ (
   echo NOTE:   Repository github.com\go-gl\glfw\v3.2\glfw is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/go-gl/glfw/v3.2/glfw'

   go get github.com/go-gl/glfw/v3.2/glfw

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/go-gl/glfw/v3.2/glfw
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/go-gl/glfw/v3.2/glfw was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\go-gl\gltext\ (
   echo NOTE:   Repository github.com\go-gl\gltext is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/go-gl/gltext'

   go get github.com/go-gl/gltext

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/go-gl/gltext
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/go-gl/gltext was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\blynn\nex\ (
   echo NOTE:   Repository github.com\blynn\nex is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/blynn/nex'

   go get github.com/blynn/nex

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/blynn/nex
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/blynn/nex was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\cznic\goyacc\ (
   echo NOTE:   Repository github.com\cznic\goyacc is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/cznic/goyacc'

   go get github.com/cznic/goyacc

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/cznic/goyacc
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/cznic/goyacc was installed successfully
   )
)

IF NOT EXIST %GOPATH%\src\github.com\skycoin\cx\ (
   echo NOTE:   Repository github.com\skycoin\cx is not present in %GOPATH%
   echo NOTE:   Downloading the repository and installing the package via 'go get github.com/skycoin/cx'
   
   git clone https://github.com/skycoin/cx.git %GOPATH%\src\github.com\skycoin\cx

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/skycoin/cx
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/skycoin/cx was installed successfully
   )
)

IF EXIST %GOPATH%\src\github.com\skycoin\cx\ (
   echo NOTE:   Re-compiling CX
   cd %GOPATH%\src\github.com\skycoin\cx && git pull 
   %GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.nex

   %GOPATH%\bin\goyacc -o %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.go %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.y

   %GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.nex

   %GOPATH%\bin\goyacc -o %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.go %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.y
   
   go build -i -o %GOPATH%/bin/cx.exe github.com/skycoin/cx/cxgo/
) ELSE (
   %GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.nex

   %GOPATH%\bin\goyacc -o %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.go %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.y

   %GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.nex

   %GOPATH%\bin\goyacc -o %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.go %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo.y

   go build -i -o %GOPATH%/bin/cx.exe github.com/skycoin/cx/cxgo/
   echo NOTE:   Compiling CX
)

IF NOT DEFINED CX_PATH (SET CX_PATH=%USERPROFILE%\cx)

IF NOT EXIST %CX_PATH% (
   echo NOTE:   CX's workspace %CX_PATH% does not exist
   echo NOTE:   Creating CX's workspace at %CX_PATH%

   mkdir %CX_PATH%
   mkdir %CX_PATH%\src
   mkdir %CX_PATH%\bin
   mkdir %CX_PATH%\pkg

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't create CX's workspace at %CX_PATH%
      EXIT /B 2
   ) ELSE (
      echo OK:     CX's workspace was successfully created
   )
)

echo 
cx.exe -v
echo 

IF ERRORLEVEL 1 (
   echo FAIL:   Is CX executable's directory %INSTALLATION_PATH%\bin in your \%PATH\% environment variable?
EXIT /B 2
)

echo NOTE: We recommend you to test your CX installation by running 'cx.exe %GOPATH%\src\github.com\skycoin\cx\tests'
