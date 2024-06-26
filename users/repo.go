package users
import (
	"anonymous/commons"
	"anonymous/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
}

func Repo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) MustInsert(tx *sqlx.Tx, user *models.User) error {
	_, err := tx.NamedExec(
		`
		INSERT INTO users (
    id, email, username, password_hash, active, profile_picture, joined_at , email_verified , email_verification_token
    )
    VALUES (
    :id, :email, :username, :password_hash, :active,:profile_picture, :joined_at, :email_verified , :email_verification_token
    );
    `,
		user,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}
func (r *UserRepo) GetUser(field, value string) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf("select * from users where %s=$1", field)
	err := r.db.Get(
		user,
		query,
		value,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("Error while retrieving user by %s: %w", field, err)
	}
	return user, nil
}
func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := "SELECT * FROM users WHERE email = $1"
    err := r.db.Get(user, query, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, commons.Errors.ResourceNotFound
        }
        return nil, fmt.Errorf("error while retrieving user by email: %w", err)
    }
    return user, nil
}


func (r *UserRepo) CheckDuplicates(email string) (string, error) {
	result := ""
	err := r.db.QueryRow(
		`
      SELECT
      CASE
        WHEN EXISTS (
            SELECT 1
            FROM users
            WHERE email = $1
        )
        THEN 'email'
        ELSE 'none'
      END AS taken_by;
    `,
		email,
	).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("Error while checking for duplicates: %w", err)
	}
	return result, nil
}
func (r *UserRepo) CheckDuplicatesU(username string) (string, error) {
	result := ""
	err := r.db.QueryRow(
		`
      SELECT
      CASE
        WHEN EXISTS (
            SELECT 1
            FROM users
            WHERE username = $1
        )
        THEN 'username'
        ELSE 'none'
      END AS taken_by;
    `,
		username,
	).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("Error while checking for duplicates: %w", err)
	}
	return result, nil
}



func (r *UserRepo) GetUserDataByID(id string) (*models.LoggedInUser, error) {
	user := models.LoggedInUser{}
	err := r.db.QueryRowx(
		`
      SELECT
        id, username, password_hash, email, email_verified, joined_at, active, profile_picture , email_verification_token
      FROM
        users
      WHERE 
        id = $1
    `,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.EmailVerified,
		&user.JoinedAt,
		&user.Active,
		&user.ProfilePicture,
		&user.EmailVerificationToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("Error while getting logged in user data: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) ChangePassword(password, id string) error {
	_, err := r.db.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", password, id)
	if err != nil {
		return fmt.Errorf("Error while changing user password: %w", err)
	}
	return nil
}

func (r *UserRepo) ToggleStatus(users []string, status bool) error {
	_, err := r.db.Exec("UPDATE users SET active = $1 WHERE id = ANY($2)", status, pq.Array(users))
	if err != nil {
		return fmt.Errorf("Error while changing accounts status: %w", err)
	}
	return nil
}

func (r *UserRepo) GetAllUsersData() (*[]models.LoggedInUser, error) {
	data := []models.LoggedInUser{}
	rows, err := r.db.Queryx(
		`
      SELECT
        id, username, password_hash, email, email_verified, joined_at, active, profile_picture , email_verification_token
      FROM
        users
    `,
	)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving users data: %w", err)
	}
	for rows.Next() {
		user := models.LoggedInUser{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.EmailVerified,
			&user.JoinedAt,
			&user.Active,
			&user.ProfilePicture,
			&user.EmailVerificationToken,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while retrieving users data: error while scanning row: %w", err)
		}
		data = append(data, user)
	}
	return &data, nil
}

func (r *UserRepo) SetContactVerified(userId string) error {
	query := "UPDATE users SET email_verified = true WHERE id = $1"
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("Error while setting user contact to verified: %w", err)
	}
	return nil
}

func (r *UserRepo) SetEmailVerificationToken(userID, token string) error {
    query := `UPDATE users SET email_verification_token = $1 WHERE id = $2`
    _, err := r.db.Exec(query, token, userID)
    return err
}

func (r *UserRepo) GetUserByVerificationToken(token string) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE email_verification_token = $1`
    err := r.db.Get(&user, query, token)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
func (r *UserRepo) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email_verified = $1, email_verification_token = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, user.EmailVerified, user.EmailVerificationToken, user.ID)
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}
	return nil
}

// FindByVerificationToken finds a user by their email verification token.
func (r *UserRepo) FindByVerificationToken(token string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, joined_at, active, profile_picture, email_verified, email_verification_token
		FROM users
		WHERE email_verification_token = $1
	`
	var user models.User
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password, &user.JoinedAt, &user.Active,
		&user.ProfilePicture, &user.EmailVerified, &user.EmailVerificationToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) VerifyEmail(token string) error {
	var user models.User

	// Find the user by the verification token
	query := `SELECT * FROM users WHERE email_verification_token = $1`
	err := r.db.Get(&user, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commons.Errors.ResourceNotFound
		}
		return fmt.Errorf("could not find user by verification token: %w", err)
	}
	user.EmailVerified = true
	user.EmailVerificationToken = ""
	query = `
        UPDATE users
        SET email_verified = $1, email_verification_token = $2
        WHERE id = $3
    `
	_, err = r.db.Exec(query, user.EmailVerified, user.EmailVerificationToken, user.ID)
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	return nil
}