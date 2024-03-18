package test

import (
	"bytes"
	"encoding/json"
	"filmoteka/internal/models"
	"filmoteka/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp_Success(t *testing.T) {
	services.IsTesting = true
	req, err := http.NewRequest("GET", "/sign-up", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	services.SignUp(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

func TestSignUp_AlreadyAuthorized(t *testing.T) {
	services.IsTesting = true
	req, err := http.NewRequest("GET", "/sign-up", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Infinity jwt token
	req.Header.Set("Authorization",
		"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjpbInVzZXIiLCJhZG1pbiJdLCJzdWIiOiIxIn0._oXiRDQ0Z82_HDyjwHD_fwH7vPek6fEEmlFWMZvr3Fg")

	rr := httptest.NewRecorder()

	services.SignUp(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestRegistration_Success(t *testing.T) {
	services.IsTesting = true
	user := models.User{
		FirstName: "FirstNameTest",
		LastName:  "LastNameTest",
		Email:     "EmailTest",
		Password:  "PasswordTest",
	}

	reqBody, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	services.Registration(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}
