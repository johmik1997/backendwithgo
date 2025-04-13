package db

import (
	"john/security"
	"log"
)
func seedDatabase() {
    // Create employees table
    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS employees (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL,
            is_admin BOOLEAN NOT NULL DEFAULT FALSE,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );
    `)
    if err != nil {
        log.Fatal("Error creating employees table:", err)
    }

    // Create employee_details table with proper foreign key
    _, err = DB.Exec(`
      CREATE TABLE IF NOT EXISTS employee_details (
    id SERIAL PRIMARY KEY,
    emp_id INTEGER NOT NULL,
    emp_name VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    experience INTEGER NOT NULL,
    address TEXT NOT NULL,
    birthdate DATE NOT NULL,
    employee_photo TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
    `)
    if err != nil {
        log.Fatal("Error creating employee_details table:", err)
    }

    // Create events table
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            date TIMESTAMP NOT NULL,
            location VARCHAR(255),
            host VARCHAR(255),
            description TEXT
        );
    `)
    if err != nil {
        log.Fatal("Error creating events table:", err)
    }

    // Create participants table with foreign keys
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS participants (
            id SERIAL PRIMARY KEY,
            event_id INTEGER REFERENCES events(id),
            user_id INTEGER REFERENCES employees(id),
            name VARCHAR(255) NOT NULL,
            class VARCHAR(100) DEFAULT 'Adventurer'
        );
    `)
    if err != nil {
        log.Fatal("Error creating participants table:", err)
    }

    // Insert employees first (with hashed passwords)
    users := []struct {
        username string
        password string
        isAdmin  bool
    }{
        {"melkam", "64382", true},
        {"mom", "7845", false},
        {"bar", "1234", false},
    }

    var employeeIDs []int
    for _, u := range users {
        // Hash the password
        hashed, err := security.HashPassword(u.password)
        if err != nil {
            log.Fatalf("Error hashing password for %s: %v", u.username, err)
        }

        // Insert employee and get the generated ID
        var id int
        err = DB.QueryRow(`
            INSERT INTO employees (username, password, is_admin)
            VALUES ($1, $2, $3)
            RETURNING id
        `, u.username, hashed, u.isAdmin).Scan(&id)
        if err != nil {
            log.Println("Insert employee failed:", err)
            continue
        }
        employeeIDs = append(employeeIDs, id)
    }

    // Only insert employee details if we have matching employee IDs
    if len(employeeIDs) >= 3 {
        employees := []struct {
            emp_id       int
            emp_name     string
            department   string
            experience   int
            address      string
            birthdate    string
            employePhoto string
        }{
            {employeeIDs[0], "Melkam", "Software Engineer", 5, "San Jose", "1990-04-15", "melkam_photo.jpg"},
            {employeeIDs[1], "Mom", "Product Manager", 3, "Mountain View", "1992-06-22", "mom_photo.jpg"},
            {employeeIDs[2], "Bar", "UX Designer", 4, "Palo Alto", "1991-08-09", "bar_photo.jpg"},
        }

        for _, e := range employees {
            _, err = DB.Exec(`
                INSERT INTO employee_details 
                (emp_id, emp_name, department, experience, address, birthdate, employePhoto)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
            `, e.emp_id, e.emp_name, e.department, e.experience, e.address, e.birthdate, e.employePhoto)
            if err != nil {
                log.Println("Insert employee details failed:", err)
            }
        }
    }

    // Insert sample events
    events := []struct {
        name        string
        date        string
        location    string
        host        string
        description string
    }{
        {
            "Quarterly Guild Meeting",
            "2025-05-15 14:00:00",
            "Great Hall",
            "Elder Council",
            "Discuss guild affairs and upcoming quests",
        },
        {
            "Magic Training Workshop",
            "2025-05-20 10:00:00",
            "Arcane Tower",
            "Master Wizard",
            "Learn new spells and improve your magical abilities",
        },
    }

    for _, e := range events {
        var eventID int
        err = DB.QueryRow(`
            INSERT INTO events (name, date, location, host, description)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `, e.name, e.date, e.location, e.host, e.description).Scan(&eventID)
        if err != nil {
            log.Println("Insert event failed:", err)
            continue
        }

        // Register some participants for each event
        if len(employeeIDs) > 0 {
            _, err = DB.Exec(`
                INSERT INTO participants (event_id, user_id, name, class)
                VALUES ($1, $2, $3, $4)
            `, eventID, employeeIDs[0], "Melkam", "Wizard")
            if err != nil {
                log.Println("Insert participant failed:", err)
            }
        }
    }

    // Additional sample events (fixed syntax)
    _, err = DB.Exec(`
        INSERT INTO events (name, date, location, host, description)
        VALUES 
            ('Quarterly Guild Meeting', NOW() + INTERVAL '7 days', 'Great Hall', 'Elder Council', 'Discuss guild affairs'),
            ('Magic Training', NOW() + INTERVAL '14 days', 'Arcane Tower', 'Master Wizard', 'Learn new spells')
        ON CONFLICT DO NOTHING;
    `)
    if err != nil {
        log.Println("Error seeding events:", err)
    }
}