# Qube Coding Challenge 2016

This Go project provides a system to manage distributors, including adding sub-distributors and checking authorization for specific regions. It includes HTTP API endpoints for managing distributor data and performing location-based checks.

## Features

- **Add Distributor**: Add a new distributor with inclusion and exclusion criteria for locations.
- **Add Sub-Distributor**: Add sub-distributors under existing distributors with inherited permissions.
- **Search Authorization**: Check if a distributor is authorized for a given location (city, state, or country).
- **Data Initialization**: Initialize location data from a CSV file.

## Project Structure

- **`main.go`**: Contains the HTTP server setup and endpoint registration.
- **`controller.go`**: Contains core logic for adding distributors, sub-distributors, and checking authorization.
- **`helper.go`**: Provides helper functions for sending HTTP responses.


## Endpoints

1. **POST /distributor**
   - Adds a new distributor with include and exclude location criteria.

2. **POST /distributor/{name}/sub-distributor**
   - Adds a sub-distributor to an existing distributor.

3. **POST /distributor/{name}/search**
   - Checks if a distributor is authorized for a given location.

## Data Initialization

- The `init()` function loads distributor data and reads location data from a CSV file (`cities.csv`), mapping locations by city, state, and country.

## Example Requests

### Add a Distributor
```json
{
  "name": "Distributor1",
  "include": [
    {"city": "mysuru", "state": "KARNATAKA", "country": "INDIA"}
  ],
  "exclude": [
    {"city": "chennai", "state": "TAMIL NADU", "country": "INDIA"}
  ]
}
```
### Add a Sub-Distributor
```json
{
  "name": "SubDistributor1",
  "include": [
    {"state": "KARNATAKA", "country": "INDIA"}
  ],
  "exclude": [
    {"city": "chennai", "state": "TAMIL NADU", "country": "INDIA"}
  ]
}
```

### Check Authorization
```json
{
  "city": "mysuru"
}
```

## Running the Project

1. Ensure you have Go installed.
2. Run the server:
   ```bash
   go run main.go
   ```
3. Access the endpoints at `http://localhost:8080`.

## Dependencies

- Go standard library (`net/http`, `strings`, `encoding/json`, etc.)
- CSV file (`cities.csv`) for initializing location data.

## Steps To Run

- Step 1
    Initialize the go modules (run `go mod init` command)
- Step 2
    Build the project (run `go build` command)
- Step 3
     Run `go run .`

Voila The program is running on 8080 port.