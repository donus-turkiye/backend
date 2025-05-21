package domain

type Appointment struct {
	AppointmentId     int                `json:"appointment_id"`
	Note              string             `json:"note"`
	Photo             []byte             `json:"photo"`
	IsApproved        *bool              `json:"is_approved"` // Nullable boolean
	IsCompleted       bool               `json:"is_completed"`
	UserId            int                `json:"user_id,omitempty"`
	AppointmentWastes []AppointmentWaste `json:"appointment_wastes,omitempty"`
	Calendars         []Calendar         `json:"calendars,omitempty"` // Calendar slots withdrawn via N:M relationship
}

type AppointmentWaste struct {
	AppointmentWasteId int `json:"appointment_waste_id"`
	CategoryId         int `json:"category_id"`
	Amount             int `json:"amount"`
	AppointmentId      int `json:"appointment_id"`
}
