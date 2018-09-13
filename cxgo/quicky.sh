#! /bin/sh

rm *.go
nex -e cxgo.nex
goyacc -p "p1" -o cxgo.go cxgo.y
goyacc -p "p0" -o cxgo0.go cxgo0.y
go build -i -o cx .
