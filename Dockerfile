FROM golang:1.16-alpine@sha256:5616dca835fa90ef13a843824ba58394dad356b7d56198fb7c93cbe76d7d67fe AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o piscon_portal main.go

FROM ubuntu:22.04@sha256:4f838adc7181d9039ac795a7d0aba05a9bd9ecd480d294483169c5def983b64d
WORKDIR /app
EXPOSE 4000
RUN apt update \
&& apt install -y tzdata \
&& apt install -y ca-certificates \
&& rm -rf /var/lib/apt/lists/* \
&& update-ca-certificates
COPY --from=build /go/src/github.com/traPtitech/piscon-portal/piscon_portal \
									/go/src/github.com/traPtitech/piscon-portal/.env  ./
ENTRYPOINT ["/app/piscon_portal"]
