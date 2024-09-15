package users

import (
	"net/http"
	"site/internal/validation"
)

type SignUpForm struct {
	Name     string                `json:"name"`
	Email    string                `json:"email"`
	Password string                `json:"password"`
	V        validation.Validation `json:"-"`
}

type LoginForm struct {
	Email    string                `json:"email"`
	Password string                `json:"password"`
	V        validation.Validation `json:"-"`
}

func (form *SignUpForm) Decode(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	form.Name = r.PostForm.Get("name")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	return nil
}

func (form *SignUpForm) Validate() {
	form.V.Check(validation.NotBlank(form.Name), "name", "name is required")
	form.V.Check(validation.NotBlank(form.Email), "email", "email is required")
	form.V.Check(validation.NotBlank(form.Password), "password", "password is required")
	form.V.Check(validation.Matches(form.Email, validation.EmailRX), "email", "email is invalid")
	form.V.Check(validation.MinChars(form.Password, 8), "password", "password must be at least 8 charaters long")
}

func (form *LoginForm) Validate() {
	form.V.Check(validation.NotBlank(form.Email), "email", "email is required")
	form.V.Check(validation.NotBlank(form.Password), "password", "password is required")
	form.V.Check(validation.Matches(form.Email, validation.EmailRX), "email", "email is invalid")
	form.V.Check(validation.MinChars(form.Password, 8), "password", "password must be at least 8 charaters long")
}
