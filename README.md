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

#### Using `make` (Recommended):

```bash
make migrations
```

#### Without `make`:

```bash
go run cmd/internal/db/migrate/migration.go
```

### 4. **Run the Seeder**

#### Using `make`:

```bash
make seed
```

#### Without `make`:

```bash
go run cmd/internal/db/seeder/seeder.go
```

Seeder populates the database with sample authentication accounts:
Can be used for testing the API.

- **1234567891234567**
- **3998442793406687**
- **7180218105191773**
- **6129856725721562**

### 5. **Run the Application**

#### Using `make`:

```bash
make run
```

#### Without `make`:

```bash
go run cmd/server/main.go
```

### 6. **Generate Swagger Documentation** (Optional)

If you want to use Swagger UI for API documentation and testing:

1. Run the following commands:

#### Using `make`:

```bash
make swagger
```

#### Without `make`:

```bash
	swag init -g cmd/internal/app/app.go
```
