package auth

import(
	"anonymous/commons"
	"anonymous/validator"
)

type registrationPayload struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type loginPayload struct {
	Method   string `json:"method"`
	Username    string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *registrationPayload) Validate() (err map[string]string) {
	err = map[string]string{}
	if validator.IsEmptyString(p.Username) {
		err["username"] = commons.Codes.EmptyField
		return
	}
	if validator.IsEmptyString(p.Password) {
		err["password"] = commons.Codes.EmptyField
		return
	}
	if validator.IsEmptyString(p.Email) {
		err["email"] = commons.Codes.EmptyField
		return
	}

	return nil
}

func (p *loginPayload) Validate() (err map[string]string) {
  err = map[string]string{}
	if validator.IsEmptyString(p.Method) {
		err["method"] = commons.Codes.EmptyField
		return
	}
	if !validator.IsOneOf(p.Method, "username", "email") {
		err["method"] = commons.Codes.InvalidField
		return
	}
  if p.Method == "username"{
    if validator.IsEmptyString(p.Username) {
      err["username"] = commons.Codes.EmptyField
      return
    }
  }
  if p.Method == "email"{
    if validator.IsEmptyString(p.Email) {
      err["email"] = commons.Codes.EmptyField
      return
    }
    if !validator.IsEmail(p.Email) {
      err["email"] = commons.Codes.InvalidField
      return
    }
  }
	if validator.IsEmptyString(p.Password) {
		err["method"] = commons.Codes.EmptyField
		return
	}
	return nil
}

type forgotPasswordPayload struct {
    Email string `json:"email"`
}

func (p *forgotPasswordPayload) Validate() []string {
    var errors []string
    if p.Email == "" {
        errors = append(errors, "Email is required")
    }
    return errors
}

type resetPasswordPayload struct {
    Token           string `json:"token"`
    NewPassword     string `json:"new_password"`
    ConfirmPassword string `json:"confirm_password"`
}

func (p *resetPasswordPayload) Validate() []string {
    var errors []string
    if p.Token == "" {
        errors = append(errors, "Token is required")
    }
    if p.NewPassword == "" {
        errors = append(errors, "New password is required")
    }
    if p.ConfirmPassword == "" {
        errors = append(errors, "Confirm password is required")
    }
    if p.NewPassword != p.ConfirmPassword {
        errors = append(errors, "Passwords do not match")
    }
    if len(p.NewPassword) < 8 {
        errors = append(errors, "Password must be at least 8 characters long")
    }
    return errors
}