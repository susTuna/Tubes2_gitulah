package database

func Define() {
	Exec(`
		CREATE TABLE elements (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			image_url VARCHAR(255)
		)
	`)

	Exec(`
		CREATE TABLE recipes (
			result_id INTEGER REFERENCES elements(id) ,
			dependency1_id INTEGER REFERENCES elements(id),
			dependency2_id INTEGER REFERENCES elements(id),
			PRIMARY KEY (result_id, dependency1_id, dependency2_id)
		)
	`)
}

func IsDefined() bool {
	var result bool

	QueryRow(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = 'public'
			AND tablename = 'elements'
		) AND EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = 'public'
			AND tablename = 'recipes'
		)
	`).Scan(&result)

	return result
}
