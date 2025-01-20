FROM golang:1.23.2-alpine AS app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main .

FROM golang:1.23.2-alpine AS runtime

WORKDIR /runtime

COPY --from=app /app/main /runtime/main

ENTRYPOINT ["/runtime/main"]