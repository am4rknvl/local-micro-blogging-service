# local-micro-blogging-service
# 📝 Microblogly — Minimalist Microblogging Platform

A fast, Go-powered microblogging and social networking service built from scratch.  
Think of it as a lean, hacker-friendly Twitter clone with core features, scalable backend, and full ownership.

---

## 🚀 Features (WIP)

- ✅ **User Authentication**
  - Signup / Login / Logout
  - Password hashing and session management
- 🔐 **JWT-based Auth Middleware**
- 🧠 **User Profiles**
  - Avatar, bio, update profile info
- 📝 **Posts (aka Microblogs)**
  - Create, delete, like, reply
- 📍 **User Interactions**
  - Follow / unfollow
  - Block users
- 🔔 **Notifications System**
- 📡 **Feed Algorithm**
  - Show posts from followed users + explore
- 💬 **Real-time Messaging** *(Upcoming)*
- 🧰 Built with clean architecture & RESTful APIs

---

## 🛠️ Tech Stack

- **Backend:** Go (Golang)
- **Router:** Gorilla Mux
- **Database:** PostgreSQL
- **Auth:** JWT, bcrypt
- **ORM:** Raw SQL / pgx (or GORM if you used it)
- **Hosting:** Railway / Fly.io / Render (WIP)
- **Frontend:**  (the mvp is built with expo/reactnative but since the project is built with Go im considering using flutter)

---

## 🧑‍💻 Getting Started

```bash
# 1. Clone the repo
git clone https://github.com/am4rkncl/local-micro-blogging-service.git
cd local-micro-blogging service

# 2. Setup .env
cp .env.example .env  # Fill in your DB_URL and JWT_SECRET

# 3. Run DB (if using Docker)
docker-compose up -d

# 4. Run the server
go run main.go
