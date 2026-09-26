package patients

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreatePatient(c *fiber.Ctx) error {
	var req patientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Invalid request body",
			MessageTh: "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}

	// Validate required fields
	if req.PatientNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Patient number is required",
			MessageTh: "กรุณาระบุเลขประจำตัวผู้ป่วย",
		})
	}

	if req.TitleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Title name is required",
			MessageTh: "กรุณาระบุคำนำหน้าชื่อ",
		})
	}

	if req.FirstName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "First name is required",
			MessageTh: "กรุณาระบุชื่อ",
		})
	}

	if req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Last name is required",
			MessageTh: "กรุณาระบุนามสกุล",
		})
	}

	// Validate date of birth
	if req.DateOfBirth.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Date of birth is required",
			MessageTh: "กรุณาระบุวันเกิด",
		})
	}

	if req.DateOfBirth.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "Date of birth cannot be in the future",
			MessageTh: "วันเกิดไม่สามารถเป็นวันที่ในอนาคตได้",
		})
	}

	// Start transaction
	tx := database.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check duplicate PatientNo
	var existingPatient models.Patient

	if err := tx.Where("patient_no = ?", req.PatientNo).First(&existingPatient).Error; err == nil {
		tx.Rollback()
		return c.Status(fiber.StatusConflict).JSON(dto.Response{
			Status:    "409",
			Message:   "Patient number already exists",
			MessageTh: "เลขประจำตัวผู้ป่วยนี้มีอยู่แล้ว",
		})
	}

	// Create patient
	patient := models.Patient{
		PatientNo:   req.PatientNo,
		TitleName:   req.TitleName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
	}

	if err := tx.Create(&patient).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to create patient",
			MessageTh: "ไม่สามารถสร้างข้อมูลผู้ป่วยได้",
			Error:     err.Error(),
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Failed to commit transaction",
			MessageTh: "ไม่สามารถบันทึกข้อมูลได้",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Status:    "201",
		Message:   "Patient created successfully",
		MessageTh: "สร้างข้อมูลผู้ป่วยสำเร็จ",
		Data:      patient,
	})
}
