FROM golang AS builder

ENV SRCPATH /src
COPY main.go $SRCPATH/
COPY internal $SRCPATH/internal

WORKDIR $SRCPATH
RUN go get -d

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o wifilogin

FROM scratch
ENV GIN_MODE=release
COPY --from=builder /src/wifilogin wifilogin
CMD ["./wifilogin"]
