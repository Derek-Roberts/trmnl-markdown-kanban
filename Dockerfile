# Build stage: compile static Go binary
FROM golang:1.24.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk add --no-cache git && \
    CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Final stage: distroless container
FROM gcr.io/distroless/static-debian12
# Set a working directory so relative paths resolve
WORKDIR /app
# Copy the server binary
COPY --from=build /server /app/server
# Copy your templates into the same image
COPY --from=build /src/templates /app/templates
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
