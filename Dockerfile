FROM golang:1.14-alpine as build-container

ENV PROJECT_DIR=${GOPATH}/src/github.com/oleg-balunenko/logs-converter

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git musl-dev make gcc

RUN mkdir -p ${PROJECT_DIR}

COPY ./  ${PROJECT_DIR}
WORKDIR ${PROJECT_DIR}
# check vendor
#RUN make dependencies
# vet project
# RUN make vet
# test project
RUN make test
# compile executable
RUN make compile

RUN mkdir /app
RUN cp ./bin/logs-converter /app/logs-converter_unix


FROM alpine:3.11.3 as deployment-container
RUN apk add -U --no-cache ca-certificates

COPY ./testdata  /testdata 


COPY --from=build-container /app/logs-converter_unix /logs-converter_unix

ENTRYPOINT ["/logs-converter_unix"]

# Expose port
EXPOSE $APP_PORT $HEALTH_PORT
