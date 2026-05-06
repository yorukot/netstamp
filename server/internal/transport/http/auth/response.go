package auth

type userResponse struct {
	ID          string  `json:"id" format:"uuid"`
	Email       string  `json:"email" format:"email"`
	DisplayName *string `json:"displayName,omitempty" maxLength:"100"`
}
