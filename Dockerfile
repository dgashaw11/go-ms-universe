FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 go build -o /app ./cmd/${SERVICE}

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app /app
COPY migrations /migrations

ENTRYPOINT ["/app"]
