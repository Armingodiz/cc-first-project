FROM golang:1.18-alpine AS go-build

ENV GO111MODULE=on

WORKDIR /app

COPY ./user-service/go.mod .
COPY ./user-service/go.sum .

RUN go mod download

COPY ./core .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service .

FROM alpine:3.9.5 as dns
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-build /app .
ENTRYPOINT ["./user-service"]