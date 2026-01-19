CREATE TABLE users (
	id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	username text NOT NULL UNIQUE,
	email text NULL UNIQUE,
	password text NOT NULL
	-- recovery text NOT NULL UNIQUE
)
