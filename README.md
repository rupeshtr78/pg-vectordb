# Project Title

## Overview
This repository is a Go-based application designed for managing embeddings stored in a PostgreSQL database. It uses Docker for containerization and includes a PostgreSQL setup for development and testing.

## Features
- Connect to a PostgreSQL database
- Fetch and load embeddings
- Use GORM for ORM capabilities 

## Getting Started

### Prerequisites
- Docker and Docker Compose installed on your local machine.
- Go programming environment set up (if you wish to run or modify the code).

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Start the database using Docker Compose:
    Build docker image from [pgvector repo](https://github.com/pgvector/pgvector.git)
    Run the docker compose.
   ```bash
   docker-compose up
   ```

3. Initialize the database schema:
   The SQL schema is located in `init.sql`. Ensure to run it against your PostgreSQL instance or let the application execute it.

### Project Structure
- `cmd/main.go`: The entry point of the application. Contains the `main` function to bootstrap the service.
- `docker-compose.yml`: Configuration file for Docker services.
- `go.mod`: Go module file for managing dependencies.
- `init.sql`: SQL script to initialize the database.
- `scratch/`: Contains standalone Go files, such as database connection configurations.
- `internal/pg_embed/`: Core service logic, divided into various functionality:
  - `constants.go`: Contains constant definitions used throughout the service.
  - `fetch_embeddings.go`: Logic for fetching embeddings from the database.
  - `load_data.go`: Logic for loading data into the database.
  - `pg_db_gorm.go`: GORM setup and database configuration using "gorm.io/gorm".
  - `query_vectordb.go`: Query functions specifically aimed at interacting with the vector database.

### Usage
To run the application, 

```bash
go run cmd/main.go
```
You can expand the functionalities by modifying the files in the `internal/pg_embed/` directory.

### Contributing
If you wish to contribute to this project, please fork the repository and make a pull request. Ensure you adhere to the coding standards and provide clear commit messages.

### License
This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contact
For inquiries or issues, please reach out to the project maintainer at [email@example.com].

