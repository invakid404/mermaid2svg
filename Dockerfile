FROM golang:1.20.5 as build

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build ./cmd/mermaid2svg

FROM selenium/standalone-firefox:112.0

EXPOSE 8080
WORKDIR /app

COPY --from=build /app/mermaid2svg .

ENTRYPOINT ["/app/mermaid2svg"]
