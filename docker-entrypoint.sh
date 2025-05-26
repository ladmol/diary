#!/bin/sh

# Ожидание доступности PostgreSQL
echo "Ожидание доступности PostgreSQL..."
until pg_isready -h postgres -p 5432 -U postgres; do
  echo "PostgreSQL недоступен - ожидание..."
  sleep 2
done
echo "PostgreSQL готов!"

# Запуск основного приложения
echo "Запуск основного приложения..."
exec "$@"
