package pkg

import (
	"testing"
)

func TestSignUpInput(t *testing.T) {
	// valid input
	val := isSignUpInputValid("A", "B", "abc", "a@b.com")
	if !val {
		t.Errorf("should be a valid input")
	}

	// invalid input
	val = isSignUpInputValid("A", "B", "abc", "a@b.com")
	if val {
		t.Errorf("shouldn't be a valid input")
	}

	// invalid input
	val = isSignUpInputValid("A", "B", "abc", "a@b.com")
	if val {
		t.Errorf("shouldn't be a valid input")
	}

	// invalid email
	val = isValidEmail("")
	if val {
		t.Errorf("shouldn't be a valid email")
	}

	// invalid email
	val = isValidEmail("@.com")
	if val {
		t.Errorf("shouldn't be a valid email")
	}

	// invalid email
	val = isValidEmail("a@.com")
	if val {
		t.Errorf("shouldn't be a valid email")
	}

	// invalid email
	val = isValidEmail("@b.com")
	if val {
		t.Errorf("shouldn't be a valid email")
	}

	// invalid email
	val = isValidEmail("@b.")
	if val {
		t.Errorf("shouldn't be a valid email")
	}

	// invalid email
	val = isValidEmail("a@b.")
	if val {
		t.Errorf("shouldn't be a valid email")
	}
}
