package models

// Doctor ข้อมูลแพทย์ที่สามารถรับนัดหมายได้
type Doctor struct {
	Model
	TitleName string `json:"title_name" gorm:"not null"` // คำนำหน้าชื่อ
	FirstName string `json:"first_name" gorm:"not null"` // ชื่อ
	LastName  string `json:"last_name" gorm:"not null"`  // นามสกุล
	// ตารางเวลาทำงานของแพทย์
	Schedules []DoctorSchedule `json:"schedules,omitempty" gorm:"foreignKey:DoctorID"`
	// รายการนัดหมายของแพทย์
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:DoctorID"`
}
