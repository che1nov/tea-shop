package requestctx

type Key string

const (
	UserIDKey Key = "user_id"
	EmailKey  Key = "email"
	RoleKey   Key = "role"
)
