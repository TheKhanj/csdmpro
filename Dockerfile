FROM golang:alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/google/wire/cmd/wire@latest
ARG APK_MIRROR
RUN if [ -n "$APK_MIRROR" ]; then \
        sed -i "s|https://dl-cdn.alpinelinux.org/alpine|$APK_MIRROR|g" /etc/apk/repositories; \
    fi
RUN apk update && apk add --no-cache make gcc musl-dev
COPY . .
RUN CGO_ENABLED=1 make

FROM alpine
ARG APK_MIRROR
RUN if [ -n "$APK_MIRROR" ]; then \
        sed -i "s|https://dl-cdn.alpinelinux.org/alpine|$APK_MIRROR|g" /etc/apk/repositories; \
    fi
RUN apk update && apk add --no-cache musl
COPY --from=build /src/csdmpro /usr/local/bin/
RUN mkdir /csdmpro
WORKDIR /csdmpro
ENTRYPOINT ["csdmpro"]
