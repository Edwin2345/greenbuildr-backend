# Greenbuildr Backend

A marketplace for construction sites to sell unused building materials.

## Services

| Service | Port | Purpose |
|---|---|---|
| `graphql-gateway` | 4000 | Single client entry point, routes to services |
| `auth-service` | 4001 | Register, login, issue JWTs |
| `listing-service` | 4002 | Create and browse material listings |

## Architecture

### Hexagonal (Ports & Adapters)
Each microservice (auth, listing) follows hexagonal architecture:

```
adapters/    ← concrete implementations (MySQL, JWT, bcrypt, SendGrid)
ports/       ← interfaces the domain depends on
domain/
  entities/  ← domain models
  errors/    ← domain errors
  services/  ← business logic (one file per action, single Execute() method)
graph/       ← inbound GraphQL adapter (resolvers call domain services)
main.go      ← wires everything together
```

The gateway is flat (no hexagonal) — it is pure infrastructure with no business logic:

```
graphql-gateway/
  middleware/  ← validates client JWT, mints internal JWT
  proxy/       ← forwards requests to downstream services
```

### Authentication Flow
1. Client sends request with a client-facing JWT to the gateway
2. Gateway validates it using `CLIENT_SECRET`, extracts userID + email
3. Gateway mints a short-lived internal JWT (30s) signed with `INTERNAL_SECRET`
4. Internal JWT is forwarded as `X-Internal-Token` header
5. Each downstream service validates it using `shared/middleware.InternalAuth()`
6. userID and email are injected into request context via `shared/contextkeys`

### Shared Module
`shared/` is infrastructure only. Do not add domain ports or business types here.

| Package | Contents |
|---|---|
| `shared/middleware` | `InternalAuth()` — internal JWT validation middleware |
| `shared/contextkeys` | Typed context keys: `UserID`, `Email` |
| `shared/domainerrors` | `DomainError` struct + methods — services alias this type and define their own error vars |
| `shared/gqlerrors` | `ToGQLError()` + `CodedError` interface — translates `DomainError` into GraphQL extensions |

In local dev, services reference shared via `replace` in `go.mod`. In production it will be a private Go module.

## Key Decisions

- **`services/` not `usecases/`** — same discipline (one file per action), preferred naming
- **Email/messaging** — port defined independently in each service that needs it (auth, listing). Not in `shared/` — services must stay decoupled
- **Payments** — port lives inside listing-service only, it is the only service that needs it
- **`shared/`** stays infrastructure-only — cross-cutting types (`DomainError`, `CodedError`), middleware, context keys. Never service-specific business logic or entity types.
- **Internal JWTs expire in 30 seconds** — short enough to be useless if intercepted on the private network
- **Database access — sqlx** (`github.com/jmoiron/sqlx`) chosen over GORM/Bob for auth-service (and by convention all services). Domain entities are plain structs with no ORM tags; DB adapter layer uses separate row structs with `db:""` tags and maps to domain entities.

## GraphQL Error Shape

Domain errors are translated in `graph/gql_errors.go` via `toGQLError()`. Every resolver must call `toGQLError(err)` before returning an error to the client.

Response shape:
```json
{
  "errors": [{ "message": "Email already registered", "extensions": { "code": "EMAIL_ALREADY_EXISTS" } }]
}
```

- `message` — human-readable, from `DomainError.Message`
- `extensions.code` — machine-readable, from `DomainError.Code`
- Logs still use `DomainError.Error()` which prints `[CODE] message` — intentional for debugging

## Email Templates

Templates live in `auth-service/adapters/email/templates/` as plain `.html` files and are loaded from disk at runtime by `ResendSender`. In Docker, the directory is volume-mounted so templates can be edited without rebuilding.

| Variable | Default (local `go run`) |
|---|---|
| `TEMPLATES_DIR` | `./adapters/email/templates` |

## Environment Variables

| Variable | Used by | Purpose |
|---|---|---|
| `JWT_SECRET` | auth-service | Signs client-facing JWTs |
| `CLIENT_SECRET` | gateway | Verifies incoming client JWTs |
| `INTERNAL_SECRET` | gateway + all services | Signs and verifies internal JWTs |
