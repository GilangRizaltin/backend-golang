package repositories

import (
	"Backend_Golang/internal/models"
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

func (r *PromoRepository) RepositoryGetPromo(conditions []string, page int) ([]models.PromoModel, error) {
	data := []models.PromoModel{}
	query := `select p.id as "No",
	p.promo_code as "Promo_code",
	pt.promo_type_name as "Promo_type",
	p.flat_amount as "Flat_amount",
	p.percent_amount as "Percent_amount",
	p.created_at as "Time_created",
	p.ended_at as "Time_ended"
	from 
	promos p 
	join
	promos_type pt on p.promo_type = pt.id
	`
	var conditional []string
	if conditions[0] != "" {
		conditional = append(conditional, "p.promo_code ilike '%"+conditions[0]+"%'")
	}
	if conditions[1] != "" {
		conditional = append(conditional, "p.ended_at = "+conditions[1])
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*3)
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *PromoRepository) RepositoryCreatePromo(body *models.PromoModel) error {
	query := `INSERT INTO promos (promo_code, promo_type, flat_amount, percent_amount, ended_at) 
	VALUES (:Promo_code, (select from promos_type where promo_type_name = :Promo_type), :Flat_amount, :Percent_amount, :Ended_at) RETURNING promo_code, ended_at`
	_, err := r.NamedExec(query, body)
	return err
}

// func (r *PromoRepository) RepositoryUpdatePromo(productID int, body *models.PromoModel) error {

// }

// func (r *PromoRepository) RepositoryDeletePromo(productID int) (sql.Result, error) {

// }
