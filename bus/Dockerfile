FROM golang:1.25.5 AS builder
WORKDIR /app
COPY ./bus/go.mod ./
RUN go mod download
COPY ./bus/. .
RUN go build -o bus .

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/bus .
COPY ./bus/dataset.json .
CMD ["./bus"]