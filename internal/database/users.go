package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vonrimez/TaskAPI/internal/auth"
	"github.com/vonrimez/TaskAPI/internal/models"
)

type UsersDB struct {
	db *sql.DB
}

func GetNewUserDB(db *sql.DB) *UsersDB {
	return &UsersDB{db: db}
}

// need return values: id, name, email, pass
func (udb *UsersDB) execWithError(query string, args ...any) (*models.User, error) {
	row := udb.db.QueryRow(query, args...)
	u := models.User{}
	err := row.Scan(
		&u.ID, &u.Name,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (udb *UsersDB) CreateUser(ur models.UserRegisterInput) (*models.User, error) {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES ($1, $2, $3) 
	RETURNING id, name;
	`

	encryptedPass, err := auth.HashPassword(ur.Pass)
	if err != nil {
		return nil, err
	}

	return udb.execWithError(query, ur.Name, ur.Email, encryptedPass)
}

func (udb *UsersDB) LoginUser(ul models.UserLoginInput, secret string) (*models.User, error) {
	query := `
	SELECT id, name, email, password 
	FROM users 
	WHERE email = $1;
	`
	row := udb.db.QueryRow(query, ul.Email)

	fetchedUser := models.UserLoginInput{}
	err := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.Email, &fetchedUser.Pass)
	if err != nil {
		return nil, fmt.Errorf("Email or password is invalid")
	}

	if !auth.IsCorrectPassword(ul.Pass, fetchedUser.Pass) {
		return nil, fmt.Errorf("Email or password is invalid")
	}

	return &models.User{
		ID:   fetchedUser.ID,
		Name: fetchedUser.Name,
	}, nil
}

func (udb *UsersDB) GetJWT(userID int, secret string) (string, error) {
	type Claims struct {
		UserID int
		jwt.RegisteredClaims
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now)},
	},
	)

	return token.SignedString([]byte(secret))
}
