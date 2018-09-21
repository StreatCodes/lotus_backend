INSERT INTO users (name, email, password)
	VALUES (
		'Lotus',
		'lotus@example.com',
		'$2a$10$SOzM4WbuuDOIgD5EoXfylOoLN0DakO40/YmMQEDDDifY2AGMfjDG6'
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Lotus',
		'',
		null,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Locations',
		'locations',
		null,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Sydney',
		'sydney',
		2,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Chatswood',
		'chatswood',
		3,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Circular Quay',
		'circular-quay',
		3,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Melbourne',
		'melbourne',
		2,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'Perth',
		'perth',
		2,
		0
	);

INSERT INTO pages (title, slug, parent, sort)
	VALUES (
		'About',
		'about',
		2,
		0
	);


INSERT INTO components (page_id, sort, data) 
	VALUES (
		1,
		1,
		'{"width": 2, "title": "sample text", "text": "Hello world!"}'
	);