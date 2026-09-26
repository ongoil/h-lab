package schedules

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
)

func CreateSchedules(c *fiber.Ctx) error {
	var req doctorScheduleRequest
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

	var doctor models.Doctor
	if err := database.DBConn.Model(&models.Doctor{}).First(&doctor).Error; err != nil {
		log.Error(c, "500 | Internal Server Error : failed to find doctor -> ", err, logrus.Fields{"id": req.DoctorID})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to find doctor",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
	}

	var department models.Department
	if err := database.DBConn.Model(&models.Department{}).First(&department).Error; err != nil {
		log.Error(c, "500 | Internal Server Error : failed to find department -> ", err, logrus.Fields{"id": req.DepartmentID})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to find department",
			MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
			Error:     err.Error(),
		})
	}

	schedule := models.DoctorSchedule{
		DoctorID:     doctor.Id,       // รหัสแพทย์
		DepartmentID: department.Id,   // รหัสแผนก
		DayOfWeek:    req.DayOfWeek,   // 0 = อาทิตย์, 1 = จันทร์, ..., 6 = เสาร์
		StartTime:    req.StartTime,   // เวลาเริ่มทำงาน
		EndTime:      req.EndTime,     // เวลาสิ้นสุดการทำงาน
		BreakStart:   req.BreakStart,  // เวลาเริ่มพัก
		BreakEnd:     req.BreakEnd,    // เวลาสิ้นสุดพัก
		IsAvailable:  req.IsAvailable, // เปิดให้จองหรือไม่
	}

	if err := tx.Create(&schedule).Error; err != nil {
		tx.Rollback()
		log.Error(c, "500 | Internal Server Error : failed to create doctor schedule -> ", err, logrus.Fields{"doctor_id": doctor.Id})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to create doctor schedule",
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
