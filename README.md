# Greenbuildr Backend

A microservices-based backend for construction sites to list unused materials for sale. Built with Go, GraphQL, and Docker.

## Project Structure

```
greenbuildr-backend/
├── auth-service/          # Authentication and user management service
├── listing-service/       # Material listings CRUD service
├── graphql-gateway/       # GraphQL API Gateway (Apollo Federation)
├── docker-compose.yml     # Local development environment
└── README.md             # This file
```

## Services

### Auth Service (Port 4001)
Handles user authentication, registration, and JWT token generation.

**Technology:** Go + GraphQL (gqlgen)
**Database:** MySQL (auth_db)

**Key Features:**
- User registration with email validation
- Login with password hashing
- JWT token generation
- User profile queries

### Listing Service (Port 4002)
Manages construction material listings with geolocation support.

**Technology:** Go + GraphQL (gqlgen)
**Database:** MySQL (listing_db)

**Key Features:**
- CRUD operations for material listings
- Fields: title, description, quantity, price, latitude, longitude
- Location-based queries
- Indexed for performance

### GraphQL Gateway (Port 4000)
Unified GraphQL endpoint using Apollo Federation to combine Auth and Listing service schemas.

**Technology:** Go + Apollo Gateway
**Port:** 4000

**Key Features:**
- Single entry point for frontend
- Schema federation
- Request routing to subgraphs

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- MySQL 8.0+ (if running without Docker)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repo-url>
cd greenbuildr-backend
```

2. Start the services:
```bash
docker-compose up --build
```

3. Access the GraphQL playground:
- Gateway: http://localhost:4000/
- Auth Service: http://localhost:4001/
- Listing Service: http://localhost:4002/

### Local Development

1. Install dependencies for each service:
```bash
cd auth-service
go mod download
cd ../listing-service
go mod download
cd ../graphql-gateway
go mod download
```

2. Run gqlgen for code generation (requires gqlgen CLI):
```bash
cd auth-service
go run github.com/99designs/gqlgen generate
cd ../listing-service
go run github.com/99designs/gqlgen generate
```

3. Set up databases:
```bash
mysql < auth-service/init.sql
mysql < listing-service/init.sql
```

4. Start each service in separate terminals:
```bash
# Terminal 1
cd auth-service
go run main.go

# Terminal 2
cd listing-service
go run main.go

# Terminal 3
cd graphql-gateway
go run main.go
```

## API Examples

### GraphQL Query - Get All Materials
```graphql
query {
  materials {
    id
    title
    description
    quantity
    price
    latitude
    longitude
  }
}
```

### GraphQL Mutation - Create Material
```graphql
mutation {
  createMaterial(input: {
    title: "Unused Drywall"
    description: "10 sheets of 5/8 drywall"
    quantity: 10
    price: 25.00
    latitude: 40.7128
    longitude: -74.0060
    userId: "user-123"
  }) {
    id
    title
    price
  }
}
```

### GraphQL Query - Materials by Location
```graphql
query {
  materialsByLocation(
    latitude: 40.7128
    longitude: -74.0060
    radiusKm: 5
  ) {
    id
    title
    distance
  }
}
```

## Environment Variables

### Auth Service
- `PORT`: Service port (default: 4001)
- `DB_HOST`: MySQL host (default: localhost)
- `DB_USER`: MySQL user (default: root)
- `DB_PASSWORD`: MySQL password
- `DB_NAME`: Database name (default: auth_db)

### Listing Service
- `PORT`: Service port (default: 4002)
- `DB_HOST`: MySQL host (default: localhost)
- `DB_USER`: MySQL user (default: root)
- `DB_PASSWORD`: MySQL password
- `DB_NAME`: Database name (default: listing_db)

### GraphQL Gateway
- `PORT`: Service port (default: 4000)
- `AUTH_SERVICE_URL`: Auth service GraphQL endpoint
- `LISTING_SERVICE_URL`: Listing service GraphQL endpoint

## Database Schema

### auth_db.users
```sql
- id: VARCHAR(36) PRIMARY KEY
- email: VARCHAR(255) UNIQUE NOT NULL
- password_hash: VARCHAR(255) NOT NULL
- created_at: TIMESTAMP
```

### listing_db.materials
```sql
- id: VARCHAR(36) PRIMARY KEY
- title: VARCHAR(255) NOT NULL
- description: TEXT
- quantity: INT NOT NULL
- price: DECIMAL(10, 2) NOT NULL
- latitude: DECIMAL(10, 8) NOT NULL
- longitude: DECIMAL(11, 8) NOT NULL
- user_id: VARCHAR(36) NOT NULL
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- Indexes: user_id, location (lat/lng), created_at
```

## Development Notes

### Adding gqlgen to Services

gqlgen is configured to auto-generate GraphQL types and resolver stubs:

```bash
cd auth-service
go run github.com/99designs/gqlgen@latest generate
```

This creates:
- `graph/generated/generated.go` - Generated types
- `schema.resolvers.go` - Resolver method stubs

### Apollo Federation Setup

The gateway uses Apollo Federation to combine schemas from Auth and Listing services. Each service marks its types with `@key` directives:

```graphql
type User @key(fields: "id") {
  id: ID!
  email: String!
}
```

## Testing

TODO: Add integration tests and unit tests.

## Performance Optimization

- Database indexes on frequently queried fields
- Connection pooling for database connections
- GraphQL query complexity analysis (TODO)
- Caching layer (TODO)

## Next Steps

1. Implement resolver methods for each service
2. Add database connection pooling
3. Set up request validation and error handling
4. Add JWT middleware for protected queries
5. Implement location-based query optimization
6. Add comprehensive testing
7. Set up CI/CD pipeline
8. Deploy to Kubernetes or cloud platform

## Contributing

TODO: Add contribution guidelines

## License

TODO: Add license information
