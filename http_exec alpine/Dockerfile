# Start by building the application.
FROM golang:alpine AS builder

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Now copy it into our base image.
FROM alpine:latest 
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app .
CMD ["./app"] 