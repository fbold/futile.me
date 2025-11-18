CREATE TABLE users (
	id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	username text NOT NULL UNIQUE,
	email text NOT NULL UNIQUE,
	password text NOT NULL UNIQUE,
	recovery text NOT NULL UNIQUE
)
