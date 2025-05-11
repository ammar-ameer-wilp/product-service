package handlers

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "database/sql"

    "github.com/jmoiron/sqlx"
    "github.com/gorilla/mux"
    "github.com/yourusername/product-service/internal/models"
)

func CreateProductHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        var p models.Product
        if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        query := `INSERT INTO products (name, description, price, sku) VALUES ($1, $2, $3, $4)`
        _, err := db.Exec(query, p.Name, p.Description, p.Price, p.SKU)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
    }
}

func BulkImportProductsHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var products []models.Product
        if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        tx := db.MustBegin()
        for _, p := range products {
            _, err := tx.Exec(`INSERT INTO products (name, description, price, sku) VALUES ($1, $2, $3, $4)`,
                p.Name, p.Description, p.Price, p.SKU)
            if err != nil {
                tx.Rollback()
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
        if err := tx.Commit(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
    }
}

func GetProductsHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        q := r.URL.Query()

        search := q.Get("search")
        minPrice := q.Get("min_price")
        maxPrice := q.Get("max_price")
        sortBy := q.Get("sort_by")
        order := q.Get("order")
        page := q.Get("page")
        limit := q.Get("limit")

        sortFields := map[string]bool{"name": true, "price": true}
        if !sortFields[sortBy] {
            sortBy = "name"
        }
        if order != "desc" {
            order = "asc"
        }

        pageVal, _ := strconv.Atoi(page)
        limitVal, _ := strconv.Atoi(limit)
        if pageVal < 1 {
            pageVal = 1
        }
        if limitVal < 1 {
            limitVal = 10
        }
        offset := (pageVal - 1) * limitVal

        baseQuery := `SELECT * FROM products WHERE 1=1`
        args := []interface{}{}
        idx := 1

        if search != "" {
            baseQuery += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", idx, idx+1)
            args = append(args, "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
            idx += 2
        }
        if minPrice != "" {
            baseQuery += fmt.Sprintf(" AND price >= $%d", idx)
            args = append(args, minPrice)
            idx++
        }
        if maxPrice != "" {
            baseQuery += fmt.Sprintf(" AND price <= $%d", idx)
            args = append(args, maxPrice)
            idx++
        }

        baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, order, idx, idx+1)
        args = append(args, limitVal, offset)

        var products []models.Product
        err := db.Select(&products, baseQuery, args...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(products)
    }
}

func GetProductByIDHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        id := mux.Vars(r)["id"]
        var product models.Product
        err := db.Get(&product, "SELECT * FROM products WHERE id = $1", id)
        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                http.Error(w, "Product not found", http.StatusNotFound)
            } else {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }
        json.NewEncoder(w).Encode(product)
    }
}
