package repositories

import (
	"Backend_Golang/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	*sqlx.DB
}

func InitializeAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) RepositoryRegister(body *models.Auth, hashedPassword string) error {
	query := `
	insert into users(full_name, email, user_type, password_user) VALUES ($1, $2, 'Normal User', $3)
    `
	values := []any{body.Full_name, body.Email, hashedPassword}
	_, err := r.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositorySelectPrivateData(body *models.AuthLogin) ([]models.Auth, error) {
	data := []models.Auth{}
	query := `select u.id as "No",
	u.full_name as "Full_name",
	u.password_user as "Password",
	u.user_type as "User_type",
	u.otp as "Otp"
	from users u
	where u.email = $1`
	values := []any{body.Email}
	err := r.Select(&data, query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func (r *AuthRepository) RepositoryForgetPassword() {

}

func (r *AuthRepository) RepositoryResetPassword() {

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
