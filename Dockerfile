FROM alpine:3.4

RUN apk add --update go git

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin:/usr/local/bin

WORKDIR /go/src/github.com/lavrs/proxy
ADD . /go/src/github.com/lavrs/proxy

RUN go install \
    && go build

ENTRYPOINT ["/go/bin/proxy/cmd/proxy/proxy"]