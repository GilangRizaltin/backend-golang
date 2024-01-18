package repositories

import (
	"Backend_Golang/internal/models"
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

type IUserRepository interface {
	RepositoryGetUser(body *models.QueryParamsUser) ([]models.UserModel, error)
	RepositoryGetUserProfile(ID int) ([]models.UserModel, error)
	RepositorySensitiveDataUser(ID int) ([]models.UserModel, error)
	RepositoryAddUser(body *models.UserModel, hashedPassword, url string) error
	RepositoryUpdateUser(productID int, body *models.UserUpdateModel, url, hashedPassword string) (int64, error)
	RepositoryUpdatePasswordUser(userID int, hashedPassword string) error
	RepositoryDeleteUser(userID int) (int64, error)
	RepositoryCountUser(body *models.QueryParamsUser) ([]int, error)
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
	WHERE u.deleted_at is null 
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
	if len(conditional) == 1 {
		query += fmt.Sprintf(" AND %s", conditional[0])
	}
	if len(conditional) > 1 {
		query += fmt.Sprintf(" AND %s", strings.Join(conditional, " AND "))
	}
	if body.SortOrder != "" {
		if body.SortOrder == "asc" {
			query += " ORDER BY u.full_name asc"
		}
		if body.SortOrder == "desc" {
			query += " ORDER BY u.full_name desc"
		}
	}
	if body.SortOrder == "" {
		query += " ORDER BY u.id asc"
	}
	var page = body.Page
	if body.Page == 0 {
		page = 1
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*6)
	// fmt.Println(body.SortOrder)
	// fmt.Println(query)
	err := r.Select(&data, query, values...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositoryGetUserProfile(ID int) ([]models.UserModel, error) {
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
	u.created_at as "created_at",
	u.otp as "Otp"
	from users u
	where u.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositorySensitiveDataUser(ID int) ([]models.UserModel, error) {
	data := []models.UserModel{}
	query := `
	select 
	u.password_user as "Password",
	u.otp as "Otp"
	from users u
	where u.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func (r *UserRepository) RepositoryRegisterUser(body *models.UserModel) (*sqlx.Rows, error) {
// 	query := `
// 	insert into users(full_name, email, user_type, password_user) VALUES (:Full_name, :Email, 'Normal User', :Password) returning id, full_name
//     `
// 	result, err := r.NamedQuery(query, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

func (r *UserRepository) RepositoryAddUser(body *models.UserModel, hashedPassword, url string) error {
	query := `
	insert into users(user_photo_profile, user_name, full_name, address, phone, email, user_type, password_user) VALUES 
	(:Photo_profile, :User_name, :Full_name, :Address,  :Phone, :Email, :User_type, :hashedPassword) returning id, full_name
    `
	params := make(map[string]interface{})
	params["Photo_profile"] = url
	params["User_name"] = body.User_name
	params["Full_name"] = body.Full_name
	params["Address"] = body.Address
	params["Phone"] = body.Phone
	params["Email"] = body.Email
	params["User_type"] = body.User_type
	params["hashedPassword"] = hashedPassword
	_, err := r.NamedExec(query, params)
	return err
}

func (r *UserRepository) RepositoryUpdateUser(productID int, body *models.UserUpdateModel, url, hashedPassword string) (int64, error) {
	var conditional []string
	query := `
        UPDATE users
        SET `
	params := make(map[string]interface{})
	if url != "" {
		conditional = append(conditional, "user_photo_profile = :Url")
		params["Url"] = url
	}
	if body.User_name != "" {
		conditional = append(conditional, "user_name = :User_name")
		params["User_name"] = body.User_name
	}
	if body.Full_name != "" {
		conditional = append(conditional, "full_name = :Full_name")
		params["Full_name"] = body.Full_name
	}
	if body.Phone != nil {
		conditional = append(conditional, "phone = :Phone")
		params["Phone"] = body.Phone
	}
	if hashedPassword != "" {
		conditional = append(conditional, "password_user = :Password")
		params["Password"] = hashedPassword
	}
	if body.Address != nil {
		conditional = append(conditional, "address = :Address")
		params["Address"] = body.Address
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
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func (r *UserRepository) RepositoryUpdatePasswordUser(userID int, hashedPassword string) error {
	query := `UPDATE users SET password_user = $1 where id = $2`
	values := []any{hashedPassword, userID}
	_, err := r.Query(query, values)
	if err != nil {
		return nil
	}
	return nil
}

func (r *UserRepository) RepositoryDeleteUser(userID int) (int64, error) {
	query := `
        update users
        set 
		deleted_at = now ()
        where id = $1
		returning full_name;
    `
	result, err := r.Exec(query, userID)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func (r *UserRepository) RepositoryCountUser(body *models.QueryParamsUser) ([]int, error) {
	var totalData = []int{}
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
	err := r.Select(&totalData, query, values...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return totalData, nil
}
