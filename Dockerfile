FROM golang:alpine AS build-backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o majestic_rss

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-backend /app/majestic_rss ./
CMD ["./majestic_rss"]