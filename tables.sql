CREATE TABLE IF NOT EXISTS pages (
	id			SERIAL PRIMARY KEY,
	title		TEXT,
	slug		TEXT,
	sort		INTEGER
);

CREATE TABLE IF NOT EXISTS users (
	id			SERIAL PRIMARY KEY,
	name		TEXT,
	email		TEXT,
	password	TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
	user_id 	INTEGER REFERENCES users,
	key 		TEXT,
	created_ip 	TEXT
);
CREATE INDEX IF NOT EXISTS session_index ON sessions(user_id);

CREATE TABLE IF NOT EXISTS components (
	id			SERIAL PRIMARY KEY,
	page_id 	INTEGER REFERENCES pages,
	sort		INTEGER,
	data		TEXT
);