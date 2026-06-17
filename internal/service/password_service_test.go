package service

import (
	"errors"
	"testing"
)

type mockRepo struct {
	saveFunc        func(userID string, password string) error
	getLastFiveFunc func(userID string) ([]string, error)
}

func (m *mockRepo) Save(userID string, password string) error {
	if m.saveFunc != nil {
		return m.saveFunc(userID, password)
	}
	return nil
}

func (m *mockRepo) GetLastFive(userID string) ([]string, error) {
	if m.getLastFiveFunc != nil {
		return m.getLastFiveFunc(userID)
	}
	return []string{}, nil
}

func TestNewPasswordService(t *testing.T) {
	repo := &mockRepo{}
	service := NewPasswordService(repo)

	if service == nil {
		t.Fatal("expected service to be created, got nil")
	}

	if service.repo != repo {
		t.Error("expected repo to be set correctly")
	}
}

func TestGenerateAndSave(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		saveErr error
		wantErr bool
	}{
		{
			name:    "successful generation and save",
			userID:  "user123",
			saveErr: nil,
			wantErr: false,
		},
		{
			name:    "save error",
			userID:  "user456",
			saveErr: errors.New("save failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var savedPassword string
			var savedUserID string

			repo := &mockRepo{
				saveFunc: func(userID string, password string) error {
					savedUserID = userID
					savedPassword = password
					return tt.saveErr
				},
			}

			service := NewPasswordService(repo)
			password, err := service.GenerateAndSave(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAndSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(password) != 16 {
					t.Errorf("expected password length 16, got %d", len(password))
				}

				if password != savedPassword {
					t.Error("returned password should match saved password")
				}

				if savedUserID != tt.userID {
					t.Errorf("expected userID %s, got %s", tt.userID, savedUserID)
				}

				charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
				for _, char := range password {
					found := false
					for _, validChar := range charset {
						if char == validChar {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("password contains invalid character: %c", char)
					}
				}
			}
		})
	}
}
