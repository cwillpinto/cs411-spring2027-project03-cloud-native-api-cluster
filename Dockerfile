# TODO: use a named Go builder stage, then a minimal non-root runtime image.
FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o /app/api ./src
EXPOSE 8080
ENTRYPOINT ["/app/api"]
