FROM golang:alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/google/wire/cmd/wire@latest
RUN apk update
RUN apk add --no-cache make
COPY . .
RUN make

FROM alpine
COPY --from=build /src/csdmpro /usr/local/bin/
ENTRYPOINT ["csdmpro"]
