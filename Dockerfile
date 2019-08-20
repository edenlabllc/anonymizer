# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o anonimizator

# final stage
FROM alpine:3.9
WORKDIR /app
COPY config/ /app/config
COPY --from=build-env /src/anonimizator /app/

ENTRYPOINT ./anonimizator