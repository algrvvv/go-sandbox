#!/bin/sh

echo "Loading Docker image go-runner"

# Проверяем доступность Docker Daemon
if ! docker info > /dev/null 2>&1; then
    echo "Docker Daemon is not running. Please ensure Docker Daemon is running."
    exit 1
fi

# Загружаем образ go-runner
if ! docker load -i /go-runner.tar; then
    echo "Failed to load Docker image go-runner"
    exit 1
fi

echo "Docker image go-runner loaded successfully"

# Запускаем основное приложение
exec "$@"
