
ARG GO_VER
ARG ALPINE_VER

FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder
LABEL stage=server-intermediate
ARG WORKDIR
WORKDIR ${WORKDIR}
RUN apk --no-cache --update add \
    git
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -mod=vendor -o /tmp/recipe-count-task ./cmd/cli/*.go

FROM alpine:${ALPINE_VER} as base
COPY --from=builder /tmp/recipe-count-task /usr/local/bin/recipe-count-task
ENTRYPOINT ["/usr/local/bin/recipe-count-task"]
