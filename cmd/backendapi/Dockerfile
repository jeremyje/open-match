# Golang application builder steps
FROM alpine:3.8

ARG APP_ROOT=/go/src/github.com/GoogleCloudPlatform/open-match
ARG APP_NAME=backendapi

RUN apk --update add ca-certificates \
  && adduser -D openmatch \
  && mkdir -p ${APP_ROOT}/config \
  && mkdir -p ${APP_ROOT}/cmd/${APP_NAME} \
  && chown -R openmatch:openmatch /go

COPY --chown=openmatch:root ./cmd/${APP_NAME}/${APP_NAME} ${APP_ROOT}/cmd/${APP_NAME}/${APP_NAME}

USER openmatch
WORKDIR ${APP_ROOT}
ENTRYPOINT ["/go/src/github.com/GoogleCloudPlatform/open-match/cmd/backendapi/backendapi"]
