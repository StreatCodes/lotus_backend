CREATE TABLE IF NOT EXISTS pages (
	id				SERIAL PRIMARY KEY,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	title			TEXT,
	slug			TEXT,
	sort			INTEGER
);

CREATE TABLE IF NOT EXISTS users (
	id				SERIAL PRIMARY KEY,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	name			TEXT,
	email			TEXT,
	password		TEXT
);

CREATE TABLE IF NOT EXISTS permissions (
	user_id 		INTEGER REFERENCES users,
	manage_users 	BOOL,
	view_all_logs	BOOL,
	manage_pages	BOOL,
	manage_media	BOOL
);
CREATE INDEX ON permissions(user_id);

CREATE TABLE IF NOT EXISTS sessions (
	user_id 		INTEGER REFERENCES users,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	created_ip 		TEXT,
	key 			TEXT
);
CREATE INDEX ON sessions(user_id);

CREATE TABLE IF NOT EXISTS components (
	id				SERIAL PRIMARY KEY,
	page_id 		INTEGER REFERENCES pages,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	sort			INTEGER,
	data			TEXT
);
CREATE INDEX ON components(page_id, created_at);
-- CREATE INDEX ON components(created_at);

CREATE TABLE IF NOT EXISTS logs (
	id				SERIAL PRIMARY KEY,
	user_id 		INTEGER REFERENCES users,
	created_at		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	action			TEXT
);
CREATE INDEX ON logs(user_id);