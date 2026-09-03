package database

import (
	"database/sql"
	"errors"

	"github.com/vonrimez/TaskAPI/domain"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type UsersDB struct {
	db *sql.DB
}

func GetNewUserDB(db *sql.DB) *UsersDB {
	return &UsersDB{db: db}
}

// required return values: id, name, pass
func (udb *UsersDB) execWithError(query string, args ...any) (*models.UserOutput, *domain.AppError) {
	row := udb.db.QueryRow(query, args...)
	outputUser := models.UserOutput{}
	err := row.Scan(
		&outputUser.ID, &outputUser.Name, &outputUser.Pass,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("there are no matches")
		}
		return nil, domain.NewInternalError(err)
	}
	return &outputUser, nil
}

func (udb *UsersDB) CreateUser(inputUser models.UserRegisterInput) (*models.UserOutput, *domain.AppError) {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES ($1, $2, $3) 
	RETURNING id, name, pass;
	`

	return udb.execWithError(query, inputUser.Name, inputUser.Email, inputUser.Pass)
}

func (udb *UsersDB) LoginUser(inputUser models.UserLoginInput) (*models.UserOutput, *domain.AppError) {
	query := `
	SELECT id, name, email, password 
	FROM users 
	WHERE email = $1;
	`

	return udb.execWithError(query, inputUser.Email)
}
