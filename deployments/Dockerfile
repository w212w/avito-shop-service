FROM golang:1.23

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY . .

RUN go build -o avito-shop-service ./cmd/shop-service \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/app/avito-shop-service"]
