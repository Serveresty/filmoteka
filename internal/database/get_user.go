package database

func GetUser(email string) (int, string, error) {
	var id int
	var pass string
	row := DB.QueryRow(`SELECT user_id, password FROM "users_data" WHERE email=$1`, email)
	err := row.Scan(&id, &pass)
	if err != nil {
		return 0, "", err
	}
	return id, pass, nil
}
