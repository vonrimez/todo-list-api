DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;

CREATE TYPE TASK_STATUS AS ENUM ('todo', 'in-progress', 'done');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(50) UNIQUE,
    password VARCHAR(255)
);

CREATE TABLE tasks (
    id BIGINT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, id),
    title VARCHAR(50),
    description VARCHAR(255) NOT NULL,
    status TASK_STATUS,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);