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
	responseBody := `{"id": 1, "name": "John", "email": "john.doe@example.com"}`
	body := io.NopCloser(bytes.NewBufferString(responseBody))

	mockHTTPClient := new(mocks.HTTPClient)
	mockHTTPClient.On("Post", authURL, "application/json", mock.Anything).
		Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

	client := New(Config{AuthURL: authURL}, mockHTTPClient)
	user, err := client.Authenticate("sample_token")

	assert.NoError(t, err)
	assert.Equal(t, CurrentUser{
		ID:    1,
		Name:  "John",
		Email: "john.doe@example.com",
	}, *user)

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
	user, err := client.Authenticate("invalid_token")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "authentication failed")

	mockHTTPClient.AssertExpectations(t)
}
