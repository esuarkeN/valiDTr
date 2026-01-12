# --- build stage ---
FROM golang:1.23-bookworm AS builder
WORKDIR /src

# leverage module cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# adjust if your main package path differs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/valiDTr ./cmd/valiDTr

# --- runtime stage ---
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates git gnupg \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/valiDTr /usr/local/bin/valiDTr
ENTRYPOINT ["valiDTr"]
