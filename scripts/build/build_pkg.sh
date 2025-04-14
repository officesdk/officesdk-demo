set -e

export CGO_ENABLED=0
export GOOS=$TARGETOS
export GOARCH=$TARGETARCH
#export GOARCH=amd64
export APP=turbo-demo

# Detect system type and install dependencies
if [ -f /etc/alpine-release ]; then
    # Alpine Linux
    apk add --no-cache bash make git
elif [ -f /etc/debian_version ]; then
    # Debian/Ubuntu
    apt-get update && apt-get install -y bash make git
else
    echo "Unsupported system type"
    exit 1
fi

make build

mkdir -p ./data/logs/ ./data/files/ ./data/leveldb/ ./data/resource ./data/config
mv ./bin/*  ./data/ && cp -r ./config/* ./data/config/ && cp -r ./resource/* ./data/resource/

echo "./${APP} server" >> ./data/start.sh