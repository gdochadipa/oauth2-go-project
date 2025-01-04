package repository

type JWTRepository struct {
	secretKey []byte
}

func (jwt *JWTRepository) sign() (*string, error) {

}
