FROM golang:1.21 as builder
WORKDIR app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=builder /app/cloudrun .
ENTRYPOINT ["./cloudrun"]