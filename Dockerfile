FROM asia.gcr.io/werun-project/go-builder:latest as builder 

RUN mkdir -p /go/src/github.com/werunclub/baymax
COPY . /go/src/github.com/werunclub/baymax
WORKDIR /go/src/github.com/werunclub/baymax
RUN glide install
RUN cd rpc/examples && go run rpc_client.go