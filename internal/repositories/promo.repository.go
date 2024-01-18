package repositories

import (
	"Backend_Golang/internal/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PromoRepository struct {
	*sqlx.DB
}

func InitializePromoRepository(db *sqlx.DB) *PromoRepository {
	return &PromoRepository{db}
}

type IPromoRepository interface {
	RepositoryGetPromo(body *models.QueryParamsPromo) ([]models.PromoModel, error)
	RepositoryCreatePromo(body *models.PromoModel) error
	RepositoryUpdatePromo(productID int, body *models.UpdatePromoModel) (int64, error)
	RepositoryDeletePromo(productID int) (int64, error)
	RepositoryCountPromo(body *models.QueryParamsPromo) ([]int, error)
}

func (r *PromoRepository) RepositoryGetPromo(body *models.QueryParamsPromo) ([]models.PromoModel, error) {
	data := []models.PromoModel{}
	query := `select p.id as "No",
	p.promo_code as "Promo_code",
	pt.promo_type_name as "Promo_type",
	p.flat_amount as "Flat_amount",
	p.percent_amount as "Percent_amount",
	p.created_at as "Time_created",
	p.ended_at as "Ended_at"
	from 
	promos p 
	join
	promos_type pt on p.promo_type = pt.id
	`
	var conditional []string
	values := []any{}
	if body.Promo_code != "" {
		conditional = append(conditional, "p.promo_code ilike $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Promo_code+"%")
	}
	if body.Time_end != nil && body.Time_end.String() != "" {
		conditional = append(conditional, "p.ended_at < $"+fmt.Sprint(len(values)+1))
		values = append(values, body.Time_end.String())
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((body.Page-1)*6)
	err := r.Select(&data, query, values...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *PromoRepository) RepositoryCreatePromo(body *models.PromoModel) error {
	query := `INSERT INTO promos (promo_code, promo_type, flat_amount, percent_amount, ended_at) 
	VALUES (:Promo_code, (select id from promos_type where promo_type_name = :Promo_type), :Flat_amount, :Percent_amount, NOW() + interval '7 days') RETURNING promo_code`
	_, err := r.NamedExec(query, body)
	return err
}

func (r *PromoRepository) RepositoryUpdatePromo(productID int, body *models.UpdatePromoModel) (int64, error) {
	var conditional []string
	query := `
        UPDATE promos
        SET `
	params := make(map[string]interface{})
	if body.Promo_code != "" {
		conditional = append(conditional, "promo_code = :Promo_code")
		params["Promo_code"] = body.Promo_code
	}
	if body.Ended_at != "" {
		conditional = append(conditional, "ended_at = :Ended_at")
		params["Ended_at"] = body.Ended_at
	}
	if len(conditional) == 1 {
		query += conditional[0] + ", "
	}
	if len(conditional) > 1 {
		query += strings.Join(conditional, ", ")
	}
	params["Id"] = productID
	query += ` update_at = NOW() WHERE id = :Id`
	result, err := r.NamedExec(query, params)
	// fmt.Println(query)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (r *PromoRepository) RepositoryDeletePromo(productID int) (int64, error) {
	query := `
	update products set deleted_at = now() where id = $1;
	`
	result, err := r.Exec(query, productID)
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, err
}

func (r *PromoRepository) RepositoryCountPromo(body *models.QueryParamsPromo) ([]int, error) {
	var total_data = []int{}
	query := `
		SELECT
			COUNT(*) AS "Total_promo"
		FROM
			promos p `
	var conditional []string
	values := []any{}
	if body.Promo_code != "" {
		conditional = append(conditional, "p.promo_code ilike $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Promo_code+"%")
	}
	if body.Time_end != nil && body.Time_end.String() != "" {
		conditional = append(conditional, "p.ended_at < $"+fmt.Sprint(len(values)+1))
		values = append(values, body.Time_end.String())
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	err := r.Select(&total_data, query, values...)
	if err != nil {
		// log.Fatalln(err)
		return nil, err
	}
	return total_data, nil
}
