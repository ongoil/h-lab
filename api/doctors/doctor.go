package doctors

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

func CreateDoctor(c *fiber.Ctx) error {
	var req DoctorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid request body",
			MessageTh: "ข้อมูลไม่ถูกต้อง",
		})
	}

	if req.TitleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid title name",
			MessageTh: "กรุณาระบุคำนำหน้าชื่อ",
		})
	}

	if req.FirstName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid first name",
			MessageTh: "กรุณาระบุชื่อ",
		})
	}

	if req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Response{
			Status:    "400",
			Message:   "invalid last name",
			MessageTh: "กรุณาระบุนามสกุล",
		})
	}

	tx := database.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	doctor := models.Doctor{
		TitleName: req.TitleName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := tx.Create(&doctor).Error; err != nil {
		tx.Rollback()
		log.Error(c, "500 | Internal Server Error : failed to create doctor -> ", err, logrus.Fields{"first name": req.FirstName})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to create doctor",
			MessageTh: "ไม่สามารถสร้างข้อมูลแพทย์ได้",
		})
	}

	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Error(c, "500 | Internal Server Error : failed to commit transaction -> ", err, logrus.Fields{"first name": req.FirstName})
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
