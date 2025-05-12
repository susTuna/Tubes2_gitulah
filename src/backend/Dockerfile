#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o /go/bin/ -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/Tubes2_gitulah /app
ENTRYPOINT [ "/app" ]
LABEL Name=tubes2-be-gitulah Version=1.0.0
EXPOSE 5761
