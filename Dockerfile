FROM golang:1.10
WORKDIR /go/src/github.com/blockloop/darksky-alexa
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM ubuntu:bionic
RUN apt-get update && apt-get install -y ca-certificates sqlite3
EXPOSE 3000
RUN mkdir /app
WORKDIR /root/
ADD CHECKS /app/CHECKS
ADD geodb.sqlite /app/geodb.sqlite
COPY --from=0 /go/src/github.com/blockloop/darksky-alexa/app .
CMD ["./app"]
