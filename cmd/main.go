package main

import (
    "log"
    "net/http"
    "os"
    "fmt"

    "github.com/gorilla/mux"
    "product-service/internal/db"
    "product-service/internal/handlers"
)

func main() {
    database := db.Connect()
    defer database.Close()

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Product service is running")
    }).Methods("GET")
    r.HandleFunc("/admin/products", handlers.CreateProductHandler(database)).Methods("POST")
    r.HandleFunc("/admin/products/bulk", handlers.BulkImportProductsHandler(database)).Methods("PUT")
    r.HandleFunc("/products", handlers.GetProductsHandler(database)).Methods("GET")
    r.HandleFunc("/products/{id}", handlers.GetProductByIDHandler(database)).Methods("GET")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on port %s...", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
