package service

import (
	"crypto/rand"
	"math/big"

	"github.com/MagnunAVF/random-pass/internal/domain"
)

type PasswordService struct {
	repo domain.PasswordRepository
}

func NewPasswordService(repo domain.PasswordRepository) *PasswordService {
	return &PasswordService{repo: repo}
}

func (s *PasswordService) GenerateAndSave(userID string) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	length := 16
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	pwd := string(result)
	err := s.repo.Save(userID, pwd)
	return pwd, err
}
