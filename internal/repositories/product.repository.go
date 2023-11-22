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

type ProductRepository struct {
	*sqlx.DB
}

func InitializeRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) RepositoryGet(conditions []string, page int) ([]models.ProductModel, error) {
	data := []models.ProductModel{}
	query := `
		SELECT
			p.id as "No",
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
	if conditions[0] != "" {
		conditional = append(conditional, "p.product_name ilike '%"+conditions[0]+"%'")
	}
	if conditions[1] != "" {
		conditional = append(conditional, "p.price_default < "+conditions[1])
	}
	if conditions[2] != "" {
		conditional = append(conditional, "p.price_default > "+conditions[2])
	}
	// if conditions[3] != "" {
	// 	conditional = append(conditional, "p.category = (SELECT id FROM categories c WHERE c.category_name = 'Coffee')")
	// }
	if conditions[3] != "" {
		conditional = append(conditional, "c.category_name = 'Coffee'")
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	if conditions[4] != "" {
		query += " ORDER BY "
		if conditions[4] == "Cheapest" {
			query += " p.price_default asc"
		}
		if conditions[4] == "Most Expensive" {
			query += " p.price_default desc"
		}
		if conditions[4] == "New Product" {
			query += " p.created_at desc"
		}
		if conditions[4] == "Oldest" {
			query += " p.created_at asc"
		}
		if conditions[4] == "" {
			query += " p.id asc"
		}
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*3)
	err := r.Select(&data, query)
	fmt.Println(query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProductRepository) RepositoryGetDetail(ID int) ([]models.ProductModel, error) {
	data := []models.ProductModel{}
	query := `SELECT
		p.id as "No",
		p.product_name as "Product",
		c.category_name as "Categories",
		p.description as "Description",
		p.price_default as "Price"
	FROM
		products p
	JOIN
		categories c ON p.category = c.id
	WHERE p.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProductRepository) RepositoryCreateProduct(body *models.ProductModel) error {
	query := `
		INSERT INTO products (product_name, category, description, price_default)
        VALUES (:Product, (SELECT id FROM categories WHERE category_name = :Categories), :Description, :Price);
    `
	_, err := r.NamedExec(query, body)
	return err
}

func (r *ProductRepository) RepositoryUpdateProduct(productID int, body *models.ProductModel) error {
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
	_, err := r.NamedExec(query, params)
	fmt.Println(query)
	return err
}

func (r *ProductRepository) RepositoryDeleteProduct(productID int) (sql.Result, error) {
	query := `
        DELETE FROM products
        WHERE
            id = $1
		returning product_name;
    `
	result, err := r.Exec(query, productID)
	return result, err
}

func (r *ProductRepository) RepositoryCountProduct(conditions []string) ([]int, error) {
	var total_data = []int{}
	query := `
		SELECT
			COUNT(*) AS "Total_product"
		FROM
			products p `
	var conditional []string
	if conditions[0] != "" {
		conditional = append(conditional, "p.product_name ilike '%"+conditions[0]+"%'")
	}
	if conditions[2] != "" {
		maxprice, _ := strconv.Atoi(conditions[1])
		conditional = append(conditional, "p.price_default < "+strconv.Itoa(maxprice))
	}
	if conditions[2] != "" {
		minprice, _ := strconv.Atoi(conditions[1])
		conditional = append(conditional, "p.price_default > "+strconv.Itoa(minprice))
	}
	if conditions[3] != "" {
		conditional = append(conditional, "p.category = "+conditions[3])
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	err := r.Select(&total_data, query)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return total_data, nil
}
