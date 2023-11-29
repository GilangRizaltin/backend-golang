CREATE TABLE "golang_db".orders_products (
	id serial4 NOT NULL,
	order_id int4 NOT NULL,
	product_id int4 NOT NULL,
	size_id int4 NOT NULL,
	hot_or_not bool NOT NULL,
	price int4 NOT NULL,
	quantity int4 NOT NULL,
	subtotal int4 NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_ad timestamp NULL,
	CONSTRAINT pk_orders_products PRIMARY KEY (id)
);

ALTER TABLE "golang_db".orders_products ADD CONSTRAINT orders_products_order_id_fkey FOREIGN KEY (order_id) REFERENCES "golang_db".orders(id);
ALTER TABLE "golang_db".orders_products ADD CONSTRAINT orders_products_product_id_fkey FOREIGN KEY (product_id) REFERENCES "golang_db".products(id);
ALTER TABLE "golang_db".orders_products ADD CONSTRAINT orders_products_size_id_fkey FOREIGN KEY (size_id) REFERENCES "golang_db".sizes(id);