FROM --platform=$BUILDPLATFORM golang:bullseye

ARG TARGETARCH
WORKDIR /go/src/github.com/luoch/progress.go/

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.com.cn,direct

COPY . /go/src/github.com/luoch/progress.go/

RUN GOOS=linux GOARCH=$TARGETARCH go build -o progress .
EXPOSE 8000
RUN chmod +x /go/src/github.com/luoch/progress.go/progress
CMD ["/go/src/github.com/luoch/progress.go/progress"]
