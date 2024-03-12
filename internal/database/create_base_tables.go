package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateBaseTables() error {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS "users_data" (user_id serial PRIMARY KEY, first_name VARCHAR(50) NOT NULL, last_name VARCHAR(50) NOT NULL, email VARCHAR(255) UNIQUE, password VARCHAR(255) NOT NULL);
	CREATE TABLE IF NOT EXISTS "roles" (role_id serial PRIMARY KEY, role_name VARCHAR(20) UNIQUE NOT NULL);
	CREATE TABLE IF NOT EXISTS "users_roles" (user_id int references users_data (user_id) on delete cascade, role_id int references roles (role_id) on delete cascade);
	`)
	if err != nil {
		return err
	}

	_, _ = DB.Exec(`INSERT INTO "roles" ("role_name") VALUES ('user'), ('admin');`)
	return nil
}
