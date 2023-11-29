CREATE TABLE "golang_db".jwt IF NOT EXISTS (
	id serial4 NOT NULL,
	jwt_code text NOT NULL,
	CONSTRAINT pk_jwt PRIMARY KEY (id)
);