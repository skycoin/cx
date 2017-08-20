#! /bin/sh
go get github.com/blynn/nex
go get github.com/cznic/goyacc
go get github.com/mndrix/golog
nex src/golike/cx.nex && goyacc -o src/golike/cx.go src/golike/cx.y && go build src/golike/cx.go src/golike/cx.nn.go
