CREATE TABLE "golang_db".orders (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	subtotal int4 NULL,
	promo_id int4 NOT NULL,
	percent_discount float8 NOT NULL DEFAULT 0,
	flat_discount int4 NOT NULL DEFAULT 0,
	serve_id int4 NOT NULL,
	fee int4 NOT NULL,
	tax float8 NOT NULL DEFAULT 0,
	total_transactions int4 NULL,
	payment_type int4 NOT NULL,
	status text NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT orders_status_check CHECK ((status = ANY (ARRAY['On progress'::text, 'Pending'::text, 'Done'::text, 'Cancelled'::text]))),
	CONSTRAINT pk_orders PRIMARY KEY (id)
);

ALTER TABLE "golang_db".orders ADD CONSTRAINT orders_payment_type_fkey FOREIGN KEY (payment_type) REFERENCES "golang_db".payment_type(id);
ALTER TABLE "golang_db".orders ADD CONSTRAINT orders_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES "golang_db".promos(id);
ALTER TABLE "golang_db".orders ADD CONSTRAINT orders_serve_id_fkey FOREIGN KEY (serve_id) REFERENCES "golang_db".serve(id);
ALTER TABLE "golang_db".orders ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES "golang_db".users(id);