FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gr-scanner ./cmd/gr-scanner

FROM gcr.io/distroless/static-debian12

USER nonroot:nonroot

COPY --from=builder /app/gr-scanner /gr-scanner

ENTRYPOINT ["/gr-scanner"]