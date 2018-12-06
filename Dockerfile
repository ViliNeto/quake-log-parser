FROM golang:1.11-alpine3.8 as builder

# Set up GOPATH
ENV GOPATH /go

# Install Git
RUN apk update && \
    apk upgrade && \
    apk add git && \
    rm -rf /var/cache/apk/*

# Workdir
WORKDIR /app

# Add current working directory
COPY . /app

# Build
RUN go build main.go

#========================== Runtime Image ==========================
FROM golang:1.11-alpine3.8 as runtime

RUN apk update && \
    apk upgrade && \
    apk add tzdata && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app .

LABEL Name=quake-log-api Version=0.0.1

ENTRYPOINT ["/app/main"]

EXPOSE 9000