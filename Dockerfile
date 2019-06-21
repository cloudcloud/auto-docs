FROM node:12.0-slim AS fe
LABEL maintainer="Al <allan.shone@gmail.com>"
COPY . .
RUN yarn && yarn build

FROM golang:alpine AS bin
WORKDIR /ad
COPY . .
COPY --from=fe ["dist/", "dist/"]
RUN apk add ca-certificates git && \
    update-ca-certificates && \
    GO111MODULE=off go get -u github.com/kevinburke/go-bindata/... && \
    go-bindata -o auto-docs/server/assets.go --prefix dist/ dist/... && \
    sed -i "s/package main/package server/" auto-docs/server/assets.go && \
    GO111MODULE=on go mod download && \
    GO111MODULE=on go build -o /auto-docs ./auto-docs

FROM golang:alpine AS release
ENTRYPOINT ["/auto-docs"]
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=bin ["/auto-docs", "/auto-docs"]

