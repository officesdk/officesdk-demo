
set -e

export CGO_ENABLED=0
export GOOS=$TARGETOS
export GOARCH=$TARGETARCH
#export GOARCH=amd64
export APP=office-demo

apk add --no-cach bash make git && make build

mkdir ./data/ && mkdir ./data/logs/ && mkdir ./data/files/ && mkdir ./data/leveldb/ && mkdir ./data/resource && mkdir ./data/config
mv ./bin/*  ./data/ && cp -r ./config/ ./data/config/ && cp -r ./resource ./data/resource

echo "./${APP} server" >> ./data/start.sh