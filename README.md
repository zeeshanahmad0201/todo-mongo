# Todo App using GO

![Todo App](https://via.placeholder.com/800x400.png?text=Go+Todo+App)

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Technologies](#technologies-used)
- [Contributing](#contributing)
- [Contact](#contact)

## Features

- Create, read, update, and delete todos
- Persistent storage with MongoDB
- RESTful API
- Lightweight and fast
- Simple and clean architecture

## Getting Started

### Prerequisites

- Go 1.16 or higher
- MongoDB

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/zeeshanahmad0201/go-todo-app.git
    cd go-todo-app
    ```

2. **Install the dependencies:**
    ```sh
    go mod tidy
    ```

3. **Run the application:**
    ```sh
    go run main.go
    ```

The application will be available at `http://localhost:8080`.

## Usage

You can interact with the API using a tool like `curl` or Postman. Below are the available endpoints.

## API Endpoints

 - ```POST``` /login
 - ```POST``` /signup
 - ```GET``` /todos
 - ```POST``` /todos
 - ```GET``` /todos/{id}
 - ```PUT``` /todos/{id}
 - ```DELETE```  /todos/{id}

## Project Structure

- **Controller**: Handles the incoming HTTP requests.
  - `todo_controller.go`
  - `user_controller.go`
- **Service**: Contains the business logic.
  - `todo_service.go`
  - `user_service.go`
- **Router**: Defines the routes for the application.
  - `router.go`
- **Common**: Utility and helper functions.
  - `error_handler.go`
  - `timestamp_handler.go`
- **Database**: Handles the connection with MongoDB.
  - `mongo_db_connection.go`
- **Helpers**: Additional helper functions.
  - `password_helper.go`
  - `token_helper.go`
- **Models**: Defines the data models.
  - `user.go`
  - `todo.go`
- **Main**: Entry point of the application.
  - `main.go`

## Technologies Used

- Golang: The main programming language for the application
- MongoDB: The NoSQL database used for storing tasks
- Gorilla Mux: A powerful HTTP router and URL matcher for building Go web servers

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request.

## Contact

Created by [Zeeshan Ahmad](mailto:xeeshan.ahmad.dev@gmail.com) - feel free to contact me!
