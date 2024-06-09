package domain

type AuthPurpose int64

const (
	PurposeAccess = AuthPurpose(iota)
	PurposeRefresh
)

type (
	Account struct {
		Id        string `json:"id"`
		Role      Role   `json:"role"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Email     string `json:"email"`
		Verified  bool   `json:"verified"`
	}
)
