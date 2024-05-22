package authclient

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/mgrigoriev/goauth/authclient/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticatePositive(t *testing.T) {
	authURL := "https://example.com/auth"
	responseBody := `{"id": 12345, "name": "John Doe", "email": "john.doe@example.com"}`
	body := io.NopCloser(bytes.NewBufferString(responseBody))

	mockHTTPClient := new(mocks.HTTPClient)
	mockHTTPClient.On("Post", authURL, "application/json", mock.Anything).
		Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

	client := New(Config{AuthURL: authURL}, mockHTTPClient)
	userID, err := client.Authenticate("sample_token")

	assert.NoError(t, err)
	assert.Equal(t, 12345, userID)

	mockHTTPClient.AssertExpectations(t)
}

func TestAuthenticateInvalidToken(t *testing.T) {
	authURL := "https://example.com/auth"
	responseBody := `{"error": "invalid token"}`
	body := io.NopCloser(bytes.NewBufferString(responseBody))

	mockHTTPClient := new(mocks.HTTPClient)
	mockHTTPClient.On("Post", authURL, "application/json", mock.Anything).
		Return(&http.Response{StatusCode: http.StatusUnauthorized, Body: body}, nil)

	client := New(Config{AuthURL: authURL}, mockHTTPClient)
	userID, err := client.Authenticate("invalid_token")

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "authentication failed")

	mockHTTPClient.AssertExpectations(t)
}
