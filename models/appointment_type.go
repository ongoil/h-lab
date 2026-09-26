package models

import (
	"time"

	"github.com/google/uuid"
)

// Doctor คือข้อมูลแพทย์ที่สามารถรับนัดหมายได้
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

// Department คือข้อมูลแผนกของโรงพยาบาล
type Department struct {
	Model
	ENGName string `json:"eng_name" gorm:"not null"` // ชื่อ eng
	THName  string `json:"th_name" gorm:"not null"`  // ชื่อ TH

	// ตารางเวลาของแพทย์ในแผนก
	Schedules []DoctorSchedule `json:"schedules,omitempty" gorm:"foreignKey:DepartmentID"`

	// รายการนัดหมายของแผนก
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:DepartmentID"`
}

// DoctorSchedule คือกำหนดการทำงานของแพทย์
type DoctorSchedule struct {
	Model

	DoctorID     uuid.UUID `json:"doctor_id" gorm:"not null"`     // รหัสแพทย์
	DepartmentID uuid.UUID `json:"department_id" gorm:"not null"` // รหัสแผนก

	DayOfWeek   int    `json:"day_of_week" gorm:"not null"`               // 0 = อาทิตย์, 1 = จันทร์, ..., 6 = เสาร์
	StartTime   string `json:"start_time" gorm:"type:TIME;not null"`      // เวลาเริ่มทำงาน
	EndTime     string `json:"end_time" gorm:"type:TIME;not null"`        // เวลาสิ้นสุดการทำงาน
	BreakStart  string `json:"break_start,omitempty" gorm:"type:TIME"`    // เวลาเริ่มพัก
	BreakEnd    string `json:"break_end,omitempty" gorm:"type:TIME"`      // เวลาสิ้นสุดพัก
	IsAvailable bool   `json:"is_available" gorm:"not null;default:true"` // เปิดให้จองหรือไม่

	// แพทย์
	Doctor Doctor `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`

	// แผนก
	Department Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

// Patient คือข้อมูลผู้ป่วย
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

// AppointmentType คือประเภทของการนัดหมาย
type AppointmentType struct {
	Model

	Name            string `json:"name" gorm:"not null"`                   // ชื่อประเภท เช่น New Patient, Follow-up
	DurationMinutes int    `json:"duration_minutes" gorm:"not null"`       // ระยะเวลานัดหมายเป็นนาที
	IsActive        bool   `json:"is_active" gorm:"not null;default:true"` // เปิดใช้งานประเภทนี้หรือไม่

	// รายการนัดหมายที่ใช้ประเภทนี้
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:AppointmentTypeID"`
}

// Appointment คือข้อมูลการนัดหมายผู้ป่วย
type Appointment struct {
	Model

	// รหัสที่ใช้เชื่อมกับข้อมูลอื่น
	PatientID         uuid.UUID `json:"patient_id" gorm:"not null"`          // รหัสผู้ป่วย
	DoctorID          uuid.UUID `json:"doctor_id" gorm:"not null"`           // รหัสแพทย์
	DepartmentID      uuid.UUID `json:"department_id" gorm:"not null"`       // รหัสแผนก
	AppointmentTypeID uuid.UUID `json:"appointment_type_id" gorm:"not null"` // รหัสประเภทการนัดหมาย

	// วันและเวลาของการนัดหมาย
	AppointmentDate time.Time `json:"appointment_date" gorm:"type:date;not null"` // วันที่นัดหมาย
	StartTime       string    `json:"start_time" gorm:"type:time;not null"`       // เวลาเริ่มนัดหมาย
	EndTime         string    `json:"end_time" gorm:"type:time;not null"`         // เวลาสิ้นสุดนัดหมาย

	// ข้อมูลการนัดหมาย
	Status string `json:"status" gorm:"type:varchar(20);not null;index"` // BOOKED, CANCELLED, COMPLETED
	Reason string `json:"reason,omitempty" gorm:"type:text"`             // เหตุผลหรือรายละเอียดการนัดหมาย

	CreatedBy string `json:"created_by" gorm:"not null"` // ผู้ที่สร้างรายการนัดหมาย

	// ข้อมูลการยกเลิก
	CancelledBy        *string    `json:"cancelled_by,omitempty"`                         // ผู้ที่ยกเลิก
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`                         // วันที่และเวลาที่ยกเลิก
	CancellationReason string     `json:"cancellation_reason,omitempty" gorm:"type:text"` // เหตุผลที่ยกเลิก

	// ความสัมพันธ์กับข้อมูลอื่น
	Patient         Patient         `json:"patient,omitempty" gorm:"foreignKey:PatientID"`                  // ผู้ป่วย
	Doctor          Doctor          `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`                    // แพทย์
	Department      Department      `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`            // แผนก
	AppointmentType AppointmentType `json:"appointment_type,omitempty" gorm:"foreignKey:AppointmentTypeID"` // ประเภทการนัดหมาย
}
