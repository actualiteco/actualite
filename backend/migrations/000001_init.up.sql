CREATE TABLE users (
  id          VARCHAR(36)     PRIMARY KEY,
  email       VARCHAR(100)    NOT NULL      UNIQUE CHECK (email ~ '^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\.[a-zA-Z]{2,}$'),
  password    TEXT,
  role        VARCHAR(15)     NOT NULL      CHECK (role IN ('admin', 'user')),
  created_at  TIMESTAMP                     DEFAULT NOW(),
  username    VARCHAR(15)     NOT NULL      UNIQUE CHECK (char_length(username) >= 4),
  avatar_url  TEXT,
  bio         TEXT
);

INSERT INTO users(id, email, password, role, username)
VALUES('8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'test.user@email.com', '$2a$14$ZHHlsHyTMgDQMNPa5uMMCO61rbFJk8ax5llJrZMLn8tNSf4gB9k4q', 'user', 'username');

INSERT INTO users(id, email, password, role, username)
VALUES('249be7cf-d419-4c54-97f2-d04107806e36', 'test.admin@email.com', '$2a$14$ZHHlsHyTMgDQMNPa5uMMCO61rbFJk8ax5llJrZMLn8tNSf4gB9k4q', 'admin', 'adminame');