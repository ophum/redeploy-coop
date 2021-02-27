FROM golang:1.15 as builder
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY ./main /usr/local/bin/agent
