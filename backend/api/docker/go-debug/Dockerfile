# Compile stage
FROM golang:1.23 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.1

ADD . /dockerdev
WORKDIR /dockerdev

# Compile the application with the optimizations turned off
# This is important for the debugger to correctly work with the binary
RUN go build -gcflags "all=-N -l" -o /server

# Final stage
FROM debian:buster

RUN apt-get update -y
RUN apt-get install -y libc6

EXPOSE 8000 40000

WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server /

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]
