# 💸 EasySplit

A backend service that helps users easily track, manage, and split expenses with friends. Whether you're planning a trip, sharing a household, or managing group events, EasySplit makes it simple to keep everyone on the same page financially.

## ✨ Features

- ✅ User registration and authentication
- 👥 Create and manage groups
- 💰 Add, split, and track expenses
- 📊 Track who owes whom and how much
- 📬 Send and accept friend requests
- 🔄 Supports equal and unequal splits
- 🗂 Detailed APIs for front-end integration

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT
- **Framework:** Gorilla Mux
- **Validation:** go-playground/validator

## 🧪 Getting Started

1. **Clone the repository:**

    git clone https://github.com/Basu008/EasySplit.git
    cd EasySplit

2. **Setting up config file:**
    - Go to conf folder in root directory
    - create a new file as default.toml
    - refer to example.toml and add configuration as per your need.

3. **Installing dependency:**
    go mod tidy

4. **Running Code:**
    go run main.go

5. **Sanity Check**
    Hit GET {{baseURL}}/health-check. If everything is working fine, you should get `true` in payload
