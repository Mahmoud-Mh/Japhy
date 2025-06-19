# Japhy Backend

Pet breeds management API built with Go and MySQL.

## Setup

1. Clone the repository
2. Run with Docker:
   ```bash
   docker compose up -d
   ```
3. API available at: http://localhost:50010

## Endpoints

- `GET /health` - Health check
- `GET /breeds` - Get all breeds
- `GET /breeds/{id}` - Get breed by ID
- `POST /breeds` - Create breed
- `PUT /breeds/{id}` - Update breed
- `DELETE /breeds/{id}` - Delete breed

## Search

Filter breeds by species and weight:
```
GET /breeds?species=dog&weight_min=20&weight_max=40
```

## Tech Stack

- Go
- MySQL
- Docker
