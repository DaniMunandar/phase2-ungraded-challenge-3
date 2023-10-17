package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Inventory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ItemCode    string `json:"item_code"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type InventoryHandler struct {
	db *sql.DB
}

func NewInventoryHandler(db *sql.DB) *InventoryHandler {
	return &InventoryHandler{db}
}

func (ih *InventoryHandler) GetAllInventories(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := ih.db.Query("SELECT id, name, item_code, stock, description, status FROM inventories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	inventories := []Inventory{}
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(&inv.ID, &inv.Name, &inv.ItemCode, &inv.Stock, &inv.Description, &inv.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		inventories = append(inventories, inv)
	}

	response, err := json.Marshal(inventories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (ih *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	row := ih.db.QueryRow("SELECT id, name, item_code, stock, description, status FROM inventories WHERE id = ?", id)

	var inv Inventory
	err := row.Scan(&inv.ID, &inv.Name, &inv.ItemCode, &inv.Stock, &inv.Description, &inv.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (ih *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse data dari body request
	var inv Inventory
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert data inventory baru ke dalam database
	result, err := ih.db.Exec("INSERT INTO inventories (name, item_code, stock, description, status) VALUES (?, ?, ?, ?, ?)",
		inv.Name, inv.ItemCode, inv.Stock, inv.Description, inv.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Dapatkan ID yang baru saja dibuat
	newID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]int64{"id": newID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (ih *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// Parse data dari body request
	var inv Inventory
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update data inventory dalam database
	_, err := ih.db.Exec("UPDATE inventories SET name = ?, item_code = ?, stock = ?, description = ?, status = ? WHERE id = ?", inv.Name, inv.ItemCode, inv.Stock, inv.Description, inv.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Inventory updated successfully")
}

func (ih *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	_, err := ih.db.Exec("DELETE FROM inventories WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Inventory deleted successfully")
}
