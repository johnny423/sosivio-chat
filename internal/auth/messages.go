package auth

type TokenMsg struct {
	Token string `json:"token"`
}

type ValidateResult struct {
	UserName string `json:"user_name"`
}
