FROM golang:1.25-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build build-api

FROM debian:bookworm-slim
COPY --from=builder /app/bin/solarsage-mcp /usr/local/bin/
COPY --from=builder /app/bin/solarsage-api /usr/local/bin/
COPY --from=builder /app/third_party/swisseph/ephe /usr/local/share/swisseph/ephe

ENV SWISSEPH_EPHE_PATH=/usr/local/share/swisseph/ephe
ENTRYPOINT ["solarsage-mcp"]
