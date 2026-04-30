package contextkeys

//create context keys to avoid name collison with service userID/email in request

type contextKey string

const (
	UserID contextKey = "userID"
	Email  contextKey = "email"
)
