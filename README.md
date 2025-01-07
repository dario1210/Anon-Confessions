# **Project Title: Anon-Confessions**

A web application for anonymous confessions, providing a platform for open and anonymous communication.

---

## **Prerequisites**

### **Required**

- **Go**: [Install Go](https://go.dev/)

### **Optional (For Developers)**

- **SQLite**: Required for database operations.
- **Migrate**: For database migrations. Install via [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4).
- **Make**: Often pre-installed on Linux. If not, install it using:

  ```bash
  sudo apt install make      # Debian/Ubuntu
  sudo yum install make      # CentOS/RHEL
  sudo pacman -S make        # Arch-based systems
  ```

---

## **Run a Migration**

To create a new migration file, use the following command:

```bash
migrate create -ext sql -dir ./cmd/internal/db/migrations_files/ -seq create_user_table
```

---

## **Installation and Setup**

Follow these steps to set up and run the project:

### 1. **Clone the Repository**

```bash
git clone https://github.com/dario1210/Anon-Confessions.git
cd Anon-Confessions
```

### 2. **Install Dependencies**

```bash
go mod tidy
```

### 3. **Run Database Migrations**

#### Using `make` (Recommended)

```bash
make migrations
```

#### Without `make`

```bash
go run cmd/internal/db/migrate/migration.go
```

### 4. **Run the Seeder**

Seed the database with sample data for testing.

#### Using `make`

```bash
make seed
```

#### Without `make`

```bash
go run cmd/internal/db/seeder/seeder.go
```

The seeder populates the database with sample authentication accounts, which can be used for testing the API:

- **1234567891234567**
- **3998442793406687**
- **7180218105191773**
- **6129856725721562**

### 5. **Generate Swagger Documentation** (Optional)

To use Swagger UI for API documentation and testing, generate the Swagger docs:

#### Using `make`

```bash
make swagger
```

#### Without `make`

```bash
swag init -g cmd/internal/app/app.go
```

### 6. **Run the Application**

Run the application to start the server.

#### Using `make`

```bash
make run
```

#### Without `make`

```bash
go run cmd/server/main.go
```

### Access the Application

Visit the following URL in your browser:

[`http://localhost:9000/swagger/index.html#/`](http://localhost:9000/swagger/index.html#/)
