# Client builder
FROM node:22-slim AS client-builder

WORKDIR /src
COPY client ./client

WORKDIR /src/client
RUN npm ci
RUN npm run all

# Server builder
FROM golang:1.26 AS server-builder

WORKDIR /src
COPY server ./server

WORKDIR /src/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vordl

# Final image
FROM scratch

WORKDIR /app

COPY --from=client-builder /src/client/dist ./static
COPY --from=server-builder /src/server/vordl ./

ENTRYPOINT ["/app/vordl"]
