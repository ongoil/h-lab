package appointments

import (
	uuid "github.com/google/uuid"
)

const (
	StatusBooked    = "BOOKED"
	StatusCancelled = "CANCELLED"
	StatusCompleted = "COMPLETED"
)

type appointmentRequest struct {
	PatientID       uuid.UUID `json:"patient_id"`
	DoctorID        uuid.UUID `json:"doctor_id"`
	DepartmentID    uuid.UUID `json:"department_id"`
	AppointmentType string    `json:"appointment_type"`
	AppointmentDate string    `json:"appointment_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	Reason          string    `json:"reason"`
	CreatedBy       string    `json:"created_by"`
}
