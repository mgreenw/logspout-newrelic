# ## Multi-stage build

#
# Init stage, includes logspout source code
# and triggers the build.sh script
#
FROM gliderlabs/logspout:v3.2.13 as logspout

#
# Build stage, build logspout with fluentd adapter
#
FROM golang:1.12.5-alpine3.9 as builder
RUN apk add --update go build-base git mercurial ca-certificates git
ENV GO111MODULE=on
WORKDIR /go/src/github.com/gliderlabs/logspout
COPY --from=logspout /go/src/github.com/gliderlabs/logspout /go/src/github.com/gliderlabs/logspout
COPY modules.go .
ADD . /go/src/github.com/mgreenw/logspout-newrelic
RUN cd /go/src/github.com/mgreenw/logspout-newrelic; go mod download
RUN cd /go/src/github.com/gliderlabs/logspout; go mod download
RUN echo "replace github.com/mgreenw/logspout-newrelic => /go/src/github.com/mgreenw/logspout-newrelic" >> go.mod

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$1" -o /bin/logspout


# #
# # Final stage
# #
FROM alpine
WORKDIR /app
COPY --from=builder /bin/logspout /app/
CMD ["./logspout"]
