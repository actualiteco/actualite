# actualite - frontend

SPA: React 19 + TypeScript + Zustand + Webpack. 

## Запуск

```bash
npm install
API_URL=http://127.0.0.1:8080 npm run start
npm run typecheck
npm run build:prod
```

`API_URL` задаётся переменной окружения (NF_Ops_1), по умолчанию — `http://127.0.0.1:8080`.
В dev-режиме клиент ходит на `/api/*`, а dev-сервер проксирует это на `API_URL` с отрезанием `/api`,
поэтому CORS на бэкенде не нужен. В prod-сборке клиент обращается на `API_URL` напрямую (тогда бэкенду
нужен CORS) или на `/api`, если его проксирует веб-сервер. Пример — в `.env.example`.

Типы API лежат в `src/generated/data-contracts.ts` и соответствуют `swagger.json`. Пересоздать:
`npm run generate-api`. В коде используются короткие алиасы из `src/models/api.types.ts`

## Как устроена сессия

- Вход возвращает пару токенов: access (10 минут) и refresh. Они хранятся в `localStorage` (`access_token`, `refresh_token`).
- При 401 на защищённом запросе клиент один раз обновляет токены через `POST /auth/refresh` и повторяет запрос.
  Параллельные запросы ждут один общий refresh. Если refresh недействителен - сессия сбрасывается один раз
  (одно уведомление), роутер уводит на `/auth`.
- Выход вызывает `POST /auth/logout` (инвалидирует refresh-токен) и чистит локальную сессию.
- Данные пользователя (id, никнейм, аватар) берутся из `GET /profile/me`, т.к. ответ входа их не содержит.
