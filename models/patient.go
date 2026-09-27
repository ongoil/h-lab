package models

import "time"

// Patient ข้อมูลผู้ป่วย
type Patient struct {
	Model
	PatientNo   string    `json:"patient_no" gorm:"uniqueIndex;not null"` // เลขประจำตัวผู้ป่วย (HN)
	TitleName   string    `json:"title_name" gorm:"not null"`             // คำนำหน้าชื่อ
	FirstName   string    `json:"first_name" gorm:"not null"`             // ชื่อ
	LastName    string    `json:"last_name" gorm:"not null"`              // นามสกุล
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date"`         // วันเกิดผู้ป่วย
	// รายการนัดหมายของผู้ป่วย
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:PatientID"`
}
