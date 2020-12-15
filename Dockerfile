FROM golang:1.15 as builder
WORKDIR /workspace/app
COPY . .
RUN go mod tidy \
    && go get -u -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o bookshelf

FROM scratch
WORKDIR /
COPY --from=builder /workspace/app/bookshelf .
EXPOSE 8080
ENTRYPOINT ["/bookshelf"]
