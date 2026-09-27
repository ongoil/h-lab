package patients

import (
	"fmt"
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"strconv"
	"strings"
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

	// Generate Patient No: HN + year + running number
	year := time.Now().Year()
	prefix := fmt.Sprintf("HN%d", year)
	var lastPatient models.Patient
	err := tx.Where("patient_no LIKE ?", prefix+"%").Order("patient_no DESC").First(&lastPatient).Error
	nextNo := 1
	if err == nil {
		lastNo, err := strconv.Atoi(lastPatient.PatientNo[len(prefix):])
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
				Status:    "500",
				Message:   "Failed to generate patient number",
				MessageTh: "ไม่สามารถสร้างเลขประจำตัวผู้ป่วยได้",
				Error:     err.Error(),
			})
		}
		nextNo = lastNo + 1
	}

	patientNo := fmt.Sprintf("%s%04d", prefix, nextNo)
	dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	// Create patient
	patient := models.Patient{
		PatientNo:   patientNo,
		TitleName:   req.TitleName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: dateOfBirth,
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

func GetPatient(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	baseQuery := database.DBConn.Model(&models.Patient{})
	if search != "" {
		// searchWithoutSpaces := strings.ReplaceAll(search, " ", "")
		baseQuery = baseQuery.Where(`INSTR(REPLACE(CONCAT_WS('|',patient_no,first_name,last_name ), ' ', ''), ?)`, search)
	}

	var patient []models.Patient
	if err := baseQuery.Find(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to get patient",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:    "200",
		Message:   "Success",
		MessageTh: "สำเร็จ",
		Data:      patient,
	})
}
