package domain

type SessionKey string

const (
	UserDataKey SessionKey = "user_data"
)

type UserData struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

// TODO: Add more session data as needed
