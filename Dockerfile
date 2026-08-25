FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/admission ./cmd/server
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/admission /app/admission
EXPOSE 8080
ENTRYPOINT ["/app/admission"]
