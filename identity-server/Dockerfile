FROM golang:1.17-alpine AS build

WORKDIR /src/identity-server

COPY go.mod ./
COPY go.sum ./

COPY /api/*.go ./api/
COPY /command/*.go ./command/
COPY /config/*.go ./config/
COPY /mail/*.go ./mail/
COPY /model/*.go ./model/
COPY /storage/*.go ./storage/

COPY *.go ./

RUN CGO_ENABLED=0 go build -o server


FROM alpine:3.9 
RUN apk update && apk add bash

WORKDIR /

COPY --from=build /src/identity-server/server /usr/dapper-labs/identity-server/bin/server
COPY config/config.ini /usr/dapper-labs/identity-server/config/config.ini
COPY runInDocker/identity-server-runner.sh /identity-server-runner.sh
COPY runInDocker/wait-for.sh /wait-for.sh
COPY runInDocker/config.ini /usr/dapper-labs/identity-server/config/config.ini
RUN chmod +x /identity-server-runner.sh
RUN chmod +x /wait-for.sh
EXPOSE 8080

# ENTRYPOINT [ "/identity-server-runner.sh" ]


