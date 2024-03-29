FROM golang:alpine AS builder

RUN apk --update add ca-certificates git

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

# Run the exe file
FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app .


EXPOSE 8080

ENV DATABASE_DSN="host=localhost user=user password=password dbname=airline port=5432"
ENV JWT_TOKEN="m1__sup3rs3cr3t!!"
ENV DEBUG=0

CMD ["./main"]
