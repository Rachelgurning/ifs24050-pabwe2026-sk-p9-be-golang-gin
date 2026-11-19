FROM golang:1.27-alpine AS build
WORKDIR /src
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/server /app/server
ENV APP_ENV=production
ENV APP_PORT=8081
EXPOSE 8081
CMD ["/app/server"]
