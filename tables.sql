CREATE TABLE pages (
	id			INTEGER PRIMARY KEY,
	title		TEXT,
	slug		TEXT
);

CREATE TABLE users (
	id			INTEGER PRIMARY KEY,
	name		TEXT,
	email		TEXT,
	password	TEXT
);

CREATE TABLE sessions (
	user_id 	INTEGER REFERENCES users,
	key 		TEXT,
	created_ip 	TEXT
);
CREATE INDEX session_index ON sessions(user_id);

CREATE TABLE components (
	id			INTEGER PRIMARY KEY,
	page_id 	INTEGER REFERENCES pages,
	sort		INTEGER,
	data		STRING
);