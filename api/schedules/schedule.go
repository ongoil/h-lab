package schedules

import (
	"h-lab/database"
	"h-lab/dto"
	"h-lab/models"
	"strconv"
	"strings"

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

	// var doctor models.Doctor
	// if err := database.DBConn.Model(&models.Doctor{}).First(&doctor).Error; err != nil {
	// 	log.Error(c, "500 | Internal Server Error : failed to find doctor -> ", err, logrus.Fields{"id": req.DoctorID})
	// 	return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
	// 		Status:    "500",
	// 		Message:   "failed to find doctor",
	// 		MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
	// 		Error:     err.Error(),
	// 	})
	// }

	// var department models.Department
	// if err := database.DBConn.Model(&models.Department{}).First(&department).Error; err != nil {
	// 	log.Error(c, "500 | Internal Server Error : failed to find department -> ", err, logrus.Fields{"id": req.DepartmentID})
	// 	return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
	// 		Status:    "500",
	// 		Message:   "failed to find department",
	// 		MessageTh: "ข้อผิดพลาดภายในเซิร์ฟเวอร์",
	// 		Error:     err.Error(),
	// 	})
	// }

	schedule := models.DoctorSchedule{
		DoctorID:     req.DoctorID,     // รหัสแพทย์
		DepartmentID: req.DepartmentID, // รหัสแผนก
		DayOfWeek:    req.DayOfWeek,    // 0 = อาทิตย์, 1 = จันทร์, ..., 6 = เสาร์
		StartTime:    req.StartTime,    // เวลาเริ่มทำงาน
		EndTime:      req.EndTime,      // เวลาสิ้นสุดการทำงาน
		BreakStart:   req.BreakStart,   // เวลาเริ่มพัก
		BreakEnd:     req.BreakEnd,     // เวลาสิ้นสุดพัก
		IsAvailable:  req.IsAvailable,  // เปิดให้จองหรือไม่
	}

	if err := tx.Create(&schedule).Error; err != nil {
		tx.Rollback()
		log.Error(c, "500 | Internal Server Error : failed to create doctor schedule -> ", err, logrus.Fields{"doctor_id": req.DoctorID})
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

func GetDoctorSchedules(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	query := database.DBConn.
		Model(&models.DoctorSchedule{}).
		Preload("Doctor").
		Preload("Department")

	if search != "" {
		searchWithoutSpaces := strings.ReplaceAll(search, " ", "")

		query = query.
			Joins("JOIN doctors ON doctors.id = doctor_schedules.doctor_id").
			Where(
				`INSTR(
                REPLACE(
                    CONCAT_WS('|', doctors.first_name, doctors.last_name),
                    ' ',
                    ''
                ),
                ?
            ) > 0`,
				searchWithoutSpaces,
			)
	}
	// นับจำนวนข้อมูลทั้งหมดก่อนแบ่งหน้า
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to count doctor schedules",
			MessageTh: "ไม่สามารถนับจำนวนตารางแพทย์ได้",
			Error:     err.Error(),
		})
	}

	var schedules []models.DoctorSchedule

	if err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&schedules).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to get doctor schedules",
			MessageTh: "ไม่สามารถดึงข้อมูลตารางแพทย์ได้",
			Error:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Status:    "200",
		Message:   "success",
		MessageTh: "ดึงข้อมูลสำเร็จ",
		Data: fiber.Map{
			"items":     schedules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func DeleteSchedules(c *fiber.Ctx) error {
	Id := strings.TrimSpace(c.Query("id"))
	tx := database.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&models.DoctorSchedule{}, "id = ?", Id).Error; err != nil {
		tx.Rollback()
		log.Error(c, "500 | Internal Server Error : failed to delete doctor schedule -> ", err, logrus.Fields{"id": Id})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Status:    "500",
			Message:   "failed to delete doctor schedule",
			MessageTh: "ไม่สามารถลบตารางเวลาของแพทย์ได้",
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
		Message:   "Delete Success",
		MessageTh: "ลบข้อมูลสำเร็จ",
	})
}
