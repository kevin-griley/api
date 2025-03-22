# API Starter Template

A simple yet comprehensive API starter template built with Go, incorporating React Native for the client-side. Features include PostgreSQL database integration, Swagger documentation, and a containerized development environment.

## Getting Started

### Prerequisites

- Golang (v1.23 or higher recommended)
- Node.js and npm
- Make utility
- Postgres database

### Running the Project

To run the project, simply use:

```bash
# Configure environment variables
# There is an `.env.example` in the root directory you can use for reference
cp .env.example .env

# Push the schema to the database
make up

# Generate swag init & build & run app
make run

# Run React Native Client
make native
```

## Documentation

- **API Docs:** Once the backend server is running, the API documentation is available at the `/docs` endpoint (e.g., [http://localhost:3000/docs](http://localhost:3000/docs)).

- **Frontend Documentation:** Details for the React Native client can be found within the `/native` directory.

## Project Structure

```
api/
├── cmd
├── data
├── db
├── docs
├── handlers
├── native
├── types
├── Makefile
├── README.md
├── sst.config.ts
└── ...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
