FROM asia.gcr.io/werun-project/go-builder:latest as builder 

RUN mkdir -p /go/src/baymax
COPY . /go/src/baymax
WORKDIR /go/src/baymax
RUN glide install 
RUN go test ./...