package ports

type TokenGenerator interface {
	GenerateJWT(userID, email string) (string, error)
}
