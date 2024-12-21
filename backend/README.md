# E-Commerce Backend with PocketBase

This is the backend service for the e-commerce application built with PocketBase.

## Prerequisites

- Go 1.19 or later
- PocketBase

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

The server will start at `http://localhost:8090` by default.

## Features

- Product management
- Category management
- User authentication
- Order processing
- Admin dashboard (available at `http://localhost:8090/_/`) 