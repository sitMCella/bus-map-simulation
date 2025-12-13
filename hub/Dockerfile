FROM golang:1.25.5 AS builder
WORKDIR /app
COPY ./hub/go.mod ./hub/go.sum ./
RUN go mod download
COPY ./hub/. .
RUN go build -o hub .

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/hub .
EXPOSE 9090
CMD ["./hub"]