package users

import (
	"anonymous/commons"
	"anonymous/encryption"
	"anonymous/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
	ee *encryption.EmailEncryption
}

func Repo(db *sqlx.DB) *UserRepo {
	ee, _ := encryption.NewEmailEncryption()
	return &UserRepo{
		db: db,
		ee: ee,
	}
}

func (u *UserRepo) MustInsert(tx *sqlx.Tx, user *models.User) error {
	// Chiffrement de l'email avant insertion
	if err := user.EncryptEmail(); err != nil {
		return fmt.Errorf("erreur de chiffrement de l'email: %w", err)
	}
	
	_, err := tx.NamedExec(`
		INSERT INTO users (
		id, email, username, password_hash, active, profile_picture, joined_at, email_verified, email_verification_token
		)
		VALUES (
		:id, :encrypted_email, :username, :password_hash, :active, :profile_picture, :joined_at, :email_verified, :email_verification_token
		);
	`, map[string]interface{}{
		"id":                      user.ID,
		"encrypted_email":         user.EncryptedEmail,
		"username":                user.Username,
		"password_hash":           user.Password,
		"active":                  user.Active,
		"profile_picture":         user.ProfilePicture,
		"joined_at":               user.JoinedAt,
		"email_verified":          user.EmailVerified,
		"email_verification_token": user.EmailVerificationToken,
	})
	
	if err != nil {
		return fmt.Errorf("error while inserting user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUser(field, value string) (*models.User, error) {
	user := &models.User{}
	var query string
	var err error
	
	// Si la recherche est par email, il faut d'abord chiffrer la valeur
	if field == "email" {
		encryptedValue, err := r.ee.EncryptEmail(value)
		if err != nil {
			return nil, fmt.Errorf("erreur de chiffrement de l'email: %w", err)
		}
		value = encryptedValue
	}
	
	query = fmt.Sprintf("select * from users where %s=$1", field)
	err = r.db.Get(user, query, value)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("error while retrieving user by %s: %w", field, err)
	}
	
	// Déchiffrement de l'email
	if user.EncryptedEmail != "" {
		if err := user.DecryptEmail(); err != nil {
			return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
		}
	}
	
	return user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	encryptedEmail, err := r.ee.EncryptEmail(email)
    if err != nil {
        return nil, fmt.Errorf("erreur de chiffrement de l'email: %w", err)
    }
    
    // Log temporaire pour debug
	user := &models.User{}
	query := `SELECT * FROM users WHERE email = $1`
	err = r.db.Get(user, query, encryptedEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur par email: %w", err)
	}

	return user, nil
}

func (r *UserRepo) CheckDuplicates(email string) (string, error) {
	// Chiffrer l'email pour vérifier les doublons
	encryptedEmail, err := r.ee.EncryptEmail(email)
	if err != nil {
		return "", fmt.Errorf("erreur de chiffrement de l'email: %w", err)
	}
	
	result := ""
	err = r.db.QueryRow(`
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
		encryptedEmail,
	).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("error while checking for duplicates: %w", err)
	}
	return result, nil
}

func (r *UserRepo) CheckDuplicatesU(username string) (string, error) {
	result := ""
	err := r.db.QueryRow(`
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
		return "", fmt.Errorf("error while checking for duplicates: %w", err)
	}
	return result, nil
}

func (r *UserRepo) GetUserDataByID(id string) (*models.LoggedInUser, error) {
	user := models.LoggedInUser{}
	err := r.db.QueryRowx(`
      SELECT
        id, email, username, password_hash, email_verified, joined_at, active, profile_picture, email_verification_token
      FROM
        users
      WHERE 
        id = $1
    `,
		id,
	).Scan(
		&user.ID,
		&user.EncryptedEmail,
		&user.Username,
		&user.Password,
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
		return nil, fmt.Errorf("error while getting logged in user data: %w", err)
	}
	
	// Déchiffrer l'email
	if err := user.DecryptEmail(); err != nil {
		return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
	}
	
	return &user, nil
}

func (r *UserRepo) ChangePassword(password, id string) error {
    _, err := r.db.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", password, id)
    if err != nil {
        return fmt.Errorf("error while changing user password: %w", err)
    }
    return nil
}

func (r *UserRepo) ToggleStatus(users []string, status bool) error {
	_, err := r.db.Exec("UPDATE users SET active = $1 WHERE id = ANY($2)", status, pq.Array(users))
	if err != nil {
		return fmt.Errorf("error while changing accounts status: %w", err)
	}
	return nil
}

func (r *UserRepo) GetAllUsersData() (*[]models.LoggedInUser, error) {
	data := []models.LoggedInUser{}
	rows, err := r.db.Queryx(`
      SELECT
        id, email, username, password_hash, email_verified, joined_at, active, profile_picture, email_verification_token
      FROM
        users
    `)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving users data: %w", err)
	}
	for rows.Next() {
		user := models.LoggedInUser{}
		err := rows.Scan(
			&user.ID,
			&user.EncryptedEmail,
			&user.Username,
			&user.Password,
			&user.EmailVerified,
			&user.JoinedAt,
			&user.Active,
			&user.ProfilePicture,
			&user.EmailVerificationToken,
		)
		if err != nil {
			return nil, fmt.Errorf("error while retrieving users data: error while scanning row: %w", err)
		}
		
		// Déchiffrer l'email
		if err := user.DecryptEmail(); err != nil {
			return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
		}
		
		data = append(data, user)
	}
	return &data, nil
}

func (r *UserRepo) SetContactVerified(userId string) error {
	query := "UPDATE users SET email_verified = true WHERE id = $1"
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error while setting user contact to verified: %w", err)
	}
	return nil
}

func (r *UserRepo) SetEmailVerificationToken(userID, token string) error {
    query := "UPDATE users SET email_verification_token = $1 WHERE id = $2"
    _, err := r.db.Exec(query, token, userID)
    return err
}

func (r *UserRepo) GetUserByVerificationToken(token string) (*models.User, error) {
    var user models.User
    query := "SELECT * FROM users WHERE email_verification_token = $1"
    err := r.db.Get(&user, query, token)
    if err != nil {
        return nil, err
    }
    
    // Déchiffrer l'email
    if err := user.DecryptEmail(); err != nil {
        return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
    }
    
    return &user, nil
}

func (r *UserRepo) Update(user *models.User) error {
	// Chiffrer l'email avant d'enregistrer
	if err := user.EncryptEmail(); err != nil {
		return fmt.Errorf("erreur de chiffrement de l'email: %w", err)
	}
	
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

func (r *UserRepo) FindByVerificationToken(token string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, joined_at, active, profile_picture, email_verified, email_verification_token
		FROM users
		WHERE email_verification_token = $1
	`
	
	var user models.User
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.EncryptedEmail, &user.Username, &user.Password, &user.JoinedAt, &user.Active,
		&user.ProfilePicture, &user.EmailVerified, &user.EmailVerificationToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commons.Errors.ResourceNotFound
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// Déchiffrer l'email
	if err := user.DecryptEmail(); err != nil {
		return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
	}
	
	return &user, nil
}

func (r *UserRepo) VerifyEmail(token string) error {
	var user models.User

	// Trouver l'utilisateur par le token de vérification
	query := "SELECT * FROM users WHERE email_verification_token = $1"
	err := r.db.Get(&user, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commons.Errors.ResourceNotFound
		}
		return fmt.Errorf("could not find user by verification token: %w", err)
	}
	
	// Déchiffrer l'email
	if err := user.DecryptEmail(); err != nil {
		return fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
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

func (r *UserRepo) SetPasswordResetToken(email, token string) error {
    // Chiffrer l'email pour la recherche
    encryptedEmail, err := r.ee.EncryptEmail(email)
    if err != nil {
        return fmt.Errorf("erreur de chiffrement de l'email: %w", err)
    }
    
    query := "UPDATE users SET password_reset_token = $1 WHERE email = $2"
    result, err := r.db.Exec(query, token, encryptedEmail)
    if err != nil {
        return fmt.Errorf("error setting password reset token: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return commons.Errors.ResourceNotFound
    }
    
    return nil
}

// FindByPasswordResetToken trouve un utilisateur par son token de réinitialisation de mot de passe
func (r *UserRepo) FindByPasswordResetToken(token string) (*models.User, error) {
    query := `
        SELECT id, email, username, password_hash, joined_at, active, profile_picture, 
               email_verified, email_verification_token, password_reset_token
        FROM users
        WHERE password_reset_token = $1
    `
    
    var user models.User
    err := r.db.QueryRow(query, token).Scan(
        &user.ID, &user.EncryptedEmail, &user.Username, &user.Password, &user.JoinedAt, &user.Active,
        &user.ProfilePicture, &user.EmailVerified, &user.EmailVerificationToken, &user.PasswordResetToken,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, commons.Errors.ResourceNotFound
        }
        return nil, fmt.Errorf("error finding user by reset token: %w", err)
    }
    
    // Déchiffrer l'email
    if err := user.DecryptEmail(); err != nil {
        return nil, fmt.Errorf("erreur de déchiffrement de l'email: %w", err)
    }
    
    return &user, nil
}

// UpdatePassword met à jour le mot de passe d'un utilisateur
func (r *UserRepo) UpdatePassword(userID, newPasswordHash string) error {
    query := `UPDATE users SET password_hash = $1, password_reset_token = NULL WHERE id = $2`
    result, err := r.db.Exec(query, newPasswordHash, userID)
    if err != nil {
        return fmt.Errorf("error updating password: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return commons.Errors.ResourceNotFound
    }
    
    return nil
}