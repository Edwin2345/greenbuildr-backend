# Greenbuildr Backend

A marketplace for construction sites to sell unused building materials.

## Architecture Overview

```
Client
  │
  ▼ port 4000 (only exposed port)
apollo-auth           ← Go reverse proxy
  │  strips X-User-ID/X-User-Email (anti-spoofing)
  │  validates client JWT (HS256, CLIENT_SECRET)
  │  injects X-User-ID + X-User-Email from claims
  ▼ port 4100 (internal)
Apollo Router                     ← federation gateway
  │  propagates X-User-ID + X-User-Email to subgraphs
  ▼ port 4001 (internal)
auth-service                      ← register, login, verify email
listing-service (port 4002)       ← listings (coming soon)
```

**Security model:** Only the auth proxy is exposed to the host. Apollo Router and all subgraphs live on a private Docker network with no external ports. The proxy strips client-sent identity headers before injecting trusted ones — clients cannot spoof `X-User-ID`.

## Services

| Service | Port | Exposed | Purpose |
|---|---|---|---|
| `apollo-auth` | 4000 | Yes | Auth proxy |
| Apollo Router | 4100 | No | Federation gateway |
| `auth-service` | 4001 | No | Register, login, issue JWTs |
| `listing-service` | 4002 | No | Listings (coming soon) |

## Quick Start

```bash
cp .env.example .env   # fill in your secrets
docker compose up --build
```

GraphQL entry point: `http://localhost:4000/graphql`

## Development

### Regenerating GraphQL code (after schema changes)

```bash
cd auth-service
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate
```

### Recomposing the supergraph (after subgraph schema changes)

```bash
~/.rover/bin/rover supergraph compose --config supergraph.yaml --elv2-license accept 2>/dev/null > supergraph.graphql
```

### Running tests

```bash
cd auth-service && go test ./...
cd shared && go test ./...
```

### Clearing the database (local dev)

```bash
docker exec auth-db mysql -uroot -p$DB_PASSWORD auth_db -e "DELETE FROM auth_tokens; DELETE FROM users;"
```

## Environment Variables

Copy `.env.example` to `.env`. Docker Compose loads it automatically.

| Variable | Purpose |
|---|---|
| `DB_PASSWORD` | MySQL root password |
| `JWT_SECRET` | Signs client-facing JWTs (auth-service) |
| `CLIENT_SECRET` | Verifies client JWTs (auth proxy) |
| `RESEND_API_KEY` | Resend email API key |

## Project Structure

```
greenbuildr-backend/
├── auth-service/                  # Auth subgraph (hexagonal architecture)
│   ├── adapters/                  # MySQL, bcrypt, JWT, Resend email
│   ├── domain/
│   │   ├── entities/              # User, Token
│   │   ├── errors/                # Service-specific error vars
│   │   └── services/              # register.go, login.go
│   ├── graph/                     # gqlgen resolvers + federation
│   ├── ports/                     # Repository, hasher, email, JWT interfaces
│   └── schema.graphql
├── apollo-auth/       # Client-facing auth proxy
├── shared/                        # Cross-cutting infrastructure
│   ├── contextkeys/               # UserID, Email context keys
│   ├── domainerrors/              # DomainError struct
│   ├── gqlerrors/                 # ToGQLError() translator
│   └── middleware/                # InjectUserContext()
├── bruno/                         # Bruno API collections
├── supergraph.yaml                # Rover composition config
├── supergraph.graphql             # Composed supergraph schema (generated)
├── router.yaml                    # Apollo Router config
└── docker-compose.yml
```
