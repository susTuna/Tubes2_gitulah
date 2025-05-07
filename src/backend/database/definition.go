package database

func Define() {
	Exec(`
		CREATE TABLE Elements (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			image_url VARCHAR(255)
		)
	`)

	Exec(`
		CREATE TABLE Recipes (
			result_id INTEGER REFERENCES Elements(id) ,
			dependency1_id INTEGER REFERENCES Elements(id),
			dependency2_id INTEGER REFERENCES Elements(id),
			PRIMARY KEY (result_id, dependency1_id, dependency2_id)
		)
	`)
}

func IsDefined() bool {
	var result bool

	QueryRow(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = "public"
			AND tablename = "Elements"
		) AND EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = "public"
			AND tablename = "Recipes"
		)
	`).Scan(&result)

	return result
}
