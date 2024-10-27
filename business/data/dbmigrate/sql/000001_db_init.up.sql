
-- Description: Create table users

CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         citext UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
	created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
	updated_at  TIMESTAMP   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (user_id)
);

-- SEED USER
INSERT INTO users (user_id, name, email, roles, password_hash, department, enabled) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', NULL, true),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', NULL, true)
ON CONFLICT DO NOTHING;

-- ==============================================================================
