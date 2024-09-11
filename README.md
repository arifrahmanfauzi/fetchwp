# README

## Overview

This Go program fetches data from a WordPress REST API, processes the JSON response, and inserts the data into a MySQL database. The program uses the `net/http` package for making HTTP requests, the `encoding/json` package for unmarshalling JSON, and the `database/sql` package for interacting with MySQL.

## Installation

1. Install Go: https://golang.org/doc/install
2. Install the required packages:
    - `github.com/go-sql-driver/mysql` for MySQL driver: `go get -u github.com/go-sql-driver/mysql`
3. Clone this repository: `git clone https://github.com/your-username/your-repository.git`
4. Navigate to the project directory: `cd your-repository`

## Usage

1. Update the MySQL connection string in the `main` function with your own credentials.
2. Run the program: `go run main.go`

## Notes

- The program assumes that the MySQL database is already created and has the necessary tables.
- The program fetches data from the WordPress REST API using a GET request with the specified query parameters.
- The JSON response is unmarshaled into a slice of `Post` structs.
- The program prepares an SQL statement for inserting data into the MySQL database.
- The program loops through the `posts` slice and executes the prepared statement for each post, inserting the data into the database.
- Error handling is included for HTTP requests, JSON unmarshalling, and MySQL operations.

## Contributing

Contributions are welcome! If you find any bugs or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.