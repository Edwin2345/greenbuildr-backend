package graph

import "github.com/greenbuildr/auth-service/domain/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RegisterService *services.RegisterService
	LoginService    *services.LoginService
}
