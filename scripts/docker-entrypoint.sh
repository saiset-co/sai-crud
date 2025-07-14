#!/bin/sh
set -e

echo "Starting SAI Service container initialization..."

CONFIG_TEMPLATE="config.yaml.template"

echo "Using config template: $CONFIG_TEMPLATE"

if [ ! -f "./$CONFIG_TEMPLATE" ]; then
    echo "Error: Config template ./$CONFIG_TEMPLATE not found!"
    exit 1
fi

echo "Processing configuration template with environment variables..."
envsubst < "./$CONFIG_TEMPLATE" > "./config.yaml"

echo "Configuration file generated successfully:"
echo "--- Generated config.yaml ---"
cat "./config.yaml"
echo "--- End of config ---"

echo "Environment validation passed."
echo "Starting application with command: $@"

exec "$@"