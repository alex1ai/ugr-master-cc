# Use latest stable golang base container, alpine version because it is much smaller than the normal one
FROM golang:1.11-alpine AS build

# Install tools required for project
RUN apk update
RUN apk add --no-cache git

# Get project from github master
RUN go get github.com/alex1ai/ugr-master-cc

# Switch to project dir and download dependencies
WORKDIR /go/src/github.com/alex1ai/ugr-master-cc
RUN go get -d

# Set environment variables to compile go application to use in scratch below (more below)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /bin/infogration

# Get a totally empty image ("start from scratch")
FROM scratch

# Copy the compiled binary from the intermediate container above to the new scratch container
COPY --from=build /bin/infogration /infogration

# Start server
ENTRYPOINT ["/infogration"]