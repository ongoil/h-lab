package appointments

import (
	"errors"
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAppointment(c *fiber.Ctx) error {
	var req appointmentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid request body",
			MessageTh: "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}

	// -------------------------
	// Validation
	// -------------------------
	if req.PatientID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Patient ID is required",
			MessageTh: "กรุณาระบุผู้ป่วย",
		})
	}

	if req.DoctorID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Doctor ID is required",
			MessageTh: "กรุณาระบุแพทย์",
		})
	}

	if req.DepartmentID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Department ID is required",
			MessageTh: "กรุณาระบุแผนก",
		})
	}

	if req.AppointmentDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Appointment date is required",
			MessageTh: "กรุณาระบุวันที่นัดหมาย",
		})
	}

	if req.StartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Start time is required",
			MessageTh: "กรุณาระบุเวลาเริ่มนัดหมาย",
		})
	}

	if req.EndTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "End time is required",
			MessageTh: "กรุณาระบุเวลาสิ้นสุดนัดหมาย",
		})
	}

	if req.CreatedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Created by is required",
			MessageTh: "กรุณาระบุผู้สร้างรายการ",
		})
	}

	// -------------------------
	// Parse appointment date
	// -------------------------
	appointmentDate, err := time.Parse(
		"2006-01-02",
		req.AppointmentDate,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid appointment date",
			MessageTh: "รูปแบบวันที่นัดหมายไม่ถูกต้อง",
		})
	}

	// -------------------------
	// Parse appointment time
	// -------------------------
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid start time",
			MessageTh: "รูปแบบเวลาเริ่มนัดหมายไม่ถูกต้อง",
		})
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid end time",
			MessageTh: "รูปแบบเวลาสิ้นสุดนัดหมายไม่ถูกต้อง",
		})
	}

	// เวลาเริ่มต้องน้อยกว่าเวลาสิ้นสุด
	if !startTime.Before(endTime) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid appointment time",
			MessageTh: "เวลาเริ่มต้นต้องน้อยกว่าเวลาสิ้นสุด",
		})
	}

	// -------------------------
	// Combine Date + StartTime
	// -------------------------
	appointmentDateTime, err := time.Parse(
		"2006-01-02 15:04",
		req.AppointmentDate+" "+req.StartTime,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid appointment date or time",
			MessageTh: "รูปแบบวันที่หรือเวลาไม่ถูกต้อง",
		})
	}

	// ห้ามจองย้อนหลัง
	if appointmentDateTime.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Cannot book an appointment in the past",
			MessageTh: "ไม่สามารถจองนัดหมายย้อนหลังได้",
		})
	}

	// -------------------------
	// Day of Week
	// 0 = Sunday
	// 1 = Monday
	// ...
	// 6 = Saturday
	// -------------------------
	dayOfWeek := int(appointmentDate.Weekday())

	// -------------------------
	// Transaction
	// -------------------------
	tx := database.DBConn.Begin()

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to start transaction",
			MessageTh: "ไม่สามารถเริ่ม Transaction ได้",
			Error:     tx.Error.Error(),
		})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// -------------------------
	// Check Patient
	// -------------------------
	var patient models.Patient

	if err := tx.
		Where("id = ?", req.PatientID).
		First(&patient).Error; err != nil {

		tx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Patient not found",
			MessageTh: "ไม่พบข้อมูลผู้ป่วย",
		})
	}

	// -------------------------
	// Check Doctor
	// -------------------------
	var doctor models.Doctor

	if err := tx.
		Where("id = ?", req.DoctorID).
		First(&doctor).Error; err != nil {

		tx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Doctor not found",
			MessageTh: "ไม่พบข้อมูลแพทย์",
		})
	}

	// -------------------------
	// Check Department
	// -------------------------
	var department models.Department

	if err := tx.
		Where("id = ?", req.DepartmentID).
		First(&department).Error; err != nil {

		tx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Department not found",
			MessageTh: "ไม่พบข้อมูลแผนก",
		})
	}

	// -------------------------
	// Check Doctor Schedule
	// -------------------------
	var schedule models.DoctorSchedule

	err = tx.
		Where(
			"doctor_id = ? AND day_of_week = ? AND is_available = ?",
			req.DoctorID,
			dayOfWeek,
			true,
		).
		First(&schedule).Error

	if err != nil {

		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
				Status:    "400",
				Message:   "Doctor is not available on this day",
				MessageTh: "แพทย์ไม่มีตารางทำงานในวันที่เลือก",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to check doctor schedule",
			MessageTh: "ไม่สามารถตรวจสอบตารางแพทย์ได้",
			Error:     err.Error(),
		})
	}

	// -------------------------
	// Parse Doctor Schedule Time
	// -------------------------
	scheduleStart, err := time.Parse(
		"15:04",
		schedule.StartTime[:5],
	)

	if err != nil {
		tx.Rollback()

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Invalid doctor schedule start time",
			MessageTh: "เวลาเริ่มงานของแพทย์ไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	scheduleEnd, err := time.Parse(
		"15:04",
		schedule.EndTime[:5],
	)

	if err != nil {
		tx.Rollback()

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Invalid doctor schedule end time",
			MessageTh: "เวลาเลิกงานของแพทย์ไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	// -------------------------
	// Check Working Hours
	// -------------------------
	if startTime.Before(scheduleStart) ||
		endTime.After(scheduleEnd) {

		tx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Appointment time is outside doctor working hours",
			MessageTh: "เวลาที่เลือกอยู่นอกเวลาทำงานของแพทย์",
			Data: fiber.Map{
				"work_start": schedule.StartTime[:5],
				"work_end":   schedule.EndTime[:5],
			},
		})
	}

	// -------------------------
	// Check Break Time
	// -------------------------
	if schedule.BreakStart != "" &&
		schedule.BreakEnd != "" {

		breakStart, err := time.Parse(
			"15:04",
			schedule.BreakStart[:5],
		)

		if err != nil {
			tx.Rollback()

			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
				Status:    "500",
				Message:   "Invalid break start time",
				MessageTh: "เวลาเริ่มพักไม่ถูกต้อง",
				Error:     err.Error(),
			})
		}

		breakEnd, err := time.Parse(
			"15:04",
			schedule.BreakEnd[:5],
		)

		if err != nil {
			tx.Rollback()

			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
				Status:    "500",
				Message:   "Invalid break end time",
				MessageTh: "เวลาสิ้นสุดพักไม่ถูกต้อง",
				Error:     err.Error(),
			})
		}

		// ตรวจสอบว่าเวลานัดทับกับเวลาพักหรือไม่
		if startTime.Before(breakEnd) &&
			endTime.After(breakStart) {

			tx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
				Status:    "400",
				Message:   "Appointment time is during doctor break",
				MessageTh: "เวลาที่เลือกอยู่ในช่วงพักของแพทย์",
				Data: fiber.Map{
					"break_start": schedule.BreakStart[:5],
					"break_end":   schedule.BreakEnd[:5],
				},
			})
		}
	}

	// -------------------------
	// Check Existing Appointment
	// -------------------------
	var count int64

	err = tx.
		Model(&models.Appointment{}).
		Where(`
			doctor_id = ?
			AND DATE(appointment_date) = ?
			AND status = ?
			AND start_time < ?
			AND end_time > ?
		`,
			req.DoctorID,
			req.AppointmentDate,
			StatusBooked,
			req.EndTime,
			req.StartTime,
		).
		Count(&count).Error

	if err != nil {
		tx.Rollback()

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to check existing appointment",
			MessageTh: "ไม่สามารถตรวจสอบการนัดหมายเดิมได้",
			Error:     err.Error(),
		})
	}

	if count > 0 {
		tx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Appointment time is already booked",
			MessageTh: "วันที่และเวลาที่เลือกมีการจองแล้ว",
			Data: fiber.Map{
				"date":       req.AppointmentDate,
				"start_time": req.StartTime,
				"end_time":   req.EndTime,
			},
		})
	}

	// -------------------------
	// Create Appointment
	// -------------------------
	appointment := models.Appointment{
		PatientID:       patient.Id,
		DoctorID:        doctor.Id,
		DepartmentID:    department.Id,
		AppointmentType: req.AppointmentType,
		AppointmentDate: appointmentDateTime,
		Status:          StatusBooked,
		Reason:          req.Reason,
		CreatedBy:       req.CreatedBy,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
	}

	if err := tx.Create(&appointment).Error; err != nil {
		tx.Rollback()

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to create appointment",
			MessageTh: "ไม่สามารถสร้างการนัดหมายได้",
			Error:     err.Error(),
		})
	}

	// -------------------------
	// Commit
	// -------------------------
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to commit transaction",
			MessageTh: "ไม่สามารถบันทึก Transaction ได้",
			Error:     err.Error(),
		})
	}

	// -------------------------
	// Response
	// -------------------------
	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Status:    "201",
		Message:   "Appointment created successfully",
		MessageTh: "สร้างการนัดหมายสำเร็จ",
		Data:      appointment,
	})
}
func GetAppointment(c *fiber.Ctx) error {
	doctorID := strings.TrimSpace(c.Query("doctor_id"))
	date := strings.TrimSpace(c.Query("date"))
	startTime := strings.TrimSpace(c.Query("start_time"))
	endTime := strings.TrimSpace(c.Query("end_time"))

	// =========================
	// Validate required params
	// =========================
	if doctorID == "" || date == "" || startTime == "" || endTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Missing required parameters",
			MessageTh: "กรุณาระบุแพทย์ วันที่ เวลาเริ่มต้น และเวลาสิ้นสุด",
		})
	}

	// =========================
	// Parse date
	// =========================
	appointmentDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid date format",
			MessageTh: "รูปแบบวันที่ไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	// =========================
	// Parse time
	// =========================
	requestStart, err := time.Parse("15:04", startTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid start time format",
			MessageTh: "รูปแบบเวลาเริ่มต้นไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	requestEnd, err := time.Parse("15:04", endTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid end time format",
			MessageTh: "รูปแบบเวลาสิ้นสุดไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	// เวลาเริ่มต้องน้อยกว่าเวลาสิ้นสุด
	if !requestStart.Before(requestEnd) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid appointment time",
			MessageTh: "เวลาเริ่มต้นต้องน้อยกว่าเวลาสิ้นสุด",
		})
	}

	// =========================
	// หา Day Of Week
	// 0 = Sunday
	// 1 = Monday
	// ...
	// 6 = Saturday
	// =========================
	dayOfWeek := int(appointmentDate.Weekday())

	// =========================
	// ตรวจสอบ Doctor Schedule
	// =========================
	var schedule models.DoctorSchedule
	err = database.DBConn.Where("doctor_id = ? AND day_of_week = ? AND is_available = ?", doctorID, dayOfWeek, true).First(&schedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusOK).JSON(dto.Response{
				Status:    "200",
				Message:   "Doctor is not available on this day",
				MessageTh: "แพทย์ไม่มีตารางทำงานในวันที่เลือก",
				Data: fiber.Map{
					"available":  false,
					"doctor_id":  doctorID,
					"date":       date,
					"start_time": startTime,
					"end_time":   endTime,
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to get doctor schedule",
			MessageTh: "ไม่สามารถตรวจสอบตารางแพทย์ได้",
			Error:     err.Error(),
		})
	}

	// =========================
	// Parse Doctor Schedule Time
	// =========================
	scheduleStart, err := time.Parse("15:04", schedule.StartTime[:5])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Invalid doctor schedule start time",
			MessageTh: "เวลาเริ่มงานของแพทย์ไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	scheduleEnd, err := time.Parse("15:04", schedule.EndTime[:5])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Invalid doctor schedule end time",
			MessageTh: "เวลาเลิกงานของแพทย์ไม่ถูกต้อง",
			Error:     err.Error(),
		})
	}

	// =========================
	// ตรวจสอบว่าเวลาอยู่ในเวลาทำงานหรือไม่
	// =========================
	if requestStart.Before(scheduleStart) || requestEnd.After(scheduleEnd) {
		return c.Status(fiber.StatusOK).JSON(dto.Response{
			Status:    "200",
			Message:   "Appointment time is outside doctor working hours",
			MessageTh: "เวลาที่เลือกอยู่นอกเวลาทำงานของแพทย์",
			Data: fiber.Map{
				"available":  false,
				"doctor_id":  doctorID,
				"date":       date,
				"start_time": startTime,
				"end_time":   endTime,
				"work_start": schedule.StartTime[:5],
				"work_end":   schedule.EndTime[:5],
			},
		})
	}

	// =========================
	// ตรวจสอบช่วงพัก
	// =========================
	if schedule.BreakStart != "" && schedule.BreakEnd != "" {

		breakStart, err := time.Parse("15:04", schedule.BreakStart[:5])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
				Status:    "500",
				Message:   "Invalid break start time",
				MessageTh: "เวลาเริ่มพักไม่ถูกต้อง",
				Error:     err.Error(),
			})
		}

		breakEnd, err := time.Parse("15:04", schedule.BreakEnd[:5])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
				Status:    "500",
				Message:   "Invalid break end time",
				MessageTh: "เวลาสิ้นสุดพักไม่ถูกต้อง",
				Error:     err.Error(),
			})
		}

		// ตรวจว่าเวลานัดทับกับช่วงพักหรือไม่
		if requestStart.Before(breakEnd) && requestEnd.After(breakStart) {
			return c.Status(fiber.StatusOK).JSON(dto.Response{
				Status:    "200",
				Message:   "Appointment time is during doctor break",
				MessageTh: "เวลาที่เลือกอยู่ในช่วงพักของแพทย์",
				Data: fiber.Map{
					"available":   false,
					"doctor_id":   doctorID,
					"date":        date,
					"start_time":  startTime,
					"end_time":    endTime,
					"break_start": schedule.BreakStart[:5],
					"break_end":   schedule.BreakEnd[:5],
				},
			})
		}
	}

	// =========================
	// ตรวจสอบ Appointment ที่มีอยู่
	// =========================
	var count int64

	err = database.DBConn.
		Model(&models.Appointment{}).
		Where(`
			doctor_id = ?
			AND DATE(appointment_date) = ?
			AND status = ?
			AND start_time < ?
			AND end_time > ?
		`,
			doctorID,
			date,
			StatusBooked,
			endTime,
			startTime,
		).
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to check appointment",
			MessageTh: "ไม่สามารถตรวจสอบเวลานัดหมายได้",
			Error:     err.Error(),
		})
	}

	// =========================
	// มี Appointment ชน
	// =========================
	if count > 0 {
		return c.Status(fiber.StatusOK).JSON(dto.Response{
			Status:    "200",
			Message:   "Appointment time is already booked",
			MessageTh: "วันที่และเวลาที่ค้นหามีการจองแล้ว",
			Data: fiber.Map{
				"available":  false,
				"doctor_id":  doctorID,
				"date":       date,
				"start_time": startTime,
				"end_time":   endTime,
			},
		})
	}

	// =========================
	// ยังว่าง
	// =========================
	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:    "200",
		Message:   "Appointment time is available",
		MessageTh: "วันที่และเวลาที่ค้นหายังว่าง",
		Data: fiber.Map{
			"available":  true,
			"doctor_id":  doctorID,
			"date":       date,
			"start_time": startTime,
			"end_time":   endTime,
		},
	})
}
