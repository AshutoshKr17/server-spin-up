package models

import (
	"reflect"
	"testing"
)

func TestUserTableName(t *testing.T) {
	var u User
	expected := "users"
	if u.TableName() != expected {
		t.Errorf("expected %s, got %s", expected, u.TableName())
	}
}

func TestUserProviderTableName(t *testing.T) {
	var up UserProvider
	expected := "user_providers"
	if up.TableName() != expected {
		t.Errorf("expected %s, got %s", expected, up.TableName())
	}
}

func TestUserFields(t *testing.T) {
	user := User{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		APIKey:   "secret-key",
		IsActive: true,
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email=test@example.com, got %s", user.Email)
	}
	if user.APIKey != "secret-key" {
		t.Errorf("expected APIKey=secret-key, got %s", user.APIKey)
	}
	if !user.IsActive {
		t.Errorf("expected IsActive=true, got %v", user.IsActive)
	}
}

func TestUserProviderFields(t *testing.T) {
	config := JSONMap{"region": "us-east-1", "quota": 5}

	up := UserProvider{
		ID:        10,
		UserID:    1,
		Provider:  "AWS",
		APIKey:    "aws-secret",
		IsEnabled: true,
		Config:    config,
	}

	if up.Provider != "AWS" {
		t.Errorf("expected Provider=AWS, got %s", up.Provider)
	}
	if !up.IsEnabled {
		t.Errorf("expected IsEnabled=true, got %v", up.IsEnabled)
	}
	if !reflect.DeepEqual(up.Config, config) {
		t.Errorf("expected Config=%v, got %v", config, up.Config)
	}
}
