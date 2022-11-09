package domain

type Signature struct {
	ID              string `json:"-" db:"id"`
	PublicKey       string `json:"public_key" binding:"required"`
	Body            string `json:"body" binding:"required"`
	Login           string `json:"login" binding:"required"`
	TimeOfSignature string `json:"time_of_signature" binding:"required"`
}

func NewSignature(publicKey string, body string, login string, time string) *Signature {
	return &Signature{PublicKey: publicKey, Body: body, Login: login, TimeOfSignature: time}
}
