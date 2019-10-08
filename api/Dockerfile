FROM golang:1.13-alpine AS build

RUN apk update && apk add --no-cache git

ADD https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN go get -u github.com/gobuffalo/packr/v2/packr2

WORKDIR $GOPATH/src/github.com/akrantz01/apcsp/api

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure --vendor-only

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr2 build -ldflags="-w -s" -o /go/bin/server .

FROM scratch
VOLUME /config.yaml
COPY --from=build /go/bin/server /
ENTRYPOINT ["/server"]
