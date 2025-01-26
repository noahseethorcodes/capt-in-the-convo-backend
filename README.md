# CAPTInTheConvo Backend

Welcome to the backend for CAPTInTheConvo, a forum application API built using Go's Gin framework and PostgreSQL. This guide provides instructions on setting up the backend and PostgreSQL database locally, ensuring you can get started quickly. This application was developed for the NUS CVWO AY24/25 Winter Assignment.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Setup](#local-setup)
    - [Step 1: Clone the Repository](#step-1-clone-the-repository)
    - [Step 2: Set Up PostgreSQL](#step-2-set-up-postgresql)
    - [Step 3: Configure Environment Variables](#step-3-configure-environment-variables)
    - [Step 4: Run the Application](#step-4-run-the-application)
3. [Testing the Setup](#testing-the-setup)
5. [Frontend Setup](#frontend-setup)

---

## Prerequisites

Before getting started, ensure you have the following installed on your system:

- [Go](https://go.dev/dl/) (version 1.20 or later)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/)

---

## Local Setup

Follow these steps to set up the backend on your local machine.

### Step 1: Clone the Repository

1. Clone the repository from GitHub:
    ```bash
    git clone https://github.com/noahseethorcodes/capt-in-the-convo-backend.git
    cd capt-in-the-convo-backend
    ```

2. Switch to the development branch:
    ```bash
    git checkout dev
    ```

### Step 2: Set Up PostgreSQL

1. Start PostgreSQL on your local machine.

2. Log in to the PostgreSQL shell:
    ```bash
    psql -U postgres
    ```

3. Create a new database and user:
    ```sql
    CREATE DATABASE *YOUR_DB_NAME*;
    CREATE USER *YOUR_DB_USERNAME* WITH PASSWORD '*YOUR_DB_PASSWORD*';
    GRANT ALL PRIVILEGES ON DATABASE *YOUR_DB_NAME* TO *YOUR_DB_USERNAME*;
    \q
    ```

4. Note the connection details for configuring your backend.

### Step 3: Configure Environment Variables

1. Locate the `.env.example` file in the project root directory.

2. Fill in the required details:
    ```env
    DB_HOST=YOUR_HOST_NAME # usually localhost for local setup
    DB_USER=YOUR_DB_USERNAME
    DB_PASSWORD=YOUR_DB_PASSWORD
    DB_NAME=YOUR_DB_NAME
    DB_PORT=YOUR_DB_PORT # usually 5432
    DB_SSLMODE=disable # leave this as 'disable'
    APP_HOST=YOUR_APP_HOST # usually localhost for localsetup
    APP_PORT=YOUR_APP_PORT # usually 8080 for local setup
    SECRET_KEY=YOUR_SECRET_KEY # can be anything
    ```

3. Rename the `.env.example` file to `.env`:
    ```bash
    mv .env.example .env
    ```

### Step 4: Run the Application

1. Install dependencies:
    ```bash
    go mod download
    go mod tidy
    ```

2. Start the application:
    ```bash
    go run main.go
    ```

3. The backend should now be running on `http://localhost:8080`.

---

## Testing the Setup

To verify the setup:

1. Use a tool like `curl` or Postman to test the `/ping` endpoint:
    ```bash
    curl http://localhost:8080/ping
    ```

2. Expected response:
    ```json
    { "message": "pong" }
    ```
---

## Frontend Setup

The frontend for CAPTInTheConvo is a separate repository. Follow these steps to set up the frontend:

1. Clone the frontend repository:
    ```bash
    git clone https://github.com/noahseethorcodes/capt-in-the-convo-frontend.git
    cd capt-in-the-convo-frontend
    ```

2. Follow the setup instructions provided in the frontend repository's README.md file.

3. Ensure the frontend connects to the backend by setting the appropriate API URL in the frontend's environment variables.

For further details, visit the [CAPTInTheConvo Frontend Repository](https://github.com/noahseethorcodes/capt-in-the-convo-frontend).

---
