package router

import (
	"h-lab/api/appointments"
	"h-lab/api/departments"
	"h-lab/api/doctors"
	"h-lab/api/patients"
	"h-lab/api/schedules"
	"h-lab/api/test"
	"h-lab/framework/fiber"
)

func SetRouter(app *fiber.FiberApp) {
	v1 := app.Group("/backend/api/v1")
	v1.Get("/test", test.Test)

	doctor := v1.Group("/doctors")
	doctor.Post("/create", doctors.CreateDoctor)
	doctor.Get("/show", doctors.GetDoctor)

	department := v1.Group("/departments")
	department.Post("/create", departments.CreateDepartment)
	department.Get("/show", departments.GetDepartment)

	schedule := v1.Group("/schedules")
	schedule.Post("/create", schedules.CreateSchedules)
	schedule.Get("/show", schedules.GetDoctorSchedules)
	schedule.Delete("/delete", schedules.DeleteSchedules)

	appointment := v1.Group("/appointments")
	appointment.Post("/create", appointments.CreateAppointment)
	appointment.Get("/show", appointments.GetAppointment)

	patient := v1.Group("/patients")
	patient.Post("/create", patients.CreatePatient)
	patient.Get("/show", patients.GetPatient)

}
