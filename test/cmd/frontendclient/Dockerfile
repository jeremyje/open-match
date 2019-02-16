# Golang application builder steps
FROM alpine:3.8

ARG APP_ROOT=/go/src/github.com/GoogleCloudPlatform/open-match
ARG APP_NAME=frontendclient

RUN apk --update add ca-certificates \
  && adduser -D openmatch \
  && mkdir -p ${APP_ROOT}/config \
  && mkdir -p ${APP_ROOT}/cmd/${APP_NAME} \
  && chown -R openmatch:openmatch /go

COPY --chown=openmatch:root ./test/cmd/${APP_NAME}/${APP_NAME} ${APP_ROOT}/test/cmd/${APP_NAME}/${APP_NAME}
COPY --chown=openmatch:root ./test/cmd/${APP_NAME}/. ${APP_ROOT}/.

USER openmatch
WORKDIR ${APP_ROOT}
ENTRYPOINT ["/go/src/github.com/GoogleCloudPlatform/open-match/test/cmd/frontendclient/frontendclient"]
