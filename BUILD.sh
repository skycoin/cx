nex -e cxgo/cxgo0/cxgo0.nex
goyacc -o cxgo/cxgo0/cxgo0.go cxgo/cxgo0/cxgo0.y
nex -e cxgo/cxgo.nex
goyacc -o cxgo/cxgo.go cxgo/cxgo.y
go build -tags full -i -o $GOPATH/bin/cx ./cxgo/
