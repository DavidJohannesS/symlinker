#!/bin/bash
set -e

GREEN="\033[38;5;82m"
YELLOW="\033[38;5;226m"
RESET="\033[0m"
IMAGE_NAME="symlinker-builder"
BINARY_NAME="symlinker"
OUTPUT_DIR="./build"

echo "Building the Docker image..."
docker build -t $IMAGE_NAME .

echo "Extracting the binary..."
CONTAINER_ID=$(docker create $IMAGE_NAME)

mkdir -p $OUTPUT_DIR

docker cp $CONTAINER_ID:/app/bin/$BINARY_NAME $OUTPUT_DIR/$BINARY_NAME

docker rm $CONTAINER_ID > /dev/null
chmod +x $OUTPUT_DIR/$BINARY_NAME

echo -e "${GREEN}Build complete! Your binary is ready at: ${YELLOW}$OUTPUT_DIR/$BINARY_NAME${RESET}"
