package schedules

import uuid "github.com/google/uuid"

type doctorScheduleRequest struct {
	DoctorID     uuid.UUID `json:"doctor_id"`
	DepartmentID uuid.UUID `json:"department_id"`
	DayOfWeek    int       `json:"day_of_week"`           // วันที่ทำงาน จ-ศ
	StartTime    string    `json:"start_time"`            // เริ่มงาน
	EndTime      string    `json:"end_time"`              // เลิกงาน
	BreakStart   string    `json:"break_start,omitempty"` // เริ่มพัก
	BreakEnd     string    `json:"break_end,omitempty"`   // เลิกพัก
	IsAvailable  bool      `json:"is_available"`          // ตารางนี้เปิดให้จองไหม
}
