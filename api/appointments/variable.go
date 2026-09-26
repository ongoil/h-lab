package appointments

import (
	"time"

	uuid "github.com/google/uuid"
)

const (
	StatusBooked    = "BOOKED"
	StatusCancelled = "CANCELLED"
	StatusCompleted = "COMPLETED"
)

type appointmentRequest struct {
	PatientID         uuid.UUID `json:"patient_id"`
	DoctorID          uuid.UUID `json:"doctor_id"`
	DepartmentID      uuid.UUID `json:"department_id"`
	AppointmentTypeID uuid.UUID `json:"appointment_type_id"`

	AppointmentDate time.Time `json:"appointment_date"`
	StartTime       string    `json:"start_time"`
	Reason          string    `json:"reason"`
	CreatedBy       string    `json:"created_by"`
}
