FROM golang:1.15.3-alpine3.12 as builder
COPY . /src
WORKDIR /src
RUN go build

FROM alpine
COPY --from=builder /src/releasebot /app/
ENTRYPOINT [/app/releasebot]