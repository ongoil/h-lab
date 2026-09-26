package appointments

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"time"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
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

	if req.AppointmentTypeID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Appointment type ID is required",
			MessageTh: "กรุณาระบุประเภทการนัดหมาย",
		})
	}

	if req.StartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Start time is required",
			MessageTh: "กรุณาระบุเวลาเริ่มนัดหมาย",
		})
	}

	if req.CreatedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Created by is required",
			MessageTh: "กรุณาระบุผู้สร้างรายการ",
		})
	}

	// ห้ามจองวันในอดีต
	today := time.Now().Truncate(24 * time.Hour)

	if req.AppointmentDate.Before(today) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Cannot book an appointment in the past",
			MessageTh: "ไม่สามารถจองนัดหมายย้อนหลังได้",
		})
	}

	// -------------------------
	// Transaction
	// -------------------------

	tx := database.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// -------------------------
	// Check Patient
	// -------------------------
	var patient models.Patient
	if err := tx.Where("id = ?", req.PatientID).First(&patient).Error; err != nil {
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
	if err := tx.Where("id = ?", req.DoctorID).First(&doctor).Error; err != nil {
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
	if err := tx.Where("id = ?", req.DepartmentID).First(&department).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Department not found",
			MessageTh: "ไม่พบข้อมูลแผนก",
		})
	}

	// -------------------------
	// Check Appointment Type
	// -------------------------
	var appointmentType models.AppointmentType
	if err := tx.Where("id = ? AND is_active = ?", req.AppointmentTypeID, true).First(&appointmentType).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Appointment type not found or inactive",
			MessageTh: "ไม่พบประเภทการนัดหมายหรือประเภทนี้ถูกปิดใช้งาน",
		})
	}

	// -------------------------
	// TODO:
	// Check Doctor Schedule
	// Check Break Time
	// Check Existing Appointment
	// Calculate EndTime
	// -------------------------
	appointment := models.Appointment{
		PatientID:         patient.Id,
		DoctorID:          doctor.Id,
		DepartmentID:      department.Id,
		AppointmentTypeID: appointmentType.Id,
		AppointmentDate:   req.AppointmentDate,
		StartTime:         req.StartTime,
		Status:            StatusBooked,
		Reason:            req.Reason,
		CreatedBy:         req.CreatedBy,
	}

	// -------------------------
	// Create Appointment
	// -------------------------
	if err := tx.Create(&appointment).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to create appointment",
			MessageTh: "ไม่สามารถสร้างการนัดหมายได้",
			Error:     err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to commit transaction",
			MessageTh: "ไม่สามารถบันทึก Transaction ได้",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Status:    "201",
		Message:   "Appointment created successfully",
		MessageTh: "สร้างการนัดหมายสำเร็จ",
		Data:      appointment,
	})
}
