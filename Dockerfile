FROM golang:alpine AS builder
MAINTAINER scoful
WORKDIR /go/build
COPY . /go/build
ENV GOPROXY https://goproxy.cn,direct
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o ocsp-proxy main.go

FROM scratch AS runner
COPY --from=builder /go/build/ocsp-proxy /
EXPOSE 8080
ENTRYPOINT ["/ocsp-proxy"]