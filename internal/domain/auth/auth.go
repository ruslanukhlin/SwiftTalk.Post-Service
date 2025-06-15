package auth

type VerifyTokenOutput struct {
	UserUUID string
}

type AuthRepository interface {
	VerifyToken(accessToken string) (*VerifyTokenOutput, error)
}
