FROM --platform=$BUILDPLATFORM golang:1.21.1 AS build

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build ./cmd/mermaid2svg

FROM --platform=$BUILDPLATFORM ubuntu AS fonts

RUN apt update && \
    apt install -y software-properties-common && \
    add-apt-repository multiverse && \
    apt update && \
    yes yes | DEBIAN_FRONTEND=dialog apt install -y ttf-mscorefonts-installer

FROM selenium/standalone-firefox:117.0

EXPOSE 8080
WORKDIR /app

COPY --from=build /app/mermaid2svg .
COPY --from=fonts /usr/share/fonts/. /usr/share/fonts/

RUN sudo fc-cache -f

ENTRYPOINT ["/app/mermaid2svg"]
