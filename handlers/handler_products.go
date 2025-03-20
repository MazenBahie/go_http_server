package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MazenBahie/go_http_server/models"
	"github.com/gorilla/mux"
)

func HandleGetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// mentioned columns specifically to avoid any future bugs during scanning the rows
		query := `
           SELECT id, name, price, description FROM products
        `
		rows, err := db.Query(query)

		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error getting products")
			return
		}
		defer rows.Close()

		products := []models.Product{}
		for rows.Next() {
			product := models.Product{}
			err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description)
			if err != nil {
				continue
			}
			products = append(products, product)
		}

		ResponseJson(w, http.StatusOK, products)
	}
}

func HandleCreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		product := models.Product{}
		json.NewDecoder(r.Body).Decode(&product)

		err := Val.Struct(product)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Validation error "+err.Error())
			return
		}

		createdProduct := models.Product{}

		query := `
            INSERT INTO products(name, price, description)
			VALUES ($1,$2,$3)
			RETURNING id, name, price, description;
        `
		err = db.QueryRow(query, product.Name, product.Price, product.Description).Scan(
			&createdProduct.ID, &createdProduct.Name, &createdProduct.Price, &createdProduct.Description,
		)

		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error adding product")
			return
		}

		ResponseJson(w, http.StatusCreated, createdProduct)
	}
}

func HandleUpdateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["product_id"]
		if id == "" {
			ResponseWithError(w, http.StatusBadRequest, "product_id is required")
			return
		}

		castedId, err := strconv.Atoi(id)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "product_id must be a number")
			return
		}

		query := `
           SELECT id FROM products
		   WHERE id = $1
        `
		rowChecker := 0
		row := db.QueryRow(query, castedId)
		row.Scan(&rowChecker)

		if row.Err() == sql.ErrNoRows || rowChecker == 0 {
			ResponseWithError(w, http.StatusNotFound, "Product not found")
			return
		} else if row.Err() != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error getting product: "+row.Err().Error())
			return
		}

		updatedProduct := models.Product{}
		json.NewDecoder(r.Body).Decode(&updatedProduct)

		err = Val.Struct(updatedProduct)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
			return
		}

		updateQuery := `
            UPDATE products
            SET name = $1, price = $2, description = $3
            WHERE id = $4
            RETURNING id, name, price, description
        `
		err = db.QueryRow(updateQuery, updatedProduct.Name, updatedProduct.Price,
			updatedProduct.Description, castedId).
			Scan(
				&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Price, &updatedProduct.Description,
			)
		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error updating product: "+err.Error())
			return
		}
		ResponseJson(w, http.StatusOK, updatedProduct)
	}
}

func HandleDeleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := mux.Vars(r)["product_id"]
		if id == "" {
			ResponseWithError(w, http.StatusBadRequest, "product_id is required")
			return
		}

		castedId, err := strconv.Atoi(id)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "product_id must be a number")
			return
		}

		query := `
            DELETE FROM products
			WHERE id = $1
        `
		_, err = db.Exec(query, castedId)

		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error delete product:"+err.Error())
			return
		}

		ResponseJson(w, http.StatusOK, "")
	}
}
