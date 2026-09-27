package models

import (
	"time"

	"github.com/google/uuid"
)

// Appointment ข้อมูลการนัดหมายผู้ป่วย
type Appointment struct {
	Model
	// รหัสที่ใช้เชื่อมกับข้อมูลอื่น
	PatientID    uuid.UUID `json:"patient_id" gorm:"not null"`    // รหัสผู้ป่วย
	DoctorID     uuid.UUID `json:"doctor_id" gorm:"not null"`     // รหัสแพทย์
	DepartmentID uuid.UUID `json:"department_id" gorm:"not null"` // รหัสแผนก
	// วันและเวลาของการนัดหมาย
	AppointmentDate time.Time `json:"appointment_date" gorm:"type:date;not null"` // วันที่นัดหมาย
	StartTime       string    `json:"start_time"`                                 // เวลาเริ่มนัดหมาย
	EndTime         string    `json:"end_time"`                                   // เวลาสิ้นสุดนัดหมาย
	AppointmentType string    `json:"appointment_type"`                           // ประเภทการนัดหมาย
	// ข้อมูลการนัดหมาย
	Status    string `json:"status" gorm:"type:varchar(20);not null;index"` // BOOKED, CANCELLED, COMPLETED
	Reason    string `json:"reason,omitempty" gorm:"type:text"`             // เหตุผลหรือรายละเอียดการนัดหมาย
	CreatedBy string `json:"created_by" gorm:"not null"`                    // ผู้ที่สร้างรายการนัดหมาย
	// ความสัมพันธ์กับข้อมูลอื่น
	Patient    Patient    `json:"patient,omitempty" gorm:"foreignKey:PatientID"`       // ผู้ป่วย
	Doctor     Doctor     `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`         // แพทย์
	Department Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"` // แผนก
}
