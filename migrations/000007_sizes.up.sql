CREATE TABLE "golang_db".sizes (
	id serial4 NOT NULL,
	size_name varchar(50) NOT NULL,
	additional_fee int4 NOT NULL DEFAULT 0,
	created_at timestamp NOT NULL DEFAULT now(),
	update_at timestamp NULL,
	CONSTRAINT pk_size PRIMARY KEY (id)
);