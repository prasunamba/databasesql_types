package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define a model struct
type Score struct {
	ID      int
	Message string
}

func main() {
	dsn := "user:yourpassword@tcp(127.0.0.1:3306)/prasuna?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Print("successfully connected to db")
	}

	// Auto-migrate the schema (optional)
	db.AutoMigrate(&Score{})

	// INSERT
	newScore := Score{Message: "Hello from GORM!"}
	result := db.Create(&newScore)
	if result.Error != nil {
		log.Fatal("Insert failed:", result.Error)
	}
	fmt.Println("Inserted row with ID:", newScore.ID)

	// UPDATE
	newScore.Message = "Updated via GORM"
	db.Save(&newScore)
	fmt.Println("Updated row with ID:", newScore.ID)

	someID := newScore.ID
	db.Delete(&Score{}, someID) // deletes by primary key

	var scores []Score
	result = db.Find(&scores)
	if result.Error != nil {
		log.Fatal("Query failed:", result.Error)
	}

	for _, score := range scores {
		fmt.Printf("Row: ID=%d, Message=%s\n", score.ID, score.Message)
	}
	scores = []Score{
		{Message: "msg11"},
		{Message: "msg21"},
		{Message: "msg31"},
	}
	db.Create(&scores)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Score{Message: "inside GORM txn"}).Error; err != nil {
			return err
		}
		return nil // commit
	})
	if err != nil {
		log.Fatal("GORM transaction failed:", err)
	}

}
