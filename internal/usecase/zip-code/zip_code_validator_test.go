package usecase

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func mockedCallApiError(url string) (*http.Response, error) {
	return &http.Response{Body: http.NoBody}, nil
}

func mockedCallApiOk(url string) (*http.Response, error) {
	closer := io.NopCloser(strings.NewReader("{\"cep\":\"11111111\", \"localidade\": \"Algum Lugar\", \"uf\": \"XR\"}"))
	return &http.Response{Body: closer}, nil
}

func TestGivenSmallZipCode_WhenVerifyZipCode_ThenShouldReceiveAnError(t *testing.T) {
	//Given
	validator := NewZipCodeValidator(mockedCallApiOk, "", []string{"11111111"})

	//When
	_, err := validator.VerifyZipCode("1111111")

	//Then
	assert.ErrorContains(t, err, "invalid zip code")
}

func TestGivenAUnavailableZipCode_WhenVerifyZipCode_ThenShouldReceiveAnError(t *testing.T) {
	//Given
	validator := NewZipCodeValidator(mockedCallApiOk, "", []string{""})

	//When
	_, err := validator.VerifyZipCode("11111111")

	//Then
	assert.ErrorContains(t, err, "unavailable zip code")
}

func TestGivenErrorOnZipCodeAPI_WhenVerifyZipCode_ThenShouldReceiveAnError(t *testing.T) {
	//Given
	validator := NewZipCodeValidator(mockedCallApiError, "", []string{"11111111"})

	//When
	_, err := validator.VerifyZipCode("11111111")

	//Then
	assert.Error(t, err)
}

func TestGivenZipCode_WhenVerifyZipCode_ThenShouldReceiveOk(t *testing.T) {
	//Given
	validator := NewZipCodeValidator(mockedCallApiOk, "", []string{"11111111"})

	//When
	zipCode, err := validator.VerifyZipCode("11111111")

	//Then
	assert.NoError(t, err)
	assert.Equal(t, "11111111", zipCode.Id)
	assert.Equal(t, "Algum Lugar", zipCode.City)
	assert.Equal(t, "XR", zipCode.State)
}
