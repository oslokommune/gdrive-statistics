# golang:1.15.5
FROM golang@sha256:cf46c759511d0376c706a923f2800762948d4ea1a9290360720d5124a730ed63 AS builder

WORKDIR /build

#COPY go.mod go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

# https://rollout.io/blog/building-minimal-docker-containers-for-go-applications/
# TL;DR: This makes the built binary run on alpine linux
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -a -installsuffix cgo -o main .

# alpine:3.12.1
FROM alpine@sha256:c0e9560cda118f9ec63ddefb4a173a2b2a0347082d7dff7dc14272e7841a5b5a

LABEL org.opencontainers.image.source https://github.com/oslokommune/gdrive-statistics

WORKDIR /app

COPY --from=builder /build/* ./

RUN chmod +x /app/main

CMD ["/app/main" ]
