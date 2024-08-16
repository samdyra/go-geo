# Project Setup and Makefile Documentation

This document outlines the setup process and usage of the Makefile for the go-geo project.

## Table of Contents
1. [Environment Setup](#environment-setup)
2. [Database Setup](#database-setup)
3. [Makefile Commands](#makefile-commands)
4. [Usage Example](#usage-example)

## Environment Setup

Before running the project, you need to set up your environment variables:

1. Copy the `.env.example` file to create a new `.env` file:
   ```
   cp .env.example .env
   ```

2. Edit the `.env` file to set your specific configuration:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=go-geo
   DB_PASSWORD=password
   DB_NAME=go-geo
   SERVER_PORT=8080
   ```

   Adjust these values as needed for your development environment.

## Database Setup

After setting up your environment variables, you can set up the PostgreSQL database:

1. Ensure you have Docker installed on your system.

2. Use the provided `start-postgres.sh` script to run a PostgreSQL instance:
   ```
   ./start-postgres.sh
   ```

   This script will use the environment variables from your `.env` file to configure the PostgreSQL container.

3. Alternatively, you can use the Makefile command (after reviewing the Makefile section below):
   ```
   make start-db
   ```

   This command will start a PostgreSQL container with PostGIS, using the environment variables from your `.env` file.

Remember to stop the database when you're done:
```
make stop-db
```

## Makefile Commands

The project includes a Makefile with various commands to simplify development tasks. Here's an overview of the available commands:

### Database Operations
- `make start-db`: Starts a PostgreSQL container with PostGIS.
- `make stop-db`: Stops and removes the PostgreSQL container.
- `make migrate`: Runs database migrations.

### Build and Run
- `make build`: Builds the application.
- `make run`: Builds and runs the application.
- `make clean`: Removes build artifacts.

### Testing and Dependencies
- `make test`: Runs all tests in the project.
- `make deps`: Downloads project dependencies.

### Documentation
- `make docs`: Generates API documentation using Swag.

### Utility Commands
- `make help`: Displays a list of available commands with descriptions.
- `make all`: Cleans and builds the project.
- `make setup`: Runs multiple commands in sequence (deps, install-swag, install-lint, migrate, docs).

## Usage Example

Here's a typical workflow for setting up and running the project:

1. Set up the environment:
   ```
   cp .env.example .env
   # Edit .env as needed
   ```

2. Start the database:
   ```
   ./start-postgres.sh
   # or
   make start-db
   ```

3. Set up the project:
   ```
   make setup
   ```

4. Run the application:
   ```
   make run
   ```

5. When you're done, stop the database:
   ```
   make stop-db
   ```

This workflow will set up your environment, start the database, prepare the project (including generating documentation), and run the application.