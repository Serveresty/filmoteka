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
	CREATE TABLE IF NOT EXISTS "films" (film_id serial PRIMARY KEY, name VARCHAR(150) NOT NULL, description VARCHAR(1000), release_date DATE, rate NUMERIC(3,1) CHECK (rate >= 0 AND rate <= 10));
	CREATE TABLE IF NOT EXISTS "actors" (actor_id serial PRIMARY KEY, name VARCHAR(50), gender CHAR(1), birth_date DATE);
	CREATE TABLE IF NOT EXISTS "actors_films" (actor_id INT references actors (actor_id) on delete cascade, film_id INT references films (film_id) on delete cascade, PRIMARY KEY (actor_id, film_id));
	`)
	if err != nil {
		return err
	}

	_, _ = DB.Exec(`INSERT INTO "roles" ("role_name") VALUES ('user'), ('admin');`)
	return nil
}
