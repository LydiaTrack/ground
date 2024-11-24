package session

type CreateSessionCommand struct {
	UserID       string `json:"userID"`
	ExpireTime   int64  `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

type DeleteSessionCommand struct {
	UserID string `json:"userID"`
}
