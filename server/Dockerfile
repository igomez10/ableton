FROM golang:1.11-alpine as builder
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
WORKDIR /src
COPY ./ ./
RUN go build -installsuffix 'static' -o /server .

FROM alpine
EXPOSE 4444
COPY --from=builder /server /server
ENTRYPOINT ["/server"]

