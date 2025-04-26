package jwt

var (
	SubjectKey  = ClaimsKey("sub")
	AudienceKey = ClaimsKey("aud")
	IssuedAtKey = ClaimsKey("iat")
	ExpireAtKey = ClaimsKey("exp")
	RolesKey    = ClaimsKey("roles")
)

type ClaimsKey string

type ClaimParsingFunc = func()
