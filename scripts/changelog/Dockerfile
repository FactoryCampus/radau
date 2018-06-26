FROM node:10-alpine

RUN apk add -U git

WORKDIR /build
COPY ./* ./
RUN npm install

VOLUME [ "/src" ]
WORKDIR /src
ENTRYPOINT ["/build/node_modules/.bin/conventional-changelog", "-n", "/build/changelog.js"]
