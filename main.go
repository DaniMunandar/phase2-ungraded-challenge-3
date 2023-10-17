package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"ungraded-challenge-3/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Koneksi ke database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/phase2_ngc3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Buat router HTTP
	router := httprouter.New()

	// Buat instance handler
	inventoryHandler := handler.NewInventoryHandler(db)

	// Atur routing untuk endpoint CRUD
	router.GET("/inventories", inventoryHandler.GetAllInventories)
	router.GET("/inventories/:id", inventoryHandler.GetInventory)
	router.POST("/inventories", inventoryHandler.CreateInventory)
	router.PUT("/inventories/:id", inventoryHandler.UpdateInventory)
	router.DELETE("/inventories/:id", inventoryHandler.DeleteInventory)

	// Jalankan server
	port := "8081" // Port server
	fmt.Printf("Server berjalan di port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
