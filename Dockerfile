FROM golang:1.10
WORKDIR /go/src/github.com/blockloop/darksky-alexa
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM ubuntu:bionic
RUN apt-get update && apt-get install -y ca-certificates
ENV TZ=America/Chicago
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
EXPOSE 3000
RUN mkdir /app
WORKDIR /root/
ADD CHECKS /app/CHECKS
COPY --from=0 /go/src/github.com/blockloop/darksky-alexa/app .
CMD ["./app"]
