package repositories

import (
	"Backend_Golang/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	*sqlx.DB
}

type IAuthRepository interface {
	RepositoryRegister(body *models.AuthRegister, hashedPassword string, OTP int) error
	RepositorySelectPrivateData(email string) ([]models.Auth, error)
	RepositoryActivateUser(email string) (int64, error)
	RepositoryLogout(token string) error
}

func InitializeAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) RepositoryRegister(body *models.AuthRegister, hashedPassword string, OTP int) error {
	query := `
	insert into users(full_name, email, user_type, password_user, otp) VALUES ($1, $2, 'Normal User', $3, $4)
    `
	values := []any{body.Full_name, body.Email, hashedPassword, OTP}
	_, err := r.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositorySelectPrivateData(email string) ([]models.Auth, error) {
	data := []models.Auth{}
	query := `select u.id as "No",
	u.full_name as "Full_name",
	u.user_photo_profile as "Photo_profile",
	u.password_user as "Password",
	u.user_type as "User_type",
	u.otp as "Otp",
	u.activated as "activated"
	from users u
	where u.email = $1`
	values := []any{email}
	err := r.Select(&data, query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func (r *AuthRepository) RepositoryActivateUser(email string) (int64, error) {
	const sql = `update users set activated = true where email = :Email`
	params := make(map[string]interface{})
	params["Email"] = email
	result, err := r.NamedExec(sql, params)
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, err
}

func (r *AuthRepository) RepositoryLogout(token string) error {
	query := `insert into jwt (jwt_code) values ($1)`
	values := []any{token}
	_, err := r.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositoryIsTokenBlacklisted(token string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM jwt WHERE jwt_code = $1`
	err := r.Get(&count, query, token)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
