FROM node:22-alpine AS client-builder

WORKDIR /build
COPY client/package.json client/package-lock.json ./
RUN npm ci --ignore-scripts
COPY client/ .

RUN npm run build

FROM golang:1.25-alpine AS go-builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o random-pass \
    ./main.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o healthcheck \
    ./cmd/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=go-builder /build/random-pass .
COPY --from=go-builder /build/healthcheck /healthcheck
COPY --from=client-builder /build/dist ./client/dist

EXPOSE 3000

USER nonroot:nonroot

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD ["/healthcheck"]

ENTRYPOINT ["/app/random-pass"]
