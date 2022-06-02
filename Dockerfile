FROM golang:1.18-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY . ./

RUN apk add --no-cache \
    gcc \
    musl-dev

RUN go mod download && go mod verify

RUN CGO_ENABLED=1 go build -o jobz -ldflags '-s -w -extldflags "-static"' ./cmd/jobz

FROM scratch

ENV GIN_MODE=release

COPY --from=builder /usr/src/app/jobz /usr/local/bin/jobz
COPY --from=builder /usr/src/app/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/jobz", "--migrations=./migrations"]
