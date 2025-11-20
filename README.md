# 📘 Notes Management System (Go + Gin + MySQL)

A clean, production-style **Notes Management System** built using **Golang**, **Gin Web Framework**, and **MySQL**, paired with a simple **HTML + JavaScript frontend**.

This project demonstrates real backend development skills — including database integration, middleware, REST API design, and a functional mini-frontend using `fetch()`.

---

## ⭐ Features

### 🔹 Backend (Go + Gin)
- ✔ Full CRUD operations (Create, Read, Update, Delete)  
- ✔ MySQL Database using **GORM ORM**  
- ✔ Auto-migrated DB tables  
- ✔ JSON request & response  
- ✔ Validation & error handling  
- ✔ Custom **CORS Middleware**  
- ✔ Custom **Request Logging Middleware**

### 🔹 Frontend (HTML + JavaScript)
- ✔ Add notes  
- ✔ Edit notes  
- ✔ Delete notes  
- ✔ View notes  
- ✔ Uses `fetch()` API  
- ✔ Fully connected to backend  

### 🔹 Tools & Technologies
`Go, Gin, MySQL, GORM, JavaScript, HTML, CSS, Postman`

---
go-notes-api/
│
├── main.go # Backend (API + DB + Middleware)
├── frontend.html # UI to interact with API
├── go.mod
├── go.sum
---
## 📂 Project Structure
## 🗄️ Database Setup (MySQL)

Run these commands inside **MySQL Workbench**:

```sql
CREATE DATABASE go_notes;

CREATE USER 'go_user'@'localhost' IDENTIFIED BY 'password123';

GRANT ALL PRIVILEGES ON go_notes.* TO 'go_user'@'localhost';

FLUSH PRIVILEGES;
```
Update the DSN in main.go if needed:

go_user:password123@tcp(127.0.0.1:3306)/go_notes

🔌 API Endpoints

➤ Get all notes

```Bash
GET /notes
```
➤ Create a note

```Bash
POST /notes
```

Body:

```json
{
  "title": "My Note",
  "content": "This is a new note"
}
```
➤ Get note by ID

```Bash
GET /notes/:id
```

➤ Update note

```Bash
PUT /notes/:id
```

➤ Delete note

```Bash
DELETE /notes/:id
```

🖥️ How to Run

1️⃣ Start server

```Bash
go run main.go
```

```Bash
http://localhost:8080/
```

You can now:

- Add notes
- Edit notes
- Delete notes
- Auto-refresh list


🧑‍💻 Author

Vykaas Verma

Backend & Full-Stack Developer

Skills: Go, Python, JavaScript, SQL, REST APIs






