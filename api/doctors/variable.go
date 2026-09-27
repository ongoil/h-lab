package doctors

type DoctorRequest struct {
	TitleName string `json:"title_name"` // คำนำหน้าชื่อ
	FirstName string `json:"first_name"` // ชื่อ
	LastName  string `json:"last_name"`  // นามสกุล
}

type responseDoctor struct {
	TitleName string `json:"title_name"` // คำนำหน้าชื่อ
	FirstName string `json:"first_name"` // ชื่อ
	LastName  string `json:"last_name"`  // นามสกุล
}
