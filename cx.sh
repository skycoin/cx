#! /bin/sh

## Version checking grabbed from dex's script https://github.com/coreos/dex/blob/master/scripts/check-go-version

VERSION=$( go version )

UNSUPPORTED=( "go1.2" "go1.3" "go1.4" "go1.5" "go1.6" "go1.7" )

MAJOR_GOVERSION=$( echo -n "$VERSION" | grep -o 'go1\.[0-9]' || true )
#FULL_GOVERSION=$( echo -n "$VERSION" | grep -o 'go1\.[0-9|\.]*' || true )

for V in "${UNSUPPORTED[@]}"; do
    if [ "$V" = "$MAJOR_GOVERSION" ]; then
        >&2 echo "FAIL: CX requires Go version 1.8+. Please update your Go installation: https://golang.org/dl/"
        exit 2
    fi
done

for V in "${KNOWN_INSECURE[@]}"; do
    if [ "$V" = "$FULL_GOVERSION" ]; then
        >&2 echo "Go version ${V} has known security vulnerabilities which impact CX. Please update your Go version."
        exit 2
    fi
done


## determining if $GOPATH is set
## if not, use $HOME/go, as `go get` uses this instead
if [ -z ${GOPATH+x} ];
then
    INSTALLATION_PATH=$HOME/go
    echo "NOTE:\t\$GOPATH not set; using: $INSTALLATION_PATH"
else
    INSTALLATION_PATH=$GOPATH
fi


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

git pull origin master

if [ -f "$GOPATH/bin/cx" ]
then
    echo "NOTE:\tRe-compiling CX"
else
    echo "NOTE:\tCompiling CX"
fi

$GOPATH/bin/nex -e $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo0/cxgo0.nex
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's lexical analyzer (first pass)"
    exit 0
fi

$GOPATH/bin/goyacc -o $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo0/cxgo0.go $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo0/cxgo0.y
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's parser (first pass)"
    exit 0
fi

$GOPATH/bin/nex -e $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo.nex
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's lexical analyzer"
    exit 0
fi

$GOPATH/bin/goyacc -o $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo.go $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo.y
if [ ! $? -eq 0 ]; then
    echo "FAIL:\tThere was a problem compiling CX's parser"
    exit 0
fi

# go build -o build/wiki wiki.go 

# go build -o $GOPATH/src/github.com/skycoin/cx/cxgo/cxgo github.com/skycoin/cx/cxgo/
go build -i -o $GOPATH/bin/cx github.com/skycoin/cx/cxgo/
chmod +x $GOPATH/bin/cx
# go install github.com/skycoin/cx/cxgo/
if [ $? -eq 0 ]; then
    echo "OK:\tCX was compiled successfully"
else
    echo "FAIL:\tThere was a problem compiling CX"
    exit 0
fi

echo "NOTE:\tWe recommend you to test your CX installation by running 'cx \$GOPATH/src/github.com/skycoin/cx/tests/test.cx'"

