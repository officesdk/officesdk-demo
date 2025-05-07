#!/bin/sh
set -ex

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd $SCRIPT_DIR

curl -LO https://dl.min.io/server/minio/release/${TARGETOS}-${TARGETARCH}/minio 

curl -LO https://dl.min.io/client/mc/release/${TARGETOS}-${TARGETARCH}/mc

mkdir -p /var/minio-data

chmod +x ./minio && ./minio -version
chmod +x ./mc && ./mc -version

cat <<EOF > ./start.sh
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin"
chmod +x ./minio && nohup ./minio server /var/minio-data > /var/log/minio.log 2>&1 &
sleep 5
chmod +x ./mc && ./mc alias set local http://localhost:9000 minioadmin minioadmin && ./mc mb local/officesdk
EOF
