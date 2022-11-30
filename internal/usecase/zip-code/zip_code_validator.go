package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ZipCode struct {
	Id    string `json:"cep"`
	City  string `json:"localidade"`
	State string `json:"uf"`
}

type ZipCodeValidator struct {
	CallApi         CallApi
	Uri             string
	ZipCodesAllowed []string
}

type CallApi func(url string) (*http.Response, error)

func NewZipCodeValidator(callApi CallApi, uri string, zipCodesAllowed []string) *ZipCodeValidator {
	return &ZipCodeValidator{
		CallApi:         callApi,
		Uri:             uri,
		ZipCodesAllowed: zipCodesAllowed,
	}
}

func (z *ZipCodeValidator) VerifyZipCode(zipCode string) (*ZipCode, error) {
	err := z.isValid(zipCode)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/json", z.Uri, zipCode)
	response, err := z.CallApi(url)
	if err != nil {
		return nil, err
	}

	return parseToZipCode(response.Body)
}

func (z *ZipCodeValidator) isValid(zipCode string) error {
	if len(zipCode) != 8 {
		return errors.New("invalid zip code")
	}

	for _, code := range z.ZipCodesAllowed {
		if code == zipCode {
			return nil
		}
	}
	return errors.New("unavailable zip code")
}

func parseToZipCode(closer io.ReadCloser) (*ZipCode, error) {
	body, err := io.ReadAll(closer)
	if err != nil {
		return nil, err
	}

	var zipCode ZipCode
	err = json.Unmarshal(body, &zipCode)
	if err != nil {
		return nil, err
	}

	return &zipCode, nil
}
