FROM golang:1.15-buster AS builder

ENV CGO_ENABLED=1

WORKDIR /build
COPY . .

RUN go get -d -v ./.. \
    && go install -v ./..

# Build the application
RUN go build -o /dist/6d6-bot

WORKDIR /dist
CMD ["/dist/6d6-bot"]
