ARG BASE=worldmaomao/hd-manager-base:1.0.0

FROM ${BASE} AS builder

WORKDIR $GOPATH/src/hd-manager

ARG MAKE='make build'

COPY . .

RUN $MAKE

FROM scratch

LABEL Name=hd-manager Version=${VERSION}

COPY --from=builder /go/src/hd-manager/cmd /hd-manager
COPY --from=builder /go/src/hd-manager/cmd/res/configuration.toml /hd-manager/res/configuration.toml

EXPOSE 48077

WORKDIR /hd-manager/

CMD ["./hd-manager"]
