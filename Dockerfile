FROM golang:alpine as builder
COPY . /src
WORKDIR /src
RUN go build

FROM alpine
COPY --from=builder /src/releasebot /app/
ENTRYPOINT [/app/releasebot]