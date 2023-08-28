## Build
FROM golang:1.20-alpine AS buildenv

ADD go.mod go.sum /
RUN go mod download

WORKDIR /app
COPY / /app/
RUN go build -o bin ./cmd/main/

## Deploy
FROM scratch

WORKDIR /

COPY --from=buildenv /app/bin /app/

EXPOSE 8080

CMD ["/app/bin"]