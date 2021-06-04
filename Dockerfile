#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git build-base
WORKDIR /go/src/app
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
COPY . .
RUN go test ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/TwelveFactorApp /TwelveFactorApp
ENTRYPOINT ./TwelveFactorApp
LABEL Name=twelvefactorapp Version=0.0.1
EXPOSE 8080
