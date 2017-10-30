#! /bin/sh
go get github.com/blynn/nex
go get github.com/cznic/goyacc
nex -e src/cxgo/cx.nex && goyacc -o src/cxgo/cx.go src/cxgo/cx.y && go build src/cxgo/cx.go src/cxgo/cx.nn.go
