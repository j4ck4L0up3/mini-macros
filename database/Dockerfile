FROM postgres:16-alpine

ARG INFISICAL_VERSION=0.28.4

RUN apk add --no-cache bash sudo
RUN curl -1sLf \
  'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' \
  | sudo -E distro=alpine bash
RUN apk update && sudo apk add infisical=$INFISICAL_VERSION


