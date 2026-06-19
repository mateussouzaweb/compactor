# syntax=docker/dockerfile:1
# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /usr/local/bin/compactor ./cmd/compactor

# Runtime stage
FROM node:20-alpine AS runtime

RUN apk add --no-cache libjpeg-turbo gcompat

RUN npm install -g --no-progress \
    gifsicle \
    jpegoptim-bin \
    cwebp-bin \
    optipng-bin \
    sass-embedded \
    terser \
    typescript \
    svgo \
    html-minifier \
    rollup

COPY --from=builder /usr/local/bin/compactor /usr/local/bin/

WORKDIR /data
ENTRYPOINT ["compactor"]
CMD ["--help"]
