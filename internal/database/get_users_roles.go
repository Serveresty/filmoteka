package database

func GetUserRoles(id int) ([]string, error) {
	rows, err := DB.Query(`SELECT role_name FROM roles r JOIN users_roles ur ON ur.user_id=$1 AND r.role_id=ur.role_id;`, id)
	if err != nil {
		return []string{}, err
	}
	var roles []string
	for rows.Next() {
		var role string

		err = rows.Scan(&role)
		if err != nil {
			return []string{}, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
