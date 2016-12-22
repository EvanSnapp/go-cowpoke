FROM golang:1.8

#install gb
RUN go get github.com/constabulary/gb/... && \
    go install github.com/constabulary/gb

#build the project and start the server
COPY ./src /app/src
COPY ./vendor /app/vendor

WORKDIR /app
RUN gb build
ENV GIN_MODE=release
CMD ["bin/cowpoke"]

