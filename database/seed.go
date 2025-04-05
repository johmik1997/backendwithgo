package db

func seedDatabase() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS employee_details (
			emp_id VARCHAR(50) PRIMARY KEY,
			emp_title VARCHAR(100),
			address TEXT
		);

		INSERT INTO employees (username, password) 
		VALUES ('Gordon', '64382'), ('Nick', '7845'),('John', '1234')
		ON CONFLICT (username) DO NOTHING;

		INSERT INTO employee_details (emp_id, emp_title, address)
		VALUES ('78', 'Software Engineer', 'San Jose')
		ON CONFLICT (emp_id) DO NOTHING;
	`)
	if err != nil {
		panic("Seeding failed: " + err.Error())
	}
}
