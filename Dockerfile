FROM golang:alpine as builder
WORKDIR /go/src/application
RUN go env -w GO111MODULE=on
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o application .

FROM alpine:latest

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /root/

COPY --from=builder /go/src/application/application .
CMD ["./application"]