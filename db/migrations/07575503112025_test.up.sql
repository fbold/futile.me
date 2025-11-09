CREATE TABLE users (
	id int PRIMARY KEY,
	username varchar(20) UNIQUE NOT NULL,
	email varchar(32) UNIQUE NOT NULL,
	password varchar(256) NOT NULL,
);
