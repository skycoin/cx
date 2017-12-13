#! /bin/sh

if [ ! -d "$GOPATH/src/github.com/skycoin/skycoin/" ]; then
    echo "NOTE:\tRepository github.com/skycoin/skycoin is not present in $GOPATH"
    echo "NOTE:\tDownloading the repository and installing the package via 'go get github.com/skycoin/skycoin/...'"

    go get github.com/skycoin/skycoin/...
    
    if [ $? -eq 0 ]; then
        echo "OK:\tPackage github.com/skycoin/skycoin was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/skycoin/skycoin"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/go-gl/gl/v2.1/gl" ]; then
    echo "NOTE:\tRepository github.com/go-gl/gl/v2.1/gl is not present in $GOPATH"
    echo "NOTE:\tInstalling it via 'go get github.com/go-gl/gl/v2.1/gl'"
    
    go get github.com/go-gl/gl/v2.1/gl

    if [ $? -eq 0 ]; then
        echo "OK:\tRepository github.com/go-gl/gl/v2.1/gl was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/go-gl/gl/v2.1/gl"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/go-gl/glfw/v3.2/glfw" ]; then
    echo "NOTE:\tRepository github.com/go-gl/glfw/v3.2/glfw is not present in $GOPATH"
    echo "NOTE:\tInstalling it via 'go get github.com/go-gl/glfw/v3.2/glfw'"
    
    go get github.com/go-gl/glfw/v3.2/glfw
    
    if [ $? -eq 0 ]; then
        echo "OK:\tRepository github.com/go-gl/glfw/v3.2/glfw was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/go-gl/glfw/v3.2/glfw"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/go-gl/gltext" ]; then
    echo "NOTE:\tRepository src/github.com/go-gl/gltext is not present in $GOPATH"
    echo "NOTE:\tInstalling it via 'go get github.com/go-gl/gltext'"
    
    go get github.com/go-gl/gltext
    
    if [ $? -eq 0 ]; then
        echo "OK:\tRepository github.com/go-gl/gltext was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/go-gl/gltext"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/blynn/nex" ]; then
    echo "NOTE:\tRepository github.com/blynn/nex is not present in $GOPATH"
    echo "NOTE:\tInstalling it via 'go get github.com/blynn/nex'"
    
    go get github.com/blynn/nex
    
    if [ $? -eq 0 ]; then
        echo "OK:\tRepository github.com/blynn/nex was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/blynn/nex"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/cznic/goyacc" ]; then
    echo "NOTE:\tRepository github.com/cznic/goyacc is not present in $GOPATH"
    echo "NOTE:\tInstalling it via 'go get github.com/cznic/goyacc'"
    
    go get github.com/cznic/goyacc
    
    if [ $? -eq 0 ]; then
        echo "OK:\tRepository github.com/cznic/goyacc was installed successfully"
    else
        echo "FAIL:\tCouldn't install github.com/cznic/goyacc"
        exit 0
    fi
fi

if [ ! -d "$GOPATH/src/github.com/skycoin/cx/" ]; then
    echo "NOTE:\tRepository github.com/skycoin/cx is not present in $GOPATH"
    echo "NOTE:\tDownloading the repository and installing the package via 'go get github.com/skycoin/cx/...'"

    go get github.com/skycoin/cx/...
    
    if [ $? -eq 0 ]; then
        echo "OK:\tPackage github.com/skycoin/cx was installed successfully"
    else
        echo "FAIL:\tCouldn't clone into github.com/skycoin/cx"
        exit 0
    fi
fi

if [ -f "$GOPATH/bin/cx" ]
then
    echo "NOTE:\tRe-compiling CX"
else
    echo "NOTE:\tCompiling CX"
fi

nex -e $GOPATH/src/github.com/skycoin/cx/cx/cx.nex
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's lexical analyzer"
    exit 0
fi

goyacc -o $GOPATH/src/github.com/skycoin/cx/cx/cx.go $GOPATH/src/github.com/skycoin/cx/cx/cx.y
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's parser"
    exit 0
fi

go install github.com/skycoin/cx/cx/
if [ $? -eq 0 ]; then
    echo "OK:\tCX was compiled successfully"
else
    echo "FAIL:\tThere was a problem compiling CX"
    exit 0
fi

echo "NOTE:\tWe recommend you to test your CX installation by running 'cx \$GOPATH/src/github.com/skycoin/cx/tests/test.cx'"

