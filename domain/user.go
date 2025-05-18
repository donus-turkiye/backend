package domain

type User struct {
	Id                int    `json:"id"`
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	RoleId            int    `json:"role_id"`
	TelNumber         string `json:"tel_number"`
	Address           string `json:"address"`
	Coordinate        string `json:"coordinate"`
	Wallet            int    `json:"wallet"`
	TotalRecycleCount int    `json:"total_recycle_count"`
}
