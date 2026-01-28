FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /vaultage .


FROM scratch AS final

COPY --from=builder /vaultage /vaultage

USER 65534:65534

ENTRYPOINT ["/vaultage"]
CMD ["watch"]
