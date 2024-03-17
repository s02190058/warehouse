FROM golang:1.22.1 AS modules

WORKDIR /modules

COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.22.1 AS builder

COPY --from=modules /go/pkg /go/pkg

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/warehouse ./cmd/warehouse

FROM alpine:latest

COPY --from=builder /app/configs/docker.yml /configs/docker.yml

COPY --from=builder /bin/warehouse /warehouse

EXPOSE 8080

CMD ["/warehouse", "--config", "/configs/docker.yml"]