#!/bin/bash

set -e

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}    Docker Hub Publishing Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Проверка зависимостей
echo -e "${YELLOW}Checking dependencies...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    exit 1
fi
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git is not installed${NC}"
    exit 1
fi
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}Warning: jq is not installed (optional for JSON formatting)${NC}"
fi
echo -e "${GREEN}✓ All dependencies found${NC}"
echo ""

# Получение информации
echo -e "${YELLOW}Getting version information...${NC}"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

echo -e "  Version:    ${GREEN}$VERSION${NC}"
echo -e "  Commit:     ${GREEN}$COMMIT${NC}"
echo -e "  Build Time: ${GREEN}$BUILD_TIME${NC}"
echo ""

# Запрос имени пользователя
read -p "Docker Hub username: " DOCKER_USERNAME
read -p "Image name [go-dockerhub]: " IMAGE_NAME
IMAGE_NAME=${IMAGE_NAME:-go-dockerhub}

echo ""
echo -e "${YELLOW}Building image...${NC}"
docker build \
    --build-arg VERSION="$VERSION" \
    --build-arg COMMIT="$COMMIT" \
    --build-arg BUILD_TIME="$BUILD_TIME" \
    -t "$IMAGE_NAME:$VERSION" \
    -t "$IMAGE_NAME:latest" \
    .

echo -e "${GREEN}✓ Build complete${NC}"
echo ""

# Локальное тестирование
echo -e "${YELLOW}Testing image locally...${NC}"
CONTAINER_ID=$(docker run -d -p 8080:8080 "$IMAGE_NAME:$VERSION")
sleep 3

if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "  ${GREEN}✓ Health check passed${NC}"
else
    echo -e "  ${RED}✗ Health check failed${NC}"
    docker logs "$CONTAINER_ID"
    docker stop "$CONTAINER_ID" >/dev/null
    exit 1
fi

VERSION_CHECK=$(curl -s http://localhost:8080/version | grep -o '"Version":"[^"]*"' | cut -d'"' -f4)
if [ "$VERSION_CHECK" = "$VERSION" ]; then
    echo -e "  ${GREEN}✓ Version check passed${NC}"
else
    echo -e "  ${YELLOW}⚠ Version mismatch: expected $VERSION, got $VERSION_CHECK${NC}"
fi

docker stop "$CONTAINER_ID" >/dev/null
echo -e "${GREEN}✓ Test complete${NC}"
echo ""

# Тегирование
echo -e "${YELLOW}Tagging for Docker Hub...${NC}"
docker tag "$IMAGE_NAME:$VERSION" "$DOCKER_USERNAME/$IMAGE_NAME:$VERSION"
docker tag "$IMAGE_NAME:latest" "$DOCKER_USERNAME/$IMAGE_NAME:latest"
docker tag "$IMAGE_NAME:$VERSION" "$DOCKER_USERNAME/$IMAGE_NAME:$COMMIT"

echo -e "  ${GREEN}✓ Tags created:${NC}"
echo -e "    $DOCKER_USERNAME/$IMAGE_NAME:$VERSION"
echo -e "    $DOCKER_USERNAME/$IMAGE_NAME:latest"
echo -e "    $DOCKER_USERNAME/$IMAGE_NAME:$COMMIT"
echo ""

# Публикация
echo -e "${YELLOW}Publishing to Docker Hub...${NC}"
docker push "$DOCKER_USERNAME/$IMAGE_NAME:$VERSION"
docker push "$DOCKER_USERNAME/$IMAGE_NAME:latest"
docker push "$DOCKER_USERNAME/$IMAGE_NAME:$COMMIT"

echo -e "${GREEN}✓ Published to Docker Hub${NC}"
echo -e "  https://hub.docker.com/r/$DOCKER_USERNAME/$IMAGE_NAME"
echo ""

# Очистка
read -p "Clean up local images? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Cleaning up...${NC}"
    docker rmi "$IMAGE_NAME:$VERSION" 2>/dev/null
    docker rmi "$IMAGE_NAME:latest" 2>/dev/null
    docker rmi "$DOCKER_USERNAME/$IMAGE_NAME:$VERSION" 2>/dev/null
    docker rmi "$DOCKER_USERNAME/$IMAGE_NAME:latest" 2>/dev/null
    echo -e "${GREEN}✓ Clean complete${NC}"
fi

echo ""
echo -e "${GREEN}✓ Script completed successfully!${NC}"