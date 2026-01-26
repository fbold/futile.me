CREATE TABLE documents (
	id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id int NOT NULL,
	content text NOT NULL,
	private boolean DEFAULT TRUE,

	CONSTRAINT fk_user
		FOREIGN KEY(user_id)
			REFERENCES users(id)
)
