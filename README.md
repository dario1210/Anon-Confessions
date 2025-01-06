# **Project Title: Anon-Confessions**

A web application for anonymous confessions. This project allows users to share their thoughts anonymously, providing a platform for open communication.

---

## **Prerequisites**

To run the project, ensure you have the following installed on your machine:

### **Required:**

- **Go** [Install Go](https://go.dev/)

### **Optional (For Developers):**

- **SQLite**: Required for database operations.
- **Migrate**: For running database migrations. Install it from [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4).
- **Make**: Commonly pre-installed on many Linux distributions. If not, you can install it:
  ```bash
  sudo apt install make      # For Debian/Ubuntu
  sudo yum install make      # For CentOS/RHEL
  sudo pacman -S make        # For Arch-based systems
  ```

---

## **Installation**

Follow these steps to install and set up the project:

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/dario1210/Anon-Repository.git
   ```

2. **Navigate to the Project Directory:**

   ```bash
   cd Anon-Repository
   ```

3. **Install Dependencies:**

   ```bash
   go mod tidy
   ```

4. **Database Migrations:**

   **Option 1 (If you have `make`):**

   ```bash
   make migrations
   ```

   **Option 2 (Manual Commands):**
   If `make` is not available on your machine, run the following command directly:

   ```bash
   go run cmd/internal/db/migrate/migration.go
   ```

5. **Run the Application:**

   **Option 1 (If you have `make`):**

   ```bash
   make run
   ```

   **Option 2 (Manual Commands):**

   ```bash
   go run cmd/server/main.go
   ```

6. **Swagger Documentation (Optional):**
   If you wish to generate Swagger documentation for the project:
   ```bash
   swag init -g cmd/server/main.go
   ```
