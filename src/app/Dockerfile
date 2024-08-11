# syntax=docker/dockerfile:1

# get golang base image
FROM golang:1.22-alpine AS build-stage

# run for updates and install node for dependencies
RUN apk -U upgrade &&\
  apk add --no-cache make nodejs npm

# add infisical for env vars
RUN apk add --no-cache bash curl && curl -1sLf \
  'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' | bash &&\
  apk add infisical

# set working directory
WORKDIR /usr/src/app

# get go, templ, air dependencies
COPY go.mod go.sum ./
RUN go install github.com/a-h/templ/cmd/templ@latest &&\
  go install github.com/air-verse/air@latest
RUN go mod download && go mod verify

# setup frontend (htmx, tailwindcss) dependencies
RUN npm install -D tailwindcss &&\
  npm install @tailwindcss/forms &&\
  npm install @tailwindcss/typography

RUN npm install htmx.org@2.0.1

# copy over src code & build app
COPY . ./
RUN make build

# run tests in container
FROM build-stage AS run-test-stage
RUN make test

EXPOSE 4000

CMD ["./bin"]
