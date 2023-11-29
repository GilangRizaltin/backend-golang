CREATE TABLE "golang_db".promos_type IF NOT EXISTS (
	id serial4 NOT NULL,
	promo_type_name varchar(55) NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	CONSTRAINT pk_promos_type PRIMARY KEY (id),
	CONSTRAINT promos_type_promo_type_name_key UNIQUE (promo_type_name)
);