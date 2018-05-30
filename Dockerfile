FROM golang:1.10-alpine AS builder

RUN apk --no-cache add git

ENV SRCPATH /src
COPY main.go $SRCPATH/
COPY internal $SRCPATH/internal

WORKDIR $SRCPATH
RUN go get -v -d

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir build && go build -o build/radau


FROM scratch
ENV GIN_MODE=release
ENV PORT 8080
EXPOSE 8080

COPY --from=builder /src/build/radau radau

CMD ["./wifilogin"]
