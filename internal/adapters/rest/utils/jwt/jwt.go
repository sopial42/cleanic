package jwt

var (
	IDKey       = ClaimsKey("jti")
	SubjectKey  = ClaimsKey("sub")
	IssuedAtKey = ClaimsKey("iat")
	ExpireAtKey = ClaimsKey("exp")
	RolesKey    = ClaimsKey("roles")
)

type ClaimsKey string

type ClaimParsingFunc = func()
