#
# STAGE 1: prepare
#
# Start from golang base image
FROM golang:1.17.1-alpine as build

# Updates the repository and installs git
RUN apk update && apk upgrade && apk --no-cache add git && apk --no-cache add tzdata && rm -rf /var/cache/apk/*


WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
COPY .env .
# RUN go mod download

#
# STAGE 2: build
#
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./api .

#########################################################
#
# STAGE 3: run
#
# The project has been successfully built and we will use a
# lightweight alpine image to run the server
FROM alpine:latest as Dev

# Adds CA Certificates to the image
RUN apk update && apk --no-cache add ca-certificates && apk --no-cache add tzdata && rm -rf /var/cache/apk/*

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /go/src/app/api /app/api
COPY --from=build /go/src/app/.env /app/.env

# Switches working directory to /app
WORKDIR "/app"

# Exposes the 8080 port from the container
EXPOSE 8080

# Runs the binary once the container starts
CMD ["./api"]