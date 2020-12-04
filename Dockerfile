FROM golang:1.15-buster AS builder

ARG discord_bot_token
ENV CGO_ENABLED=1

WORKDIR /build
COPY . .

ENV DISCORD_BOT_TOKEN=$discord_bot_token

RUN go build -o /dist/6d6-bot

WORKDIR /dist
CMD ["/dist/6d6-bot"]
