package domain

type PasswordEntry struct {
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
}

type PasswordRepository interface {
	Save(userID string, password string) error
	GetLastFive(userID string) ([]string, error)
}
