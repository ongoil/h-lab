package models

import "github.com/google/uuid"

// DoctorSchedule กำหนดการทำงานของแพทย์
type DoctorSchedule struct {
	Model
	DoctorID     uuid.UUID `json:"doctor_id" gorm:"not null"`                 // รหัสแพทย์
	DepartmentID uuid.UUID `json:"department_id" gorm:"not null"`             // รหัสแผนก
	DayOfWeek    int       `json:"day_of_week" gorm:"not null"`               // 0 = อาทิตย์, 1 = จันทร์, ..., 6 = เสาร์
	StartTime    string    `json:"start_time" gorm:"type:TIME;not null"`      // เวลาเริ่มทำงาน
	EndTime      string    `json:"end_time" gorm:"type:TIME;not null"`        // เวลาสิ้นสุดการทำงาน
	BreakStart   string    `json:"break_start,omitempty" gorm:"type:TIME"`    // เวลาเริ่มพัก
	BreakEnd     string    `json:"break_end,omitempty" gorm:"type:TIME"`      // เวลาสิ้นสุดพัก
	IsAvailable  bool      `json:"is_available" gorm:"not null;default:true"` // เปิดให้จองหรือไม่
	// แพทย์
	Doctor Doctor `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	// แผนก
	Department Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}
