package patients

type patientRequest struct {
	PatientNo   string `json:"patient_no"`
	TitleName   string `json:"title_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
}
