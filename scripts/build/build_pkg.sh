
set -e

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export APP=office-demo

apk add --no-cach bash make git && make build

mkdir ./data/ && mkdir ./data/logs/ && mkdir ./data/files/ && mkdir ./data/leveldb/ && mkdir ./data/resource
mv ./bin/*  ./data/ && cp ./config/default.toml  ./data/ && cp -r ./resource ./data/resource

echo "EGO_CONFIG_PATH=./default.toml ./${APP} server -c ./default.toml" >> ./data/start.sh