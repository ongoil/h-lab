package departments

type departmentRequest struct {
	ENGName string `json:"eng_name"` // ชื่อ eng
	THName  string `json:"th_name"`  // ชื่อ TH
}

type reponseDepartment struct {
	ENGName string `json:"eng_name"` // ชื่อ eng
	THName  string `json:"th_name"`  // ชื่อ TH
}
