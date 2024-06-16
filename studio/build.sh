#!/bin/bash

# Variables
BRANCH=$1
IMAGE_TAG=$2
ENV_FILE=${3:-".env"}  # Use the third argument or default to ".env"

# Check if .env file exists
if [ ! -f $ENV_FILE ]; then
  echo ".env file not found!"
  exit 1
fi

# Clone the repository and apply patches
./patch.sh $BRANCH

cd ./code-$BRANCH || exit

if [ -f ../$ENV_FILE ]; then
  if [ -f ./apps/studio/.env ]; then
    rm ./apps/studio/.env
  fi
  cp ../$ENV_FILE ./apps/studio/.env
fi

# Build and push for amd64
docker build \
  --file apps/studio/Dockerfile \
  --target production \
  --tag ${IMAGE_TAG} \
  .

# Cleanup
cd ../../..
# rm -rf ./code-$BRANCH
echo "Produced a patched docker image ${IMAGE_TAG} from branch ${BRANCH}."
