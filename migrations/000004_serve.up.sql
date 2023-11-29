CREATE TABLE "golang_db".serve (
	id serial4 NOT NULL,
	serve_type varchar(50) NOT NULL,
	fee int4 NOT NULL DEFAULT 0,
	created_at timestamp NOT NULL DEFAULT now(),
	update_at timestamp NULL,
	CONSTRAINT pk_serve PRIMARY KEY (id)
);