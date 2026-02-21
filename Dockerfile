FROM node:20-alpine AS web-build

WORKDIR /webapp
COPY webapp/package.json webapp/package-lock.json ./
RUN npm ci
COPY webapp/ ./
RUN npm run build


FROM golang:1.22-alpine AS go-build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /out/asd-server ./cmd/asd-server
RUN CGO_ENABLED=0 go build -o /out/asd ./cmd/asd


FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates ffmpeg

WORKDIR /app

COPY --from=go-build /out/asd-server /app/asd-server
COPY --from=go-build /out/asd /app/asd

# Embed the built SPA for production usage.
COPY --from=web-build /webapp/dist /app/webapp/dist

ENV ASD_ADDR=0.0.0.0:8080 \
    ASD_DB_PATH=/data/asd.db \
    ASD_WEB_DIST=/app/webapp/dist

EXPOSE 8080

CMD ["/app/asd-server", "-addr", "0.0.0.0:8080", "-db", "/data/asd.db"]
