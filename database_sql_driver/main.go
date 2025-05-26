package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "user:yourpassword@tcp(127.0.0.1:3306)/prasuna"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	} else {
		log.Print("successfully connected to db")
	}
	defer db.Close()

	// Check the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed:", err)
	}
	fmt.Println("Successfully connected to DB")

	// INSERT
	insertStmt := "INSERT INTO scores (message) VALUES (?)"
	res, err := db.Exec(insertStmt, "Hello from database/sql!")
	if err != nil {
		log.Fatal("Insert failed:", err)
	}
	lastID, _ := res.LastInsertId()
	fmt.Println("Inserted row with ID:", lastID)

	// UPDATE
	updateStmt := "UPDATE scores SET message = ? WHERE id = ?"
	_, err = db.Exec(updateStmt, "Updated message", lastID)
	if err != nil {
		log.Fatal("Update failed:", err)
	}
	fmt.Println("Updated row with ID:", lastID)

	someID := lastID
	_, err = db.Exec("DELETE FROM scores WHERE id = ?", someID)
	if err != nil {
		log.Fatal("Delete failed:", err)
	}
	fmt.Println("Deleted row with ID:", someID)
	//////////////////////////////////////////////////////////
	rows, err := db.Query("SELECT id, message FROM scores")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var msg string
		if err := rows.Scan(&id, &msg); err != nil {
			log.Println("Row scan failed:", err)
			continue
		}
		fmt.Printf("Row: ID=%d, Message=%s\n", id, msg)
	}
	////////////////////////////////
	stmt, err := db.Prepare("INSERT INTO scores (message) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	messages := []string{"msg1", "msg2", "msg3"}
	for _, msg := range messages {
		_, err := stmt.Exec(msg)
		if err != nil {
			log.Fatal("Batch insert failed:", err)
		}
	}
	//////////////////////////////////////////
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Begin transaction failed:", err)
	}

	_, err = tx.Exec("INSERT INTO scores (message) VALUES (?)", "okoko")
	if err != nil {
		tx.Rollback()
		log.Fatal("Insert in txn failed:", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Commit failed:", err)
	}
	fmt.Println("Transaction committed successfully")

	////////////////////JOINS//////////////////

	/* // SQL statement to create a table
		createTable := `
	    CREATE TABLE IF NOT EXISTS users (
	        id INT AUTO_INCREMENT PRIMARY KEY,
	        name VARCHAR(255) NOT NULL
	    );`

		// Execute the statement
		_, err = db.Exec(createTable)
		if err != nil {
			log.Fatal("Failed to create table:", err)
		} else {
			log.Println("Table created successfully")
		} */
	sqlStmt, _ := ioutil.ReadFile("init.sql")
	db.Exec(string(sqlStmt))

	// Inserting in a Loop (Flexible, but Slower)
	names := []string{"Alice", "Bob", "Charlie", "Dave"}
	stmt, err = db.Prepare("INSERT INTO users (name) VALUES (?)")
	if err != nil {
		log.Fatal("Prepare failed:", err)
	}
	defer stmt.Close()

	for _, name := range names {
		_, err := stmt.Exec(name)
		if err != nil {
			log.Println("Insert failed for", name, ":", err)
		}
	}
	fmt.Println("Batch insert via loop done")

	rows, err = db.Query(`
    SELECT s.id, s.message, u.name 
    FROM scores s
    JOIN users u ON s.id = u.id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var msg, username string
		rows.Scan(&id, &msg, &username)
		fmt.Println(id, msg, username)
	}

}
