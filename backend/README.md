# Запуск

### Первый раз

```
cp .env.example .env && task app:build
```

### Последующие разы

```
task app:build-up
```

# Остановить

```
task app:build-down
```

# Файловое хранилище

Аватары хранятся в SeaweedFS. Если приложение запускается локально (`task app:run`), хранилище нужно поднять отдельно:

```
task storage:up
```
