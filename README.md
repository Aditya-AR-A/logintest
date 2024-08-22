# Login System with Go and Fiber

## Description
This project is a simple login system built with Go and the Fiber web framework. It includes user registration, login, logout, password reset, and user management functionalities. The application uses a PostgreSQL database for data storage and bcrypt for password hashing.

## Prerequisites
- Go (version 1.16 or later)
- PostgreSQL
- Git

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/logintest.git
   cd logintest
   ```

2. Install the required Go packages:
   ```
   go mod tidy
   ```

   This will install the following packages:
   - github.com/gofiber/fiber/v2
   - github.com/lib/pq
   - golang.org/x/crypto/bcrypt

3. Set up the PostgreSQL database:
   - Install PostgreSQL if you haven't already
   - Create a new database for the project

4. Set up environment variables:
   Create a `.env` file in the project root with the following content:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_postgres_username
   DB_PASSWORD=your_postgres_password
   DB_NAME=your_database_name
   ```

## Database Setup

Run the following SQL command to create the necessary table:

```sql
CREATE TABLE login_credentials (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Running the Application

1. Start the server:
   ```
   go run server/main.go
   ```

2. The server will start on `http://localhost:3000`

## API Endpoints

- `POST /register`: Register a new user
- `POST /login`: Log in a user
- `POST /logout`: Log out a user
- `GET /users`: Get all users (protected route)
- `DELETE /users/:id`: Delete a user (protected route)
- `POST /reset-password`: Reset user password

## Frontend

The project includes simple HTML pages for user interaction:

- `/`: Main page (login and registration)
- `/users-page`: User management page
- `/reset-password-page`: Password reset page

## Additional Notes

- Make sure to handle CORS settings appropriately if you're running the frontend on a different domain or port.
- This project uses session-based authentication. Ensure that your production environment is configured to handle sessions securely.
- The project uses a simple in-memory session store. For production use, consider using a more robust session storage solution.

## Contributing

Feel free to fork this repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)