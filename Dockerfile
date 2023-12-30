FROM golang:1.21.3-alpine as builder
WORKDIR /app
RUN apk update && apk add --no-cache gcc musl-dev git
COPY go.* ./
RUN go mod download
COPY . .
WORKDIR /app
RUN go build -ldflags '-w -s' -a -o mongo-streamer ./main.go

# Deployment environment
# ----------------------
FROM alpine:3.18.4
WORKDIR /app
RUN chown nobody:nobody /app
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody ./app/mongo-streamer .
COPY --from=builder --chown=nobody:nobody ./app/run.sh .

ENTRYPOINT sh run.sh
