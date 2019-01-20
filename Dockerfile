FROM golang:1.11-alpine AS build

# Install tools required for project
RUN apk update
RUN apk add --no-cache git
RUN go get github.com/alex1ai/ugr-master-cc

WORKDIR /go/src/github.com/alex1ai/ugr-master-cc
RUN go get -d
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /bin/infogration
ENTRYPOINT ["/bin/infogration"]

# This results in a single layer image
FROM scratch
COPY --from=build /bin/infogration /infogration
# Create no-root user to execute the command
#RUN groupadd -r infogration && useradd --no-log-init -r -g infogration infogration
#USER infogration
#ENV MONGO_IP="data"
ENTRYPOINT ["/infogration"]
#CMD ["--help"]