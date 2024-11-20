#!/bin/bash

# Check if domain is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <domain> [build_params]"
    echo "Example: $0 example.com '--param1 value1 --param2 value2'"
    exit 1
fi

DOMAIN=$1
BUILD_PARAMS=${2:-""}

# Determine if we need to deploy registry
if [ -z "$REGISTRY" ]; then
    REGISTRY_DEPLOY="true"
else
    REGISTRY_DEPLOY="false"
fi

# Install/upgrade the Helm chart
helm upgrade --install studio-release . \
  --namespace studio \
  --create-namespace \
  --set global.domain=$DOMAIN \
  --set global.registry.deploy=$REGISTRY_DEPLOY \
  --set postgresql.enabled=true \
  --set studio.image.build.params="$BUILD_PARAMS"