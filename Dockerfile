FROM golang:1.19.0-alpine3.16 as build

RUN mkdir -p /opt/app
WORKDIR /opt/app
COPY . .

RUN apk add --no-cache gcc musl-dev
RUN go test -v ./...
RUN go build ./cmd/server

FROM alpine:3.16 as release

RUN mkdir -p /opt/app
COPY --from=build /opt/app/server /opt/app/server
WORKDIR /opt/app

ENTRYPOINT ["./server"]
