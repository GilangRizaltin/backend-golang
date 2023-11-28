package repositories

import (
	"Backend_Golang/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	*sqlx.DB
}

func InitializeRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) RepositoryGet(body *models.QueryParamsProduct) ([]models.ProductModel, error) {
	data := []models.ProductModel{}
	query := `
		SELECT
			p.id as "No",
			p.product_image_1 as "Product_photo_1",
			p.product_image_2 as "Product_photo_2",
			p.product_image_3 as "Product_photo_3",
			p.product_image_4 as "Product_photo_4",
			p.product_name as "Product",
			c.category_name as "Categories",
			p.description as "Description",
			p.price_default as "Price"
		FROM
			products p
		JOIN
			categories c ON p.category = c.id
	`
	var conditional []string
	values := []any{}
	if body.ProductId != 0 {
		conditional = append(conditional, "p.id = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.ProductId)
	}
	if body.ProductName != "" {
		conditional = append(conditional, "p.product_name ilike $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.ProductName+"%")
	}
	if body.MaximumPrice != 0 {
		conditional = append(conditional, "p.price_default <  $"+fmt.Sprint(len(values)+1))
		values = append(values, body.MaximumPrice)
	}
	if body.MinimumPrice != 0 {
		conditional = append(conditional, "p.price_default >  $"+fmt.Sprint(len(values)+1))
		values = append(values, body.MinimumPrice)
	}
	if body.ProductCategory != "" {
		conditional = append(conditional, "c.category_name = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.ProductCategory)
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	// query += " WHERE p.deleted_at is null "
	if body.Sort == "" {
		query += " ORDER BY p.id asc"
	}
	if body.Sort != "" {
		query += " ORDER BY "
		if body.Sort == "Cheapest" {
			query += " p.price_default asc"
		}
		if body.Sort == "Most Expensive" {
			query += " p.price_default desc"
		}
		if body.Sort == "New Product" {
			query += " p.created_at desc"
		}
		if body.Sort == "Oldest" {
			query += " p.created_at asc"
		}
	}
	var page = body.Page
	if body.Page == 0 {
		page = 1
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*6)
	err := r.Select(&data, query, values...)
	// fmt.Println(query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProductRepository) RepositoryGetDetail(ID int) ([]models.ProductModel, error) {
	data := []models.ProductModel{}
	query := `SELECT
		p.id as "No",
		p.product_image_1 as "Product_photo_1",
		p.product_image_2 as "Product_photo_2",
		p.product_image_3 as "Product_photo_3",
		p.product_image_4 as "Product_photo_4",
		p.product_name as "Product",
		c.category_name as "Categories",
		p.description as "Description",
		p.price_default as "Price"
	FROM
		products p
	JOIN
		categories c ON p.category = c.id
	WHERE p.id = $1`
	// and p.deleted_at is null`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProductRepository) RepositoryCreateProduct(body *models.ProductModel, dataUrl []string) error {
	query := `
		INSERT INTO products (product_name, category, description, price_default, product_image_1, product_image_2, product_image_3, product_image_4)
        VALUES (:Product, (SELECT id FROM categories WHERE category_name = :Categories), :Description, :Price, :Product_Image_1, :Product_Image_2, :Product_Image_3, :Product_Image_4);`
	params := make(map[string]interface{})
	params["Product"] = body.Product_name
	params["Categories"] = body.Category
	params["Description"] = body.Description
	params["Price"] = body.Price_default
	params["Product_Image_1"] = dataUrl[0]
	params["Product_Image_2"] = dataUrl[1]
	params["Product_Image_3"] = dataUrl[2]
	params["Product_Image_4"] = dataUrl[3]
	_, err := r.NamedExec(query, params)
	return err
}

func (r *ProductRepository) RepositoryUpdateProduct(productID int, body *models.UpdateProduct) (sql.Result, error) {
	var conditional []string
	query := `
        UPDATE products
        SET `
	params := make(map[string]interface{})
	if body.Product_name != "" {
		conditional = append(conditional, "product_name = :Product_name")
		params["Product_name"] = body.Product_name
	}
	if body.Category != "" {
		conditional = append(conditional, "category = (SELECT id FROM categories WHERE category_name = :Category)")
		params["Category"] = body.Category
	}
	if body.Description != "" {
		conditional = append(conditional, "description = :Description")
		params["Description"] = body.Description
	}
	if len(conditional) == 1 {
		query += conditional[0] + ", "
	}
	if len(conditional) > 1 {
		query += strings.Join(conditional, ", ") + ", "
	}
	params["Id"] = productID
	query += ` update_at = NOW() WHERE id = :Id`
	result, err := r.NamedExec(query, params)
	// fmt.Println(query)
	return result, err
}

func (r *ProductRepository) RepositoryDeleteProduct(productID int) (sql.Result, error) {
	query := `
        update products set deleted_at = now() where id = $1;
    `
	result, err := r.Exec(query, productID)
	return result, err
}

func (r *ProductRepository) RepositoryCountProduct(body *models.QueryParamsProduct) ([]int, error) {
	var total_data = []int{}
	query := `
		SELECT
			COUNT(*) AS "Total_product"
		FROM
			products p 
		JOIN
			categories c ON p.category = c.id`
	var conditional []string
	values := []any{}
	if body.ProductId != 0 {
		conditional = append(conditional, "p.id = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.ProductId)
	}
	if body.ProductName != "" {
		conditional = append(conditional, "p.product_name ilike $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.ProductName+"%")
	}
	if body.MaximumPrice != 0 {
		conditional = append(conditional, "p.price_default <  $"+fmt.Sprint(len(values)+1))
		values = append(values, body.MaximumPrice)
	}
	if body.MinimumPrice != 0 {
		conditional = append(conditional, "p.price_default >  $"+fmt.Sprint(len(values)+1))
		values = append(values, body.MinimumPrice)
	}
	if body.ProductCategory != "" {
		conditional = append(conditional, "c.category_name = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.ProductCategory)
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	err := r.Select(&total_data, query, values...)
	if err != nil {
		// log.Fatalln(err)
		fmt.Println(query)
		return nil, err
	}
	return total_data, nil
}

func (r *ProductRepository) RepositoryStatisticProduct(dateStart, dateEnd string) ([]models.StatisticProduct, error) {
	data := []models.StatisticProduct{}
	query := `SELECT 
                dates::date AS "OrderDate",
                SUM(op.quantity) AS "TotalQuantity"
              FROM 
                generate_series($1::timestamp, $2::timestamp, interval '1 day') dates
              LEFT JOIN 
                orders_products AS op
              ON 
                DATE(op.created_at) = dates::date
              GROUP BY 
                dates::date
              ORDER BY 
                dates::date`
	values := []any{
		dateStart, dateEnd,
	}
	err := r.Select(&data, query, values...)
	// fmt.Println(query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProductRepository) RepositoryPopularProduct(dateStart, dateEnd string) ([]models.PopularProduct, error) {
	data := []models.PopularProduct{}
	query := `SELECT
    p.product_name as "Product",
    SUM(op.quantity) as "Total_Quantity",
    SUM(op.subtotal) as "Total_Income"
FROM
    orders_products AS op
JOIN
    products AS p
ON
    op.product_id = p.id
WHERE
	op.created_at > $1::timestamp
AND 
	op.created_at < $2::timestamp
GROUP BY
    p.product_name
HAVING
    SUM(op.quantity) IS NOT NULL
ORDER BY
    "Product" ASC`
	values := []any{
		dateStart, dateEnd,
	}
	err := r.Select(&data, query, values...)
	if err != nil {
		return nil, err
	}
	return data, nil
}
