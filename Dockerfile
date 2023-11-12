FROM golang:1.20-alpine3.18 AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct\
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o app .

FROM scratch
COPY --from=builder /build/app .
COPY --from=builder /build/conf /conf

EXPOSE 8000
ENTRYPOINT ["./app"]