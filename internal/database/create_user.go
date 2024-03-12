package database

import "filmoteka/internal/models"

func CreateNewUser(user *models.User) error {
	_, err := DB.Exec(
		`INSERT INTO "users_data" (first_name, last_name, email, password) VALUES ($1,$2,$3,$4);`,
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`INSERT INTO users_roles (user_id, role_id) SELECT ud.user_id, r.role_id FROM users_data ud JOIN roles r ON ud.email = $1 AND r.role_name = $2;`, user.Email, "user")
	if err != nil {
		return err
	}
	return nil
}
