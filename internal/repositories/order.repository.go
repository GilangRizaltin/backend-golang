package repositories

import (
	"Backend_Golang/internal/models"
	"database/sql"
	"strconv"

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
    p.promo_code as "Promo_Code",
    o.percent_discount as "Discount_Percentage",
    o.flat_discount as "Discount_Flat",
    s.serve_type as "Serving_Type",
    o.fee as "Serving_Fee",
    o.tax as "Tax",
    o.total_transactions as "Total_Transactions",
    py.payment_name as "Payment_Type",
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
    op.id as "No Order",
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
    op.id = ` + ID
	// fmt.Println(query)
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *OrderRepository) RepositoryCreateTransaction(bodyOrder *models.OrderModel, bodyOrderProducts *models.OrderDetailModel) {

}

func (r *OrderRepository) RepositoryUpdateOrder(ID int, body *models.OrderModel) (sql.Result, error) {
	query := `update orders
		set status = :Status 
		where id = :ID`
	params := make(map[string]interface{})
	params["Status"] = body.Status
	params["ID"] = ID
	query += ` update_at = NOW() WHERE id = :Id`
	result, err := r.NamedExec(query, params)
	return result, err
}

// func (r *OrderRepository) RepositoryUpdateOrderDetail(ID int, body *models.OrderDetailModel) error {
// 	tx, err := r.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	queryOrderProducts := `update orders_products
// 		set quantity = :Quantity,
// 		subtotal = (select price from orders_products where id = :ID) * :Quantity
// 		where id = :ID`
// 	params := make(map[string]interface{})
// 	params["Quantity"] = body.Quantity
// 	params["ID"] = ID
// 	query += ` updated_at = NOW() WHERE id = :Id`
// 	_, err := tx.NamedExec(queryOrderProducts, params)
// 	if err != nil {
// 		return err
// 	}
// 	queryOrder := `update orders
// 		set total_transaction = (select sum (subtotal) from orders_products where order_id = :ID)
// 		where id = (select order_id from order_products where id = :ID)`
// 	_, err := tx.NamedExec(queryOrderProducts, params)
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
