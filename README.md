# Product Service

This is a Golang-based Product Service microservice for managing products, built as part of a scalable services assignment. It includes admin and customer-facing APIs, PostgreSQL database support, Dockerized setup, and Swagger documentation.

## Features

- **Admin APIs**
  - Create a product
  - Bulk import products

- **Customer APIs**
  - Get list of products with search, filter, sort, pagination
  - Get product details by ID

## Tech Stack

- Golang
- PostgreSQL
- Docker & Docker Compose
- Gorilla Mux
- SQLX

## Project Structure

```
product-service/
├── cmd/                    
├── internal/
│   ├── db/                 
│   ├── handlers/           
│   └── models/             
├── migrations/             
├── swagger/                
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Setup Instructions

1. **Clone the repository**
   ```
   git clone https://github.com/ammar-ameer-wilp/product-service.git
   cd product-service
   ```

2. **Run the service using Docker Compose**
   ```
   docker-compose up --build
   ```

3. **Service will be available at**
   ```
   http://localhost:8080
   ```

### Environment Variables

These are set in `docker-compose.yml`:

- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=user`
- `DB_PASSWORD=password`
- `DB_NAME=productdb`

### Migrations

Migrations are automatically run using the SQL file inside `migrations/`.

## API Documentation

The Swagger documentation is available in the `swagger/swagger.yaml` file. You can view it locally using a Swagger editor or import it into [https://editor.swagger.io](https://editor.swagger.io).

## License

This project is for educational purposes only.
