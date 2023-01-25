FROM golang:1.19-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest AS app

WORKDIR /app

COPY --from=builder /app/bin/kvevri .

EXPOSE 8080

CMD ["./kvevri"]
