
cd /data/minio && chmod +x ./start.sh && ./start.sh 

cd /data/sdk && chmod +x ./start.sh && ./start.sh 
cd /data/demo && chmod +x ./start.sh && EGO_CONFIG_PATH="./config/default.toml" ./start.sh
