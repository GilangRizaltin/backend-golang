CREATE TABLE "golang_db".products (
	id serial4 NOT NULL,
	product_name varchar(100) NOT NULL,
	category int4 NOT NULL,
	description text NOT NULL,
	price_default int4 NOT NULL DEFAULT 0,
	created_at timestamp NOT NULL DEFAULT now(),
	update_at timestamp NULL,
	deleted_at timestamp NULL,
	product_image_1 text NULL,
	product_image_2 text NULL,
	product_image_3 text NULL,
	product_image_4 text NULL,
	CONSTRAINT pk_products PRIMARY KEY (id),
	CONSTRAINT unique_product_name UNIQUE (product_name)
);

ALTER TABLE "golang_db".products ADD CONSTRAINT products_category_fkey FOREIGN KEY (category) REFERENCES "golang_db".categories(id);