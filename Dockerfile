# STEP 1 build executable binary

FROM golang:alpine as builder

# Install git
RUN apk update && apk add git 

COPY . $GOPATH/src/github.com/alex1ai/infogration
WORKDIR $GOPATH/src/github.com/alex1ai/infogration

# get dependancies
RUN go get -d -v

#build the binary
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /go/bin/main

# STEP 2 build a small image
# This is needed as everything else would end up in a huge container (>250MB)
# start from scratch
FROM scratch

# Copy our static executable
COPY --from=builder /go/bin/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]
