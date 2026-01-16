# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /vaultage .

# Runtime stage - minimal scratch image
FROM scratch

COPY --from=builder /vaultage /vaultage

# Run as non-root (nobody)
USER 65534:65534

ENTRYPOINT ["/vaultage"]
CMD ["watch"]
