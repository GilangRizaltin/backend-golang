package repositories

import (
	"Backend_Golang/internal/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	*sqlx.DB
}

func InitializeOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) RepositoryGetOrder(filter []string, page int) ([]models.OrderModel, error) {
	data := []models.OrderModel{}
	query := `
	select o.id as "No",
    u.full_name  as "User",
    o.subtotal as "Subtotal",
    p.promo_code as "Promo",
    o.percent_discount as "Discount_Percentage",
    o.flat_discount as "Discount_Flat",
    s.serve_type as "Serve",
    o.fee as "Serving_Fee",
    o.tax as "Tax",
    o.total_transactions as "Total_transaction",
    py.payment_name as "Payment_type",
    o.status as "Status",
    o.created_at as "Date"
    from orders o
    join users u on o.user_id = u.id
    join promos p on o.promo_id = p.id
    join serve s on o.serve_id = s.id 
    join payment_type py on o.payment_type = py.id
	`
	if filter[0] != "" {
		query += ` where o.status = '` + filter[0] + `'`
	}
	if filter[1] != "" {
		query += ` order by o.created_at`
		if filter[1] == "Newest" {
			query += ` desc`
		}
		if filter[1] == "Oldest" {
			query += ` asc`
		}
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*3)
	// fmt.Println(query)
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *OrderRepository) RepositoryGetOrderDetail(Id, page int) ([]models.OrderDetailModel, error) {
	data := []models.OrderDetailModel{}
	ID := strconv.Itoa(Id)
	query := `select
    o.id as "No Order",
    p.product_name as "Product_name",
    op.hot_or_not as "Hot_or_not",
    s.size_name as "Size",
    op.price as "Price",
    op.quantity as "Quantity"
    from
    orders_products op
    inner join
    orders o ON op.order_id = o.id
    inner join
    users u ON o.user_id = u.id
    join
    products p ON op.product_id = p.id
    join
    sizes s ON op.size_id = s.id
    where
    op.order_id = ` + ID
	// fmt.Println(query)
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func (r *OrderRepository) RepositoryCreateTransaction(bodyOrder *models.OrderModel, bodyOrderProducts *models.OrderDetailModel) error {
// 	tx, err := r.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	queryOrder := `INSERT INTO orders(user_id, subtotal, promo_id, percent_discount, flat_discount, serve_id, fee, tax, total_transactions, payment_type, status)
// 	VALUES (
// 		(SELECT id FROM users WHERE user_name = :User), :Subtotal,
// 		(select id from promos where id = 1),
// 		(select flat_amount from promos where id = 1),
// 		(select percent_amount from promos where id = 1),
// 		(SELECT id FROM serve WHERE serve_type = :Serve),
// 		(SELECT fee FROM serve WHERE serve_type = :Serve),
// 		0.1,
// 		:Total_transaction,
// 		(SELECT id FROM payment_type WHERE payment_name = :Payment_type),
// 		'On progress'
// 	) returning id`
// 	_, err = tx.NamedQuery(queryOrder, bodyOrder)
// 	if err != nil {
// 		return err
// 	}
// 	// queryOrderProduct := `INSERT INTO orders_products (order_id, product_id,hot_or_not, size_id, price, quantity, subtotal)
// 	// VALUES (
// 	// 	$1,
// 	// 	(SELECT id FROM products WHERE product_name = $2),
// 	// 	$3,
// 	// 	(SELECT id FROM sizes WHERE size_name = $4),
// 	// 	(
// 	// 		(SELECT price_default FROM products WHERE product_name = $2) +
// 	// 		(SELECT additional_fee FROM sizes WHERE size_name = $4)
// 	// 	),
// 	// 	$5,
// 	// 	(
// 	// 		(
// 	// 			(SELECT price_default FROM products WHERE product_name = $2) +
// 	// 			(SELECT additional_fee FROM sizes WHERE size_name = $4)
// 	// 		) * 3
// 	// 	)
// 	// )`
// 	queryOrderProduct := `INSERT INTO orders_products (order_id, product_id,hot_or_not, size_id, price, quantity, subtotal)
// 	VALUES (
// 		55,
// 		2,
// 		true,
// 		4,
// 		20000,
// 		5,
// 		100000
// 	)`
// 	// fmt.Println(queryOrderProduct)
// 	_, err = tx.Exec(queryOrderProduct) // bodyOrder.Id,
// 	// bodyOrderProducts.Product_name,
// 	// bodyOrderProducts.Hot_or_not,
// 	// bodyOrderProducts.Size,
// 	// bodyOrderProducts.Quantity
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *OrderRepository) RepositoryCreateOrder(bodyOrder *models.OrderModel, client *sqlx.Tx) (*sqlx.Rows, error) {
	queryOrder := `
        INSERT INTO orders(user_id, subtotal, promo_id, percent_discount, flat_discount, serve_id, fee, tax, total_transactions, payment_type, status)
        VALUES (
            (SELECT id FROM users WHERE user_name = :User), :Subtotal, 
            (SELECT id FROM promos WHERE id = 1), 
            (SELECT flat_amount FROM promos WHERE id = 1),
            (SELECT percent_amount FROM promos WHERE id = 1),
            (SELECT id FROM serve WHERE serve_type = :Serve),
            (SELECT fee FROM serve WHERE serve_type = :Serve),
            0.1,
            :Total_transaction,
            (SELECT id FROM payment_type WHERE payment_name = :Payment_type),
            'On progress'
        )
        RETURNING id
    `
	rows, err := client.NamedQuery(queryOrder, bodyOrder)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *OrderRepository) RepositoryCreateOrderProduct(bodyOrder *models.OrderModel, client *sqlx.Tx, orderId string) (*sqlx.Rows, error) {
	queryOrder := `
	INSERT INTO orders_products (order_id, product_id, hot_or_not, size_id, price, quantity, subtotal)
		VALUES `
	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["order_id"] = orderId

	j := 2
	for i := 0; i < len(bodyOrder.Product); i++ {
		filteredBody = append(filteredBody, "(:order_id")
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`(select id from products where product_name = :Product_name%d)`, j))
		filterBody[fmt.Sprintf("Product_name%d", j)] = bodyOrder.Product[i].Product_name
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`:Hot_or_not%d`, j))
		filterBody[fmt.Sprintf("Hot_or_not%d", j)] = bodyOrder.Product[i].Hot_or_not
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`(select id from sizes where size_name = :Size%d)`, j))
		filterBody[fmt.Sprintf("Size%d", j)] = bodyOrder.Product[i].Size
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`:Quantity%d`, j))
		filterBody[fmt.Sprintf("quantity%d", j)] = bodyOrder.Product[i].Quantity
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`:Price%d)`, j))
		filterBody[fmt.Sprintf("Price%d", j)] = bodyOrder.Product[i].Price
		//
		filteredBody = append(filteredBody, fmt.Sprintf(`:Subtotal_product%d)`, j))
		filterBody[fmt.Sprintf("Subtotal_product%d", j)] = bodyOrder.Product[i].Subtotal_product
		j++
	}
	if len(filteredBody) > 0 {
		queryOrder += strings.Join(filteredBody, ", ")
	}
	rows, err := client.NamedQuery(queryOrder, bodyOrder)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *OrderRepository) RepositoryUpdateOrder(ID int, body *models.OrderModel) (sql.Result, error) {
	query := `update orders
		set status = :Status`
	params := make(map[string]interface{})
	params["Status"] = body.Status
	params["Id"] = ID
	query += `, updated_at = NOW() WHERE id = :Id`
	fmt.Println(query)
	result, err := r.NamedExec(query, params)
	return result, err
}

// func (r *OrderRepository) RepositoryUpdateOrderDetail(ID int, body *models.OrderDetailModel) error {
// 	var mtx = sync.Mutex{}
// 	tx, err := r.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	queryOrderProducts := `
// 		UPDATE orders_products
// 		SET
// 			quantity = :Quantity,
// 			subtotal = (SELECT price FROM orders_products WHERE id = :ID) * :Quantity
// 		WHERE
// 			id = :ID
// 	`
// 	params := map[string]interface{}{
// 		"Quantity": body.Quantity,
// 		"ID":       ID,
// 	}
// 	mtx.Lock()
// 	_, err = tx.NamedExec(queryOrderProducts, params)
// 	if err != nil {
// 		return err
// 	}
// 	mtx.Unlock()
// 	queryOrder := `
// 		UPDATE orders
// 		SET
// 			subtotal = (
// 				SELECT SUM(subtotal) FROM orders_products WHERE order_id = (SELECT order_id FROM orders_products WHERE id = :ID)
// 			)
// 		WHERE
// 			id = (SELECT order_id FROM orders_products WHERE id = :ID)
// 	`
// 	_, err = tx.NamedExec(queryOrder, params)
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *OrderRepository) RepositoryDeleteProduct(ID int) (sql.Result, error) {
	query := `
        Update orders
		set deleted_at = now()
		where id = $1
		returning id;
    `
	result, err := r.Exec(query, ID)
	return result, err
}

func (r *OrderRepository) RepositoryCountOrder(filter []string) ([]int, error) {
	var total_data = []int{}
	query := `
		SELECT
			COUNT(*) AS "Total_order"
		FROM
			orders o `
	if filter[0] != "" {
		query += ` where o.status = '` + filter[0] + `'`
	}
	err := r.Select(&total_data, query)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return total_data, nil
}
