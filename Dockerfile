FROM golang:1.22.4 as builder
COPY ./ /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/go-oms -ldflags "-w -s" ./main.go

FROM alpine:3.19
COPY --from=builder /bin/go-oms /
RUN apk add --no-cache tzdata
ENV TZ=Asia/Bangkok

EXPOSE 8080

CMD ["/go-oms"]
