FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

FROM alpine

COPY --from=builder /server /server

CMD ["/server"]