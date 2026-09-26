package database

import (
	"h-lab/models"

	"gorm.io/gorm"
)

func DataBaseMigration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Doctor{},
		&models.Department{},
		&models.DoctorSchedule{},
		&models.Patient{},
	)
	if err != nil {
		return err
	}

	err = db.Exec(`
		ALTER TABLE doctor_schedules
		MODIFY COLUMN start_time TIME NOT NULL,
		MODIFY COLUMN end_time TIME NOT NULL,
		MODIFY COLUMN break_start TIME NULL,
		MODIFY COLUMN break_end TIME NULL
	`).Error
	if err != nil {
		return err
	}

	return nil
}
