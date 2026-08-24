package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		&u.ID, &u.Name, &u.Email, &u.Pass,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (udb *UsersDB) GetUserInfo(u *models.User) string {
	return fmt.Sprintf("[ID:%d]=%s:%s_%s\n", u.ID, u.Email, u.Name, u.Pass)
}

func (udb *UsersDB) CreateUser(ur models.UserRegisterInput) (*models.User, error) {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES ($1, $2, $3) 
	RETURNING *;
	`
	// password hashing logic

	return udb.execWithError(query, ur.Name, ur.Email, ur.Pass)
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

func (udb *UsersDB) LoginUser(ul models.UserLoginInput, secret string) (string, error) {
	query := `
	SELECT * FROM users 
	WHERE email = $1;
	`
	row := udb.db.QueryRow(query, ul.Email)

	fetchedUser := models.UserLoginInput{}
	err := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.Email, &fetchedUser.Pass)
	if err != nil {
		return "", fmt.Errorf("Email is invalid")
	}

	// compare password and proc error

	return udb.GetJWT(fetchedUser.ID, secret)
}
