package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Product struct {
	ProductId          string `json:"product_id"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/product")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/products", getAllProduct).Methods("GET")
	router.HandleFunc("/products", addProduct).Methods("POST")
	router.HandleFunc("/products/{id}", getDetailProduct).Methods("GET")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	http.ListenAndServe(":10000", router)
}
func getAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product
	result, err := db.Query("SELECT product_id, product_name, product_description from ms_product")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var product Product
		err := result.Scan(&product.ProductId, &product.ProductName, &product.ProductDescription)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add Product")
	w.Header().Set("Content-Type", "application/json")
	query := "INSERT INTO `ms_product` (`product_name`, `product_description`) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	product_name := keyVal["product_name"]
	product_description := keyVal["product_description"]
	_, err = stmt.Exec(product_name, product_description)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product Added")
}

func getDetailProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT product_id, product_name, product_description FROM ms_product WHERE product_id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var product Product
	for result.Next() {
		err := result.Scan(&product.ProductId, &product.ProductName, &product.ProductDescription)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("UPDATE ms_product SET product_name = ?, product_description = ? WHERE product_id = ?")
	params := mux.Vars(r)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newname := keyVal["product_name"]
	newdesc := keyVal["product_description"]
	_, err = stmt.Exec(newname, newdesc, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product with ProductId = %s was updated", params["id"])
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM ms_product WHERE product_id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product with ProductId = %s was deleted", params["id"])
}
