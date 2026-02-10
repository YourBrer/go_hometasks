package verify

type SendEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
