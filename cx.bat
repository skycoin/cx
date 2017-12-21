@setlocal enableextensions enabledelayedexpansion
@echo off

SET UNSUPPORTED=go1.0 go1.1 go1.2 go1.3 go1.4 go1.5 go1.6 go1.7
FOR /F "tokens=* USEBACKQ" %%F IN (`go version`) DO (SET VERSION=%%F)
FOR /F "tokens=* USEBACKQ" %%F IN (`go env GOPATH`) DO (SET GOPATH=%%F)

FOR %%v IN (%UNSUPPORTED%) DO (
    set check="!VERSION:%%v=!"
    IF "!VERSION!" NEQ !check! (
       ECHO FAIL: CX requires Go version 1.8+. Please update your Go installation: https://golang.org/dl/
       EXIT /B 2
    )
)

rem echo %GOPATH%

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

   go get github.com/skycoin/cx/...

   IF ERRORLEVEL 1 (
      echo FAIL:   Couldn't install github.com/skycoin/cx
      EXIT /B 2
   ) ELSE (
      echo OK:     Package github.com/skycoin/cx was installed successfully
   )
)

IF EXIST %GOPATH%\src\github.com\skycoin\cx\ (
   echo NOTE:   Re-compiling CX
) ELSE (
   echo NOTE:   Compiling CX
)

%GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cx\cx.nex
IF ERRORLEVEL 1 (
   echo FAIL:   There was a problem compiling CX's lexical analyzer
   EXIT /B 2
)

%GOPATH%\bin\goyacc -o %GOPATH%\src\github.com\skycoin\cx\cx\cx.go %GOPATH%\src\github.com\skycoin\cx\cx\cx.y
IF ERRORLEVEL 1 (
   echo FAIL:   There was a problem compiling CX's parser
   EXIT /B 2
)

go install github.com/skycoin/cx/cx/
IF ERRORLEVEL 1 (
   echo FAIL:   There was a problem compiling CX
   EXIT /B 2
) ELSE (
   echo OK      CX was compiled successfully
)

echo NOTE:\tWe recommend you to test your CX installation by running 'cx \$GOPATH/src/github.com/skycoin/cx/tests/test.cx'
