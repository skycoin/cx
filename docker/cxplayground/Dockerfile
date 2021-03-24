FROM golang:1.14-stretch AS build

WORKDIR /go/src/github.com/skycoin/cx/

ADD cmd/cxplayground cmd/cxplayground
ADD cx cx
ADD cxgo cxgo
ADD vendor vendor

ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/ ./cmd/...

FROM busybox

WORKDIR /var/cxplayground

COPY --from=build /usr/local/bin/cxplayground /usr/local/bin/cxplayground
COPY --from=build /go/src/github.com/skycoin/cx/cmd/cxplayground/dist /var/cxplayground/dist
COPY --from=build /go/src/github.com/skycoin/cx/cmd/cxplayground/examples /var/cxplayground/examples


ENTRYPOINT ["cxplayground"]
