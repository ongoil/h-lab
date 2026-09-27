package models

// Department ข้อมูลแผนกของโรงพยาบาล
type Department struct {
	Model
	ENGName string `json:"eng_name" gorm:"not null"` // ชื่อ eng
	THName  string `json:"th_name" gorm:"not null"`  // ชื่อ TH
	// ตารางเวลาของแพทย์ในแผนก
	Schedules []DoctorSchedule `json:"schedules,omitempty" gorm:"foreignKey:DepartmentID"`
	// รายการนัดหมายของแผนก
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:DepartmentID"`
}
