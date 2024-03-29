FROM golang:1.20.5 as build-env
WORKDIR /mbb
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env AS build
WORKDIR /mbb
COPY . .
# make config
RUN apt-get update && apt-get install -y gettext-base
RUN mkdir -p bin/api
RUN env `cat ./build/env.dev` envsubst < ./configs/api-config.json.tmpl > ./bin/api/config.json
# build
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -a -o bin/api ./cmd/mbb-api

FROM alpine
WORKDIR /mbb
COPY --from=build mbb/bin/api .
EXPOSE 5000
CMD ./mbb-api
