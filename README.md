# API Project

This is a simple API project built with Go.

## Getting Started

### Prerequisites

- Go installed on your system
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
```

## Documentation

API documentation is available at `/docs` endpoint when the server is running.

## Project Structure

```
api/
├── cmd 
├── data
├── db
├── docs
├── handlers
├── types
├── Makefile
├── README.md
└── ...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.