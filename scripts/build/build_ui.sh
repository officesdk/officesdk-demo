
export USE_NEW_FILE_PERMISSION=false
export NODE_ENV=production
export CI_COMMIT_SHORT_SHA="v1"

yarn config set registry https://registry.npmmirror.com && yarn install && yarn run build