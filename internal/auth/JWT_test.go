package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {

	uuid1 := uuid.New()
	uuid2 := uuid.New()
	tokenSecret1 := "asd"
	tokenSecret2 := "fdsdf"
	tokenSecret3 := "fast"
	expiresin := time.Second * 10
	expireFast := time.Millisecond * 10

	token1, _ := MakeJWT(uuid1, tokenSecret1, expiresin)
	token2, _ := MakeJWT(uuid2, tokenSecret2, expiresin)
	token3, _ := MakeJWT(uuid2, tokenSecret3, expireFast)
	time.Sleep(time.Millisecond * 20) // Wait for token to expire

	tests := []struct {
		name        string
		id          uuid.UUID
		token       string
		tokenSecret string
		wantErr     bool
	}{
		{
			name:        "valid token",
			id:          uuid1,
			token:       token1,
			tokenSecret: tokenSecret1,
			wantErr:     false,
		},
		{
			name:        "invalid token",
			id:          uuid.Nil,
			token:       token1,
			tokenSecret: tokenSecret2,
			wantErr:     true,
		},
		{
			name:        "valid token2",
			id:          uuid2,
			token:       token2,
			tokenSecret: tokenSecret2,
			wantErr:     false,
		},
		{
			name:        "invalid token2",
			id:          uuid.Nil,
			token:       token2,
			tokenSecret: tokenSecret1,
			wantErr:     true,
		},
		{
			name:        "expired",
			id:          uuid2,
			token:       token3,
			tokenSecret: tokenSecret3,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := ValidateJWT(tt.token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			//check that returned UUID matches
			if err == nil && uid != tt.id {
				t.Errorf("ValidateJWT() got = %v, want %v", uid, tt.id)
			}

		})
	}

}
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
