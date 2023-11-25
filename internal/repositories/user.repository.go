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

type UserRepository struct {
	*sqlx.DB
}

func InitializeUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) RepositoryGetUser(body *models.QueryParamsUser) ([]models.UserModel, error) {
	data := []models.UserModel{}
	query := `
	select u.id as "No",
	u.user_photo_profile as "Photo_profile",
	u.user_name as "User_name",
	u.full_name as "Full_name",
	u.phone as "Phone",
	u.address as "Address",
	u.email as "Email",
	u.user_type as "User_type",
	u.otp as "Otp"
	from users u
	`
	var conditional []string
	values := []any{}
	if body.Userid != "" {
		conditional = append(conditional, "u.id = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.Userid)
	}
	if body.Username != "" {
		conditional = append(conditional, "u.user_name ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Username+"%")
	}
	if body.Fullname != "" {
		conditional = append(conditional, "u.full_name ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Fullname+"%")
	}
	if body.Email != "" {
		conditional = append(conditional, "u.email ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Email+"%")
	}
	if body.Phone != "" {
		conditional = append(conditional, "u.phone ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Phone+"%")
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	if body.SortOrder != "" {
		query += " ORDER BY u.id " + body.SortOrder
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((body.Page-1)*3)
	err := r.Select(&data, query, values...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositoryGetUserProfile(ID string) ([]models.UserModel, error) {
	data := []models.UserModel{}
	query := `
	select u.id as "No",
	u.user_photo_profile as "Photo_profile",
	u.user_name as "User_name",
	u.full_name as "Full_name",
	u.phone as "Phone",
	u.address as "Address",
	u.email as "Email",
	u.user_type as "User_type",
	u.otp as "Otp"
	from users u
	where u.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositoryRegisterUser(body *models.UserModel) (*sqlx.Rows, error) {
	query := `
	insert into users(full_name, email, user_type, password_user) VALUES (:Full_name, :Email, 'Normal User', :Password) returning id, full_name
    `
	result, err := r.NamedQuery(query, body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) RepositoryAddUser(body *models.UserModel) error {
	query := `
	insert into users(user_name, full_name, address, phone, email, user_type, password_user) VALUES (:User_name, :Full_name, :Address, :Phone,  :Email, :User_type, :Password) returning id, full_name
    `
	_, err := r.NamedExec(query, body)
	return err
}

func (r *UserRepository) RepositoryUpdateUser(productID int, body *models.UserModel) (sql.Result, error) {
	var conditional []string
	query := `
        UPDATE users
        SET `
	params := make(map[string]interface{})
	if body.User_name != nil {
		conditional = append(conditional, "user_name = :User_name")
		params["User_name"] = body.User_name
	}
	if body.Full_name != nil {
		conditional = append(conditional, "full_name = :Full_name")
		params["Full_name"] = body.Full_name
	}
	if body.Phone != nil {
		conditional = append(conditional, "phone = :Phone")
		params["Phone"] = body.Phone
	}
	if body.Address != nil {
		conditional = append(conditional, "address = :Address")
		params["Address"] = body.Address
	}
	if body.Password != "" {
		conditional = append(conditional, "password_user = :Password")
		params["Password"] = body.Password
	}
	if body.User_type != "" {
		conditional = append(conditional, "user_type = :User_type")
		params["User_type"] = body.User_type
	}
	if len(conditional) == 1 {
		query += conditional[0]
	}
	if len(conditional) > 1 {
		query += strings.Join(conditional, ", ")
	}
	params["Id"] = productID
	query += ` ,update_at = NOW() WHERE id = :Id`
	fmt.Println(query)
	result, err := r.NamedExec(query, params)
	return result, err
}

func (r *UserRepository) RepositoryDeleteUser(userID int) (sql.Result, error) {
	query := `
        update users
        set 
		deleted_at = now ()
        where id = $1
		returning full_name;
    `
	result, err := r.Exec(query, userID)
	return result, err
}

func (r *UserRepository) RepositoryCountUser(body *models.QueryParamsUser) (int, error) {
	var totalData int
	query := `
		SELECT
			COUNT(*) AS "Total_user"
		FROM
			users u `
	var conditional []string
	values := []any{}
	if body.Userid != "" {
		conditional = append(conditional, "u.id = $"+fmt.Sprint(len(values)+1))
		values = append(values, body.Userid)
	}
	if body.Username != "" {
		conditional = append(conditional, "u.user_name ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Username+"%")
	}
	if body.Fullname != "" {
		conditional = append(conditional, "u.full_name ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Fullname+"%")
	}
	if body.Email != "" {
		conditional = append(conditional, "u.email ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Email+"%")
	}
	if body.Phone != "" {
		conditional = append(conditional, "u.phone ILIKE $"+fmt.Sprint(len(values)+1))
		values = append(values, "%"+body.Phone+"%")
	}
	if len(conditional) > 0 {
		query += " WHERE " + strings.Join(conditional, " AND ")
	}
	// fmt.Println(query)
	err := r.Get(&totalData, query, values...)
	if err != nil {
		log.Fatalln(err)
		return 0, err
	}
	return totalData, nil
}
