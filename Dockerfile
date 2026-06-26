# syntax=docker/dockerfile:1
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -buildvcs=false -o /usr/local/bin/compactor ./cmd/compactor

# Runtime stage
FROM node:lts-alpine AS runtime

RUN apk add --no-cache \
    gifsicle \
    optipng \
    jpegoptim \
    libjpeg-turbo \
    libwebp-tools

RUN npm install -g --no-progress \
    sass-embedded \
    terser \
    typescript \
    svgo \
    html-minifier \
    rollup

COPY --from=builder /usr/local/bin/compactor /usr/local/bin/

WORKDIR /app
ENTRYPOINT ["compactor"]
CMD ["--help"]
