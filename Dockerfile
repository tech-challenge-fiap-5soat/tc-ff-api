FROM golang:1.21 AS builder

WORKDIR /app

COPY cmd ./cmd
COPY src ./src

COPY go.mod go.sum ./

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g ./cmd/go-command.go 

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o fiap-tech-fast-food ./cmd/go-command.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/fiap-tech-fast-food /app/fiap-tech-fast-food

COPY --from=builder /app/src/external/api/infra/config/configs.yaml.sample /app/data/configs/configs.yaml

EXPOSE 8080

ENTRYPOINT ["/app/fiap-tech-fast-food"]
