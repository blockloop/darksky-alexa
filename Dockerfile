FROM golang:1.10
ENV TZ=America/Chicago
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /go/src/github.com/blockloop/darksky-alexa
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
EXPOSE 3000
RUN mkdir /app
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ADD CHECKS /app/CHECKS
COPY --from=0 /go/src/github.com/blockloop/darksky-alexa/app .
CMD ["./app"]
