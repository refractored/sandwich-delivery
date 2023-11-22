FROM golang:1.21-alpine AS build
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download
COPY . /build
RUN CGO_ENABLED=0 go build -o /build/main src/main/main.go


FROM alpine:latest
COPY --from=build /build/main /app/main
WORKDIR /app/work
CMD ["/app/main"]