#!/bin/sh
set -ex

RELEASE_TAG="version-0.0.1-beta38"
REPO_OWNER="officesdk"
REPO_NAME="officesdk"
FILE_NAME="${TARGETOS}-${TARGETARCH}.zip"
#FILE_NAME="${TARGETOS}-amd64.zip"

ROOT_URL="http://127.0.0.1:9101"
CALLBACK_URL="http://127.0.0.1:8080"

DIR=$(cd "$(dirname "$0")"; pwd)

apk add --no-cache curl unzip

curl -sSL -o /tmp/$FILE_NAME \
    "https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$RELEASE_TAG/$FILE_NAME" && \
    unzip /tmp/$FILE_NAME -d $DIR && \
    mkdir $DIR/logs && \
    sed -i "s|^\(rootURL:\).*|\1 \"$ROOT_URL\"|" "$DIR/config/api.yaml" && \
    sed -i "s|^\(\s*endpoint:\s*\).*|\1 \"$CALLBACK_URL\"|" "$DIR/config/callback.yaml" && \
    rm /tmp/$FILE_NAME && \
    ls -ll $DIR 
