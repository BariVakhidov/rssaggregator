package token

type Verifier interface {
	VerifyToken(token string) (*Payload, error)
}
