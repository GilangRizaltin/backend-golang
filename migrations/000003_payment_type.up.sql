CREATE TABLE "golang_db".payment_type (
	id serial4 NOT NULL,
	payment_name varchar(50) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	update_at timestamp NULL,
	CONSTRAINT pk_payment_type PRIMARY KEY (id)
);