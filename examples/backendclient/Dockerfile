# Golang application builder steps
FROM alpine:3.8

ARG APP_ROOT=/go/src/github.com/GoogleCloudPlatform/open-match
ARG APP_NAME=backendclient

RUN apk --update add ca-certificates \
  && adduser -D openmatch \
  && mkdir -p ${APP_ROOT}/config \
  && mkdir -p ${APP_ROOT}/cmd/${APP_NAME} \
  && mkdir -p ${APP_ROOT}/profiles \
  && chown -R openmatch:openmatch /go

COPY --chown=openmatch:root ./examples/${APP_NAME}/${APP_NAME} ${APP_ROOT}/examples/${APP_NAME}/${APP_NAME}
COPY --chown=openmatch:root ./examples/${APP_NAME}/profiles/testprofile.json ${APP_ROOT}/profiles/testprofile.json

USER openmatch
WORKDIR ${APP_ROOT}
ENTRYPOINT ["/go/src/github.com/GoogleCloudPlatform/open-match/examples/backendclient/backendclient"]
