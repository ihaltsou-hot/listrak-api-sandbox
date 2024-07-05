FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update && apk add ca-certificates && apk add tzdata

COPY . .
RUN go build -o main .

FROM scratch

ENV PORT=8000

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/main

EXPOSE 8000

ENTRYPOINT ["/app/main"]