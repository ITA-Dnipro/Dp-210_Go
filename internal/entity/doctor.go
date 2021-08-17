package entity

// Doctor struct (Model) and request
type Doctor struct {
	ID         string `json:"id" validate:"omitempty"`
	UserID     string `json:"user_id" validate:"omitempty"`
	FirstName  string `json:"name" validate:"required"`
	LastName   string `json:"last_name" validate:"omitempty"`
	Speciality string `json:"speciality" validate:"required"`
	ScheduleId string `json:"schedule_id" validate:"omitempty"`
}
