package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define a model struct
type Remark struct {
	ID      int
	Message string
	UserID  int
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
	db.AutoMigrate(&Remark{})

	// INSERT
	newRemark := Remark{Message: "Hello from GORM!", UserID: 1}
	result := db.Create(&newRemark)
	if result.Error != nil {
		log.Fatal("Insert failed:", result.Error)
	}
	fmt.Println("Inserted row with ID:", newRemark.ID)

	// UPDATE
	newRemark.Message = "Updated via GORM"
	db.Save(&newRemark)
	fmt.Println("Updated row with ID:", newRemark.ID)

	someID := newRemark.ID
	db.Delete(&Remark{}, someID) // deletes by primary key

	var scores []Remark
	result = db.Find(&scores)
	if result.Error != nil {
		log.Fatal("Query failed:", result.Error)
	}

	for _, score := range scores {
		fmt.Printf("Row: ID=%d, Message=%s\n", score.ID, score.Message)
	}
	scores = []Remark{
		{Message: "msg11", UserID: 1},
		{Message: "msg21", UserID: 2},
		{Message: "msg31", UserID: 3},
	}
	db.Create(&scores)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Remark{Message: "inside GORM txn", UserID: 23}).Error; err != nil {
			return err
		}
		return nil // commit
	})
	if err != nil {
		log.Fatal("GORM transaction failed:", err)
	}
	// /////////JOINS///////////////

}
