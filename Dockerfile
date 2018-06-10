FROM golang:1.10-alpine AS builder

RUN apk --no-cache add git

ENV GOPATH /go
ENV SRCPATH $GOPATH/src/github.com/factorycampus/radau

COPY main.go ${SRCPATH}/
COPY api ${SRCPATH}/api
COPY migrations ${SRCPATH}/migrations

WORKDIR $SRCPATH
RUN go get -v -d

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir -p /tmp/build && go build -o /tmp/build/radau


FROM scratch
ENV GIN_MODE=release
ENV PORT 8080
EXPOSE 8080

COPY --from=builder /tmp/build/radau radau

ENTRYPOINT [ "./radau" ]
