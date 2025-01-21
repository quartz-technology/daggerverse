ARG BASE_IMAGE=golang:1.23.2-alpine

FROM ${BASE_IMAGE} AS app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main .

FROM golang:1.23.2-alpine AS runtime

ARG BIN_NAME

WORKDIR /runtime

COPY --from=app /app/main /runtime/${BIN_NAME}

RUN --mount=type=secret,id=my-super-secret \
  cat /run/secrets/my-super-secret > /runtime/secret.txt

ENTRYPOINT ["/runtime/${BIN_NAME}"]