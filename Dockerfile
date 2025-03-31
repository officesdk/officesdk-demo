## UI build stage
FROM node:18-alpine3.17 AS ui-builder

ARG USE_NEW_FILE_PERMISSION=false
WORKDIR /ui

COPY ui/* ./
COPY ./scripts/build/build_ui.sh ./
RUN chmod +x ./build_ui.sh && ./build_ui.sh


FROM golang:1.21-alpine3.17 AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

ENV TZ Asia/Shanghai
WORKDIR /app

COPY --from=ui-builder /ui/dist ./ui/dist/
COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN chmod +x ./scripts/build/build_pkg.sh && ./scripts/build/build_pkg.sh


FROM alpine:3.16.2 as final

ARG TARGETOS
ARG TARGETARCH

ENV TZ Asia/Shanghai
ENV TRIAL_ENV docker

WORKDIR /data

COPY --from=builder /app/data/* ./

EXPOSE 9001
EXPOSE 8080


ENTRYPOINT ["sh","./start.sh"]
