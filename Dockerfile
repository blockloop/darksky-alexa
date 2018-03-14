FROM golang:1.10 as build
WORKDIR /go/src/github.com/blockloop/darksky-alexa
ADD . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o service .

FROM ubuntu:bionic
# disable tzdata prompts when installing with apt
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y ca-certificates sqlite3 tzdata
EXPOSE 3000
RUN mkdir /app
WORKDIR /app
ADD CHECKS /app/CHECKS
ADD geodb.sqlite /app/geodb.sqlite
COPY --from=build /go/src/github.com/blockloop/darksky-alexa/service /app/service
CMD ["/app/service"]
