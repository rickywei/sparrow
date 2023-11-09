# build
FROM golang:1.20.4-alpine3.17
COPY . /project
WORKDIR /project
RUN set -eux
RUN apk add make
RUN make build

# deploy
FROM alpine:latest
RUN set -eux
RUN apk add tzdata
# ENV TZ=Asia/Shanghai
COPY --from=0 /project/build /app
COPY --from=0 /project/conf.yaml /app/conf.yaml
COPY --from=0 /project/active.en.toml /app/active.en.toml
COPY --from=0 /project/active.zh.toml /app/active.zh.toml
WORKDIR /app
EXPOSE 80
CMD ["./app"]
