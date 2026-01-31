package session

type RegisterResponse struct {
	UID      int64  `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`

}

type LoginOutput struct {
	Token string `json:"token"`
	UserName string `json:"username"`

}
