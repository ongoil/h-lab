package departments

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
)

func CreateDepartment(c *fiber.Ctx) error {
	var req departmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
	}

	tx := database.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.ENGName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid request eng_name",
			MessageTh: "กรุณาระบุชื่อภาษาอังกฤษ",
		})
	}

	if req.THName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid request th_name",
			MessageTh: "กรุณาระบุชื่อภาษาไทย",
		})
	}

	department := models.Department{
		ENGName: req.ENGName,
		THName:  req.THName,
	}

	if err := tx.Create(&department).Error; err != nil {
		tx.Rollback()
		log.Error(c, "500 | Internal Server Error : failed to create department -> ", err, logrus.Fields{"department": req.ENGName})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to create department",
			MessageTh: "ไม่สามารถสร้างข้อมูลแผนได้",
			Error:     err.Error(),
		})
	}

	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Error(c, "500 | Internal Server Error : failed to commit transaction -> ", err, logrus.Fields{})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "Internal Server Error - POST transactions",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:    "200",
		Message:   "Success",
		MessageTh: "สำเร็จ",
	})
}

func GetDepartment(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	baseQuery := database.DBConn.Model(&models.Department{})
	if search != "" {
		// searchWithoutSpaces := strings.ReplaceAll(search, " ", "")
		baseQuery = baseQuery.Where(`INSTR(REPLACE(CONCAT_WS('|',eng_name,th_name ), ' ', ''), ?)`, search)
	}

	var department []models.Department
	if err := baseQuery.Find(&department).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to get department",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:    "200",
		Message:   "Success",
		MessageTh: "สำเร็จ",
		Data:      department,
	})
}
