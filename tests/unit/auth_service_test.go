package unit

import (
	"testing"
	"time"

	"fileuploader/internal/config"
	"fileuploader/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService(t *testing.T) {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = "test-secret-key"
	cfg.Auth.TokenExpiration = 24 * time.Hour

	authService := services.NewAuthService(cfg)

	// Test that the service was created and can generate tokens
	assert.NotNil(t, authService)

	// Test functionality by generating and validating a token
	token, err := authService.GenerateToken("test-user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := authService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "test-user", claims.UserID)
}

func TestAuthService_GenerateToken(t *testing.T) {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = "test-secret-key"
	cfg.Auth.TokenExpiration = 24 * time.Hour
	authService := services.NewAuthService(cfg)

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  "user123",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: false, // JWT allows empty claims
		},
		{
			name:    "long user ID",
			userID:  "very-long-user-id-with-special-characters-123!@#$%^&*()",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.GenerateToken(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify the token can be parsed and contains correct claims
				claims, err := authService.ValidateToken(token)
				require.NoError(t, err)
				require.NotNil(t, claims)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.True(t, claims.ExpiresAt.After(time.Now()))
				assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = "test-secret-key"
	cfg.Auth.TokenExpiration = 24 * time.Hour
	authService := services.NewAuthService(cfg)

	// Generate a valid token for testing
	validUserID := "user123"
	validToken, err := authService.GenerateToken(validUserID)
	require.NoError(t, err)

	// Create an expired token by using a different service with very short expiration
	shortCfg := &config.Config{}
	shortCfg.Auth.JWTSecret = "test-secret-key"
	shortCfg.Auth.TokenExpiration = -1 * time.Hour // Negative duration to create expired token
	shortAuthService := services.NewAuthService(shortCfg)

	expiredTokenString, err := shortAuthService.GenerateToken("expiredUser")
	require.NoError(t, err)

	// Create a token with wrong signature for testing
	wrongSignatureToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoid3JvbmdNZXRob2QiLCJleHAiOjk5OTk5OTk5OTl9.invalid-signature"

	tests := []struct {
		name       string
		token      string
		wantErr    bool
		wantClaims *services.Claims
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
			wantClaims: &services.Claims{
				UserID: validUserID,
			},
		},
		{
			name:       "expired token",
			token:      expiredTokenString,
			wantErr:    true,
			wantClaims: nil,
		},
		{
			name:       "malformed token",
			token:      "invalid.token.string",
			wantErr:    true,
			wantClaims: nil,
		},
		{
			name:       "empty token",
			token:      "",
			wantErr:    true,
			wantClaims: nil,
		},
		{
			name:       "token with wrong signature",
			token:      wrongSignatureToken,
			wantErr:    true,
			wantClaims: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tt.wantClaims != nil {
					assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
					assert.True(t, claims.ExpiresAt.After(time.Now()))
				}
			}
		})
	}
}

func TestAuthService_TokenRoundTrip(t *testing.T) {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = "test-secret-key"
	cfg.Auth.TokenExpiration = 1 * time.Hour
	authService := services.NewAuthService(cfg)

	userID := "roundtrip-user"

	// Generate token
	token, err := authService.GenerateToken(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Validate token
	claims, err := authService.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims
	assert.Equal(t, userID, claims.UserID)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
}

func TestAuthService_DifferentSecrets(t *testing.T) {
	cfg1 := &config.Config{}
	cfg1.Auth.JWTSecret = "secret1"
	cfg1.Auth.TokenExpiration = 1 * time.Hour

	cfg2 := &config.Config{}
	cfg2.Auth.JWTSecret = "secret2"
	cfg2.Auth.TokenExpiration = 1 * time.Hour

	authService1 := services.NewAuthService(cfg1)
	authService2 := services.NewAuthService(cfg2)

	userID := "test-user"

	// Generate token with service1
	token, err := authService1.GenerateToken(userID)
	require.NoError(t, err)

	// Try to validate with service2 (different secret)
	claims, err := authService2.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Validate with correct service
	claims, err = authService1.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}
