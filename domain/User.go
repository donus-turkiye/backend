package domain

type User struct {
	Id                int    `json:"id"`
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	RoleId            int    `json:"role_id"`
	TelNumber         string `json:"tel_number"`
	Address           string `json:"address"`
	Coordinate        string `json:"coordinate"`
	Wallet            int    `json:"wallet"`
	TotalRecycleCount int    `json:"total_recycle_count"`
}

type Category struct {
	CategoryId int    `json:"category_id"`
	WasteType  string `json:"waste_type"`
	UnitType   string `json:"unit_type"`
}

type Appointment struct {
	AppointmentId int     `json:"appointment_id"`
	Note          *string `json:"note,omitempty"`
	Photo         *[]byte `json:"photo,omitempty"`
	DateTime      string  `json:"date_time"`
	IsApproved    bool    `json:"is_approved"`
	IsCompleted   bool    `json:"is_completed"`
	UserId        int     `json:"user_id"`
}

type AppointmentWaste struct {
	AppointmentWasteId int `json:"appointment_waste_id"`
	CategoryId        int `json:"category_id"`
	Amount            int `json:"amount"`
	AppointmentId     int `json:"appointment_id"`
}

type RecyclingPoint struct {
	RecyclingPointId int    `json:"recycling_point_id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	Coordinate       string `json:"coordinate"`
}

type Calendar struct {
	CalendarId  int    `json:"calendar_id"`
	IsAvailable bool   `json:"is_available"`
	Date        string `json:"date"`
	Hour        string `json:"hour"`
}
