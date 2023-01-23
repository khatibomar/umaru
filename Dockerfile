FROM golang:1.19.5-alpine as build-env

ENV APP_NAME umaru

COPY . /$APP_NAME
WORKDIR /$APP_NAME

RUN CGO_ENABLED=0 go build -v /$APP_NAME

FROM alpine:3.17.1
ENV APP_NAME umaru

ARG GO_PORT=4000
ARG GO_ENV="development"

ENV GO_PORT=${GO_PORT}
ENV GO_ENV=${GO_ENV}

ENV APP_PATH=./${APP_NAME}

COPY --from=build-env /$APP_NAME/$APP_NAME .
COPY --from=build-env /$APP_NAME/templates ./templates
COPY --from=build-env /$APP_NAME/assets ./assets

CMD ${APP_PATH} -port ${GO_PORT} -env ${GO_ENV}