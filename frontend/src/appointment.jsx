import { useEffect, useState } from "react";
import "./App.css";

const dayNames = [
    "อาทิตย์",
    "จันทร์",
    "อังคาร",
    "พุธ",
    "พฤหัสบดี",
    "ศุกร์",
    "เสาร์",
];

function AppointmentPage() {
    const [doctors, setDoctors] = useState([]);
    const [departments, setDepartments] = useState([]);
    const [patients, setPatients] = useState([]);
    const [patientSearch, setPatientSearch] = useState("");
    const [schedules, setSchedules] = useState([]);
    const [scheduleSearch, setScheduleSearch] = useState("");
    const [schedulePage, setSchedulePage] = useState(1);
    const [schedulePageSize, setSchedulePageSize] = useState(10);
    const [scheduleTotal, setScheduleTotal] = useState(0);
    const [availabilityResult, setAvailabilityResult] = useState(null);
    const [availabilityLoading, setAvailabilityLoading] = useState(false);
    const [showScheduleModal, setShowScheduleModal] = useState(false);
    const [creatingSchedule, setCreatingSchedule] = useState(false);
    const [showModal, setShowModal] = useState(false);
    const [creating, setCreating] = useState(false);
    const [showPatientModal, setShowPatientModal] = useState(false);
    const [creatingPatient, setCreatingPatient] = useState(false);
    const [showDoctorModal, setShowDoctorModal] = useState(false);
    const [creatingDoctor, setCreatingDoctor] = useState(false);
    const [showDepartmentModal, setShowDepartmentModal] = useState(false);
    const [creatingDepartment, setCreatingDepartment] = useState(false);
    const [doctorSearch, setDoctorSearch] = useState("");
    const [showDoctorDropdown, setShowDoctorDropdown] = useState(false);


    const filteredDoctors = doctors.filter((doctor) => {
        const fullName = `${doctor.title_name} ${doctor.first_name} ${doctor.last_name}`;

        return fullName
            .toLowerCase()
            .includes(doctorSearch.toLowerCase());
    });

    const formatTime = (time) => {
        if (!time) return "-";
        return time.slice(0, 5).replace(":", ".");
    };

    const [departmentForm, setDepartmentForm] = useState({
        eng_name: "",
        th_name: "",
    });

    const [doctorForm, setDoctorForm] = useState({
        title_name: "",
        first_name: "",
        last_name: "",
    });

    const [patientForm, setPatientForm] = useState({
        title_name: "",
        first_name: "",
        last_name: "",
        date_of_birth: "",
    });

    const [scheduleForm, setScheduleForm] = useState({
        doctor_id: "",
        department_id: "",
        day_of_week: "",
        start_time: "",
        end_time: "",
        break_start: "",
        break_end: "",
        is_available: true,
    });

    const [availabilityForm, setAvailabilityForm] = useState({
        doctor_id: "",
        date: "",
        start_time: "",
        end_time: "",
    });

    const [form, setForm] = useState({
        patient_id: "",
        doctor_id: "",
        department_id: "",
        appointment_type: "",
        appointment_date: "",
        start_time: "",
        end_time: "",
        reason: "",
        created_by: "admin",
    });

    const handleCreateDepartment = async (e) => {
        e.preventDefault();

        if (!departmentForm.eng_name || !departmentForm.th_name) {
            alert("กรุณากรอกข้อมูลแผนกให้ครบถ้วน");
            return;
        }

        try {
            setCreatingDepartment(true);

            const response = await fetch(
                "http://127.0.0.1:8080/backend/api/v1/departments/create",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(departmentForm),
                }
            );

            const result = await response.json();

            if (!response.ok) {
                alert(
                    result.message_th ||
                    result.message ||
                    "ไม่สามารถเพิ่มข้อมูลแผนกได้"
                );
                return;
            }

            alert(result.message_th || "เพิ่มข้อมูลแผนกสำเร็จ");

            setShowDepartmentModal(false);

            setDepartmentForm({
                eng_name: "",
                th_name: "",
            });

            window.location.reload();
        } catch (error) {
            console.error("Create department error:", error);
            alert("ไม่สามารถเชื่อมต่อกับระบบได้");
        } finally {
            setCreatingDepartment(false);
        }
    };

    const handleCreateDoctor = async (e) => {
        e.preventDefault();

        if (
            !doctorForm.title_name ||
            !doctorForm.first_name ||
            !doctorForm.last_name
        ) {
            alert("กรุณากรอกข้อมูลแพทย์ให้ครบถ้วน");
            return;
        }

        try {
            setCreatingDoctor(true);

            const response = await fetch(
                "http://127.0.0.1:8080/backend/api/v1/doctors/create",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(doctorForm),
                }
            );

            const result = await response.json();

            if (!response.ok) {
                alert(
                    result.message_th ||
                    result.message ||
                    "ไม่สามารถเพิ่มข้อมูลแพทย์ได้"
                );
                return;
            }

            alert(result.message_th || "เพิ่มข้อมูลแพทย์สำเร็จ");

            setShowDoctorModal(false);

            setDoctorForm({
                title_name: "",
                first_name: "",
                last_name: "",
            });

            // Refresh หน้า
            window.location.reload();

        } catch (error) {
            console.error("Create doctor error:", error);
            alert("ไม่สามารถเชื่อมต่อกับระบบได้");
        } finally {
            setCreatingDoctor(false);
        }
    };

    const handleCreateSchedule = async (e) => {
        e.preventDefault();

        try {
            setCreatingSchedule(true);

            const response = await fetch(
                "http://127.0.0.1:8080/backend/api/v1/schedules/create",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        ...scheduleForm,
                        day_of_week: Number(scheduleForm.day_of_week),
                    }),
                }
            );

            const result = await response.json();

            if (!response.ok) {
                alert(
                    result.message_th ||
                    result.message ||
                    "ไม่สามารถเพิ่มตารางเวลาได้"
                );
                return;
            }

            alert(result.message_th || "เพิ่มตารางเวลาสำเร็จ");

            // Refresh หน้า
            window.location.reload();

        } catch (error) {
            console.error("Create schedule error:", error);
            alert("ไม่สามารถเชื่อมต่อกับระบบได้");
        } finally {
            setCreatingSchedule(false);
        }
    };

    useEffect(() => {
        fetch("http://localhost:8080/backend/api/v1/doctors/show")
            .then((response) => response.json())
            .then((result) => {
                setDoctors(result.data || []);
            })
            .catch((error) => {
                console.error("Error fetching doctors:", error);
            });
    }, []);

    useEffect(() => {
        fetch("http://127.0.0.1:8080/backend/api/v1/departments/show")
            .then((response) => response.json())
            .then((result) => {
                setDepartments(result.data || result);
            })
            .catch((error) => {
                console.error("Error fetching departments:", error);
            });
    }, []);

    useEffect(() => {
        fetch("http://127.0.0.1:8080/backend/api/v1/schedules/show")
            .then((response) => response.json())
            .then((result) => {
                setSchedules(result.data?.items || []);
                setScheduleTotal(result.data?.total || 0);
            })
            .catch((error) => {
                console.error("Error fetching schedules:", error);
                setSchedules([]);
                setScheduleTotal(0);
            });
    }, []);

    useEffect(() => {
        const fetchSchedules = async () => {
            try {
                const params = new URLSearchParams({
                    search: scheduleSearch,
                    page: schedulePage,
                    size: schedulePageSize,
                });

                const response = await fetch(
                    `http://127.0.0.1:8080/backend/api/v1/schedules/show?${params.toString()}`
                );

                const result = await response.json();

                if (!response.ok) {
                    setSchedules([]);
                    setScheduleTotal(0);
                    return;
                }

                setSchedules(result.data?.items || []);
                setScheduleTotal(result.data?.total || 0);

            } catch (error) {
                console.error("Error fetching schedules:", error);
                setSchedules([]);
                setScheduleTotal(0);
            }
        };

        fetchSchedules();
    }, [scheduleSearch, schedulePage, schedulePageSize]);

    const handleCreatePatient = async (e) => {
        e.preventDefault();

        if (
            !patientForm.title_name ||
            !patientForm.first_name ||
            !patientForm.last_name ||
            !patientForm.date_of_birth
        ) {
            alert("กรุณากรอกข้อมูลผู้ป่วยให้ครบถ้วน");
            return;
        }

        try {
            setCreatingPatient(true);

            const response = await fetch(
                "http://127.0.0.1:8080/backend/api/v1/patients/create",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(patientForm),
                }
            );

            const result = await response.json();

            if (!response.ok) {
                alert(
                    result.message_th ||
                    result.message ||
                    "ไม่สามารถเพิ่มข้อมูลผู้ป่วยได้"
                );
                return;
            }

            alert(result.message_th || "เพิ่มข้อมูลผู้ป่วยสำเร็จ");

            setShowPatientModal(false);

            setPatientForm({
                title_name: "",
                first_name: "",
                last_name: "",
                date_of_birth: "",
            });

            // Refresh หน้า
            window.location.reload();

        } catch (error) {
            console.error("Create patient error:", error);
            alert("ไม่สามารถเชื่อมต่อกับระบบได้");
        } finally {
            setCreatingPatient(false);
        }
    };

    const searchPatients = async (value) => {
        setPatientSearch(value);

        if (!value.trim()) {
            setPatients([]);
            return;
        }

        try {
            const response = await fetch(
                `http://localhost:8080/backend/api/v1/patients/show?search=${encodeURIComponent(value)}`
            );

            const result = await response.json();

            setPatients(result.data || []);
        } catch (error) {
            console.error("Search patients error:", error);
            setPatients([]);
        }
    };


    const createAppointment = async () => {
        try {
            setCreating(true);

            const response = await fetch(
                "http://127.0.0.1:8080/backend/api/v1/appointments/create",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(form),
                }
            );

            const result = await response.json();

            console.log("Create appointment response:", result);

            if (!response.ok) {
                alert(result.message_th || result.message || "ไม่สามารถสร้างการนัดหมายได้");
                return;
            }

            alert(result.message_th);

            setShowModal(false);

        } catch (error) {
            console.error("Create appointment error:", error);

            alert("ไม่สามารถเชื่อมต่อกับระบบได้");

        } finally {
            setCreating(false);
        }
    };

    const handleDeleteSchedule = async (id) => {
        const confirmDelete = window.confirm(
            "คุณต้องการลบตารางเวลานี้ใช่หรือไม่?"
        );

        if (!confirmDelete) return;

        try {
            const response = await fetch(
                `http://127.0.0.1:8080/backend/api/v1/schedules/delete?id=${id}`,
                {
                    method: "DELETE",
                }
            );

            const result = await response.json();

            if (!response.ok) {
                alert(
                    result.message_th ||
                    result.message ||
                    "ไม่สามารถลบตารางเวลาได้"
                );
                return;
            }

            alert(result.message_th || "ลบตารางเวลาสำเร็จ");

            // Refresh หน้า
            window.location.reload();

        } catch (error) {
            console.error("Delete schedule error:", error);
            alert("ไม่สามารถเชื่อมต่อกับระบบได้");
        }
    };

    const checkAvailability = async () => {
        const {
            doctor_id,
            date,
            start_time,
            end_time,
        } = availabilityForm;

        if (!doctor_id || !date || !start_time || !end_time) {
            alert("กรุณาระบุแพทย์ วันที่ เวลาเริ่มต้น และเวลาสิ้นสุด");
            return;
        }

        try {
            setAvailabilityLoading(true);
            setAvailabilityResult(null);

            const params = new URLSearchParams({
                doctor_id,
                date,
                start_time,
                end_time,
            });

            const response = await fetch(
                `http://127.0.0.1:8080/backend/api/v1/appointments/show?${params.toString()}`
            );

            const result = await response.json();

            setAvailabilityResult(result);

        } catch (error) {
            console.error("Check availability error:", error);

            setAvailabilityResult({
                status: "500",
                message_th: "ไม่สามารถตรวจสอบเวลานัดหมายได้",
                data: {
                    available: false,
                },
            });
        } finally {
            setAvailabilityLoading(false);
        }
    };

    const handleAvailabilityTimeChange = (e) => {
        let value = e.target.value.replace(/\D/g, "").slice(0, 4);

        if (value.length > 2) {
            value = value.slice(0, 2) + ":" + value.slice(2);
        }

        // แสดงค่าที่ format แล้วใน input
        e.target.value = value;

        setAvailabilityForm((prev) => ({
            ...prev,
            [e.target.name]: value,
        }));
    };

    const handleScheduleTimeChange = (e, field) => {
        let value = e.target.value.replace(/\D/g, "");

        value = value.slice(0, 4);

        if (value.length > 2) {
            value = value.slice(0, 2) + ":" + value.slice(2);
        }

        // เก็บค่าที่ผู้ใช้กรอกลง scheduleForm
        setScheduleForm((prev) => ({
            ...prev,
            [field]: value,
        }));
    };

    const handleChange = (e) => {
        setForm({
            ...form,
            [e.target.name]: e.target.value,
        });
    };

    const handleTimeChange = (e) => {
        const { name, value } = e.target;

        // เอาเฉพาะตัวเลข
        let time = value.replace(/\D/g, "");

        // จำกัด 4 ตัว
        time = time.slice(0, 4);

        // เติม :
        if (time.length > 2) {
            time = time.slice(0, 2) + ":" + time.slice(2);
        }

        setForm((prev) => ({
            ...prev,
            [name]: time,
        }));
    };

    return (
        <main className="main">

            <div className="section-title">
                <div>
                    <h2>จัดการข้อมูล</h2>
                    <p>เพิ่มข้อมูลผู้ป่วยและแพทย์เข้าสู่ระบบ</p>
                </div>

                <div className="section-actions">
                    <button
                        type="button"
                        className="primary-button"
                        onClick={() => setShowDepartmentModal(true)}
                    >
                        + เพิ่มข้อมูลแผนก
                    </button>

                    <button
                        type="button"
                        className="primary-button"
                        onClick={() => setShowPatientModal(true)}
                    >
                        + เพิ่มข้อมูลผู้ป่วย
                    </button>

                    <button
                        type="button"
                        className="primary-button"
                        onClick={() => setShowDoctorModal(true)}
                    >
                        + เพิ่มข้อมูลแพทย์
                    </button>
                </div>

                {showPatientModal && (
                    <div
                        className="modal-overlay"
                        onClick={() => setShowPatientModal(false)}
                    >
                        <div
                            className="modal patient-modal"
                            onClick={(e) => e.stopPropagation()}
                        >
                            {/* Header */}
                            <div className="modal-header">
                                <div>
                                    <h2>เพิ่มข้อมูลผู้ป่วย</h2>
                                    <p>กรอกข้อมูลพื้นฐานของผู้ป่วย</p>
                                </div>

                                <button
                                    type="button"
                                    className="modal-close"
                                    onClick={() => setShowPatientModal(false)}
                                >
                                    ×
                                </button>
                            </div>

                            <form onSubmit={handleCreatePatient}>
                                <div className="form-section">
                                    <div className="form-section-title">
                                        <h3>ข้อมูลผู้ป่วย</h3>
                                        <p>กรุณากรอกข้อมูลให้ครบถ้วน</p>
                                    </div>

                                    {/* คำนำหน้า */}
                                    <div className="form-group">
                                        <label>
                                            คำนำหน้า <span>*</span>
                                        </label>

                                        <select
                                            value={patientForm.title_name}
                                            onChange={(e) =>
                                                setPatientForm((prev) => ({
                                                    ...prev,
                                                    title_name: e.target.value,
                                                }))
                                            }
                                        >
                                            <option value="">เลือกคำนำหน้า</option>
                                            <option value="นาย">นาย</option>
                                            <option value="นาง">นาง</option>
                                            <option value="นางสาว">นางสาว</option>
                                            <option value="ด.ช.">ด.ช.</option>
                                            <option value="ด.ญ.">ด.ญ.</option>
                                        </select>
                                    </div>

                                    {/* ชื่อ / นามสกุล */}
                                    <div className="form-row">
                                        <div className="form-group">
                                            <label>
                                                ชื่อ <span>*</span>
                                            </label>

                                            <input
                                                type="text"
                                                placeholder="เช่น กิตติ"
                                                value={patientForm.first_name}
                                                onChange={(e) =>
                                                    setPatientForm((prev) => ({
                                                        ...prev,
                                                        first_name: e.target.value,
                                                    }))
                                                }
                                            />
                                        </div>

                                        <div className="form-group">
                                            <label>
                                                นามสกุล <span>*</span>
                                            </label>

                                            <input
                                                type="text"
                                                placeholder="เช่น พัฒนากุล"
                                                value={patientForm.last_name}
                                                onChange={(e) =>
                                                    setPatientForm((prev) => ({
                                                        ...prev,
                                                        last_name: e.target.value,
                                                    }))
                                                }
                                            />
                                        </div>
                                    </div>

                                    {/* วันเกิด */}
                                    <div className="form-group">
                                        <label>
                                            วันเกิด <span>*</span>
                                        </label>

                                        <input
                                            type="date"
                                            value={patientForm.date_of_birth}
                                            onChange={(e) =>
                                                setPatientForm((prev) => ({
                                                    ...prev,
                                                    date_of_birth: e.target.value,
                                                }))
                                            }
                                        />
                                    </div>
                                </div>

                                {/* Actions */}
                                <div className="modal-actions">
                                    <button
                                        type="button"
                                        className="cancel-button"
                                        onClick={() => setShowPatientModal(false)}
                                    >
                                        ยกเลิก
                                    </button>

                                    <button
                                        type="submit"
                                        className="primary-button"
                                        disabled={creatingPatient}
                                    >
                                        {creatingPatient
                                            ? "กำลังบันทึก..."
                                            : "บันทึกข้อมูลผู้ป่วย"}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {showDoctorModal && (
                    <div
                        className="modal-overlay"
                        onClick={() => setShowDoctorModal(false)}
                    >
                        <div
                            className="modal doctor-modal"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <div className="modal-header">
                                <div>
                                    <h2>เพิ่มข้อมูลแพทย์</h2>
                                    <p>กรอกข้อมูลพื้นฐานของแพทย์</p>
                                </div>

                                <button
                                    type="button"
                                    className="modal-close"
                                    onClick={() => setShowDoctorModal(false)}
                                >
                                    ×
                                </button>
                            </div>

                            <form onSubmit={handleCreateDoctor}>
                                <div className="form-section">
                                    <div className="form-section-title">
                                        <h3>ข้อมูลแพทย์</h3>
                                        <p>กรุณากรอกข้อมูลให้ครบถ้วน</p>
                                    </div>

                                    {/* คำนำหน้า */}
                                    <div className="form-group">
                                        <label>
                                            คำนำหน้า <span>*</span>
                                        </label>

                                        <select
                                            value={doctorForm.title_name}
                                            onChange={(e) =>
                                                setDoctorForm((prev) => ({
                                                    ...prev,
                                                    title_name: e.target.value,
                                                }))
                                            }
                                        >
                                            <option value="">เลือกคำนำหน้า</option>
                                            <option value="นาย">นาย</option>
                                            <option value="นาง">นาง</option>
                                            <option value="นางสาว">นางสาว</option>
                                            <option value="นพ.">นพ.</option>
                                            <option value="พญ.">พญ.</option>
                                            <option value="ดร.">ดร.</option>
                                        </select>
                                    </div>

                                    {/* ชื่อ / นามสกุล */}
                                    <div className="form-row">
                                        <div className="form-group">
                                            <label>
                                                ชื่อ <span>*</span>
                                            </label>

                                            <input
                                                type="text"
                                                placeholder="เช่น ทดสอบ3"
                                                value={doctorForm.first_name}
                                                onChange={(e) =>
                                                    setDoctorForm((prev) => ({
                                                        ...prev,
                                                        first_name: e.target.value,
                                                    }))
                                                }
                                            />
                                        </div>

                                        <div className="form-group">
                                            <label>
                                                นามสกุล <span>*</span>
                                            </label>

                                            <input
                                                type="text"
                                                placeholder="เช่น ทดสอบ3"
                                                value={doctorForm.last_name}
                                                onChange={(e) =>
                                                    setDoctorForm((prev) => ({
                                                        ...prev,
                                                        last_name: e.target.value,
                                                    }))
                                                }
                                            />
                                        </div>
                                    </div>
                                </div>

                                <div className="modal-actions">
                                    <button
                                        type="button"
                                        className="cancel-button"
                                        onClick={() => setShowDoctorModal(false)}
                                    >
                                        ยกเลิก
                                    </button>

                                    <button
                                        type="submit"
                                        className="primary-button"
                                        disabled={creatingDoctor}
                                    >
                                        {creatingDoctor
                                            ? "กำลังบันทึก..."
                                            : "บันทึกข้อมูลแพทย์"}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {showDepartmentModal && (
                    <div
                        className="modal-overlay"
                        onClick={() => setShowDepartmentModal(false)}
                    >
                        <div
                            className="modal department-modal"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <div className="modal-header">
                                <div>
                                    <h2>เพิ่มข้อมูลแผนก</h2>
                                    <p>กรอกข้อมูลพื้นฐานของแผนก</p>
                                </div>

                                <button
                                    type="button"
                                    className="modal-close"
                                    onClick={() => setShowDepartmentModal(false)}
                                >
                                    ×
                                </button>
                            </div>

                            <form onSubmit={handleCreateDepartment}>
                                <div className="form-section">
                                    <div className="form-section-title">
                                        <h3>ข้อมูลแผนก</h3>
                                        <p>กรุณากรอกข้อมูลให้ครบถ้วน</p>
                                    </div>

                                    <div className="form-group">
                                        <label>
                                            ชื่อแผนกภาษาไทย <span>*</span>
                                        </label>

                                        <input
                                            type="text"
                                            placeholder="เช่น หัวใจ"
                                            value={departmentForm.th_name}
                                            onChange={(e) =>
                                                setDepartmentForm((prev) => ({
                                                    ...prev,
                                                    th_name: e.target.value,
                                                }))
                                            }
                                        />
                                    </div>

                                    <div className="form-group">
                                        <label>
                                            ชื่อแผนกภาษาอังกฤษ <span>*</span>
                                        </label>

                                        <input
                                            type="text"
                                            placeholder="เช่น Cardiology"
                                            value={departmentForm.eng_name}
                                            onChange={(e) =>
                                                setDepartmentForm((prev) => ({
                                                    ...prev,
                                                    eng_name: e.target.value,
                                                }))
                                            }
                                        />
                                    </div>
                                </div>

                                <div className="modal-actions">
                                    <button
                                        type="button"
                                        className="cancel-button"
                                        onClick={() => setShowDepartmentModal(false)}
                                    >
                                        ยกเลิก
                                    </button>

                                    <button
                                        type="submit"
                                        className="primary-button"
                                        disabled={creatingDepartment}
                                    >
                                        {creatingDepartment
                                            ? "กำลังบันทึก..."
                                            : "บันทึกข้อมูลแผนก"}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>

            {/* ================= รายการนัดหมาย ================= */}
            <header className="topbar">
                <div>
                    <h1>รายการนัดหมาย</h1>
                    <p>จัดการและสร้างนัดหมายสำหรับผู้ป่วย</p>
                </div>

                <button
                    className="primary-button"
                    onClick={() => setShowModal(true)}
                >
                    + สร้างนัดหมายใหม่
                </button>
            </header>

            {/* Search */}
            <section className="card">
                <div className="section-title">
                    <div>
                        <h2>ค้นหาช่วงเวลาที่ว่าง</h2>
                        <p>ค้นหาเวลานัดหมายที่แพทย์ยังว่าง</p>
                    </div>
                </div>

                <div className="filters">

                    <div className="form-group">
                        <label>แพทย์</label>

                        <div className="doctor-select">
                            <input
                                type="text"
                                placeholder="ค้นหาแพทย์..."
                                value={doctorSearch}
                                onFocus={() => setShowDoctorDropdown(true)}
                                onChange={(e) => {
                                    setDoctorSearch(e.target.value);
                                    setShowDoctorDropdown(true);
                                }}
                            />

                            {showDoctorDropdown && (
                                <div className="doctor-dropdown">
                                    {filteredDoctors.length > 0 ? (
                                        filteredDoctors.map((doctor) => (
                                            <div
                                                key={doctor.id}
                                                className="doctor-option"
                                                onClick={() => {
                                                    setAvailabilityForm((prev) => ({
                                                        ...prev,
                                                        doctor_id: doctor.id,
                                                    }));

                                                    setDoctorSearch(
                                                        `${doctor.title_name} ${doctor.first_name} ${doctor.last_name}`
                                                    );

                                                    setShowDoctorDropdown(false);
                                                }}
                                            >
                                                {doctor.title_name}{" "}
                                                {doctor.first_name}{" "}
                                                {doctor.last_name}
                                            </div>
                                        ))
                                    ) : (
                                        <div className="doctor-option empty">
                                            ไม่พบแพทย์
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>

                    <div className="form-group">
                        <label>วันที่นัดหมาย</label>

                        <input
                            type="date"
                            value={availabilityForm.date}
                            onChange={(e) =>
                                setAvailabilityForm((prev) => ({
                                    ...prev,
                                    date: e.target.value,
                                }))
                            }
                        />
                    </div>

                    <div className="form-group">
                        <label>เวลาเริ่ม</label>

                        <div className="time-input">
                            <input
                                type="text"
                                name="start_time"
                                placeholder="09:00"
                                maxLength="5"
                                inputMode="numeric"
                                onChange={handleAvailabilityTimeChange}
                            />
                            <span>น.</span>
                        </div>
                    </div>

                    <div className="form-group">
                        <label>เวลาสิ้นสุด</label>

                        <div className="time-input">
                            <input
                                type="text"
                                name="end_time"
                                placeholder="17:00"
                                maxLength="5"
                                inputMode="numeric"
                                onChange={handleAvailabilityTimeChange}
                            />
                            <span>น.</span>
                        </div>
                    </div>
                    <button
                        type="button"
                        className="search-button"
                        onClick={checkAvailability}
                        disabled={availabilityLoading}
                    >
                        {availabilityLoading
                            ? "กำลังตรวจสอบ..."
                            : "ค้นหาช่วงเวลาว่าง"}
                    </button>
                </div>

                {/* ========================= ผลการค้นหา ========================= */}
                {availabilityResult && (
                    <div className="availability-result">

                        <div className="availability-result-header">
                            <h3>ผลการตรวจสอบ</h3>
                        </div>

                        <div className="availability-result-content">

                            <div className="result-item">
                                <span>แพทย์</span>
                                <strong>
                                    {doctors.find(
                                        (doctor) =>
                                            doctor.id === availabilityResult.data?.doctor_id
                                    )?.title_name}{" "}
                                    {doctors.find(
                                        (doctor) =>
                                            doctor.id === availabilityResult.data?.doctor_id
                                    )?.first_name}{" "}
                                    {doctors.find(
                                        (doctor) =>
                                            doctor.id === availabilityResult.data?.doctor_id
                                    )?.last_name}
                                </strong>
                            </div>

                            <div className="result-item">
                                <span>วันที่</span>
                                <strong>
                                    {availabilityResult.data?.date || "-"}
                                </strong>
                            </div>

                            <div className="result-item">
                                <span>เวลา</span>
                                <strong>
                                    {availabilityResult.data?.start_time || "-"}
                                    {" - "}
                                    {availabilityResult.data?.end_time || "-"}
                                </strong>
                            </div>

                            <div className="result-item">
                                <span>สถานะ</span>

                                {availabilityResult.data?.available ? (
                                    <span className="status available">
                                        ว่าง
                                    </span>
                                ) : (
                                    <span className="status unavailable">
                                        ไม่ว่าง
                                    </span>
                                )}
                            </div>

                        </div>

                        <div
                            className={
                                availabilityResult.data?.available
                                    ? "availability-message success"
                                    : "availability-message error"
                            }
                        >
                            {availabilityResult.message_th}
                        </div>

                    </div>
                )}
            </section>

            {/* Create Appointment Modal */}
            {
                showModal && (
                    <div
                        className="modal-overlay"
                        onClick={() => setShowModal(false)}
                    >
                        <div
                            className="modal"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <div className="modal-header">
                                <div>
                                    <h2>สร้างนัดหมายใหม่</h2>
                                    <p>กรอกข้อมูลสำหรับสร้างนัดหมายผู้ป่วย</p>
                                </div>

                                <button
                                    className="modal-close"
                                    onClick={() => setShowModal(false)}
                                >
                                    ×
                                </button>
                            </div>

                            <div className="modal-body">

                                {/* ผู้ป่วย + แพทย์ */}
                                <div className="form-row">

                                    <div className="form-group">
                                        <label>ผู้ป่วย</label>

                                        <input
                                            type="text"
                                            value={patientSearch}
                                            onChange={(e) => searchPatients(e.target.value)}
                                            placeholder="ค้นหาด้วย รหัส หรือ ชื่อ-นามสกุล"
                                        />

                                        {patients.length > 0 && (
                                            <div className="patient-suggestions">
                                                {patients.map((patient) => (
                                                    <div
                                                        key={patient.id}
                                                        className="patient-option"
                                                        onClick={() => {
                                                            setForm({
                                                                ...form,
                                                                patient_id: patient.id,
                                                            });

                                                            setPatientSearch(
                                                                `${patient.patient_no} - ${patient.first_name} ${patient.last_name}`
                                                            );

                                                            setPatients([]);
                                                        }}
                                                    >
                                                        <strong>
                                                            {patient.first_name} {patient.last_name}
                                                        </strong>

                                                        <small>
                                                            รหัสประจำตัวผู้ป่วย: {patient.patient_no}
                                                        </small>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>

                                    <div className="form-group">
                                        <label>แพทย์</label>

                                        <select
                                            name="doctor_id"
                                            value={form.doctor_id}
                                            onChange={handleChange}
                                        >
                                            <option value="">เลือกแพทย์</option>

                                            {doctors.map((doctor) => (
                                                <option key={doctor.id} value={doctor.id}>
                                                    {doctor.title_name} {doctor.first_name}{" "}
                                                    {doctor.last_name}
                                                </option>
                                            ))}
                                        </select>
                                    </div>

                                </div>

                                {/* แผนก + ประเภท */}
                                <div className="form-row">

                                    <div className="form-group">
                                        <label>แผนก</label>

                                        <select
                                            name="department_id"
                                            value={form.department_id}
                                            onChange={handleChange}
                                        >
                                            <option value="">เลือกแผนก</option>

                                            {departments.map((department) => {
                                                console.log("Department ID:", department.id);

                                                return (
                                                    <option
                                                        key={department.id}
                                                        value={department.id}
                                                    >
                                                        {department.th_name} ({department.eng_name})
                                                    </option>
                                                );
                                            })}
                                        </select>
                                    </div>

                                    <div className="form-group">
                                        <label>ประเภทการนัดหมาย</label>

                                        <select
                                            name="appointment_type"
                                            value={form.appointment_type}
                                            onChange={handleChange}
                                        >
                                            <option value="">เลือกประเภท</option>
                                            <option value="new">ผู้ป่วยใหม่</option>
                                            <option value="follow-up">
                                                ผู้ป่วยติดตามผล
                                            </option>
                                            <option value="consultation">
                                                ปรึกษาแพทย์
                                            </option>
                                            <option value="procedure">
                                                หัตถการ
                                            </option>
                                        </select>
                                    </div>

                                </div>

                                {/* วันที่ + เวลา */}
                                <div className="form-row">

                                    <div className="form-group">
                                        <label>วันที่นัดหมาย</label>
                                        <input
                                            type="date"
                                            name="appointment_date"
                                            value={form.appointment_date}
                                            onChange={handleChange}
                                        />
                                    </div>

                                    <div className="form-group">
                                        <label>เวลาเริ่ม</label>
                                        <div className="time-input">
                                            <input
                                                type="text"
                                                name="start_time"
                                                value={form.start_time}
                                                onChange={handleTimeChange}
                                                placeholder="09:00"
                                                maxLength="5"
                                                inputMode="numeric"
                                            />
                                            <span>น.</span>
                                        </div>
                                    </div>

                                    <div className="form-group">
                                        <label>เวลาสิ้นสุด</label>
                                        <div className="time-input">
                                            <input
                                                type="text"
                                                name="end_time"
                                                value={form.end_time}
                                                onChange={handleTimeChange}
                                                placeholder="09:00"
                                                maxLength="5"
                                                inputMode="numeric"
                                            />
                                            <span>น.</span>
                                        </div>
                                    </div>

                                </div>

                                {/* เหตุผล */}
                                <div className="form-group">
                                    <label>เหตุผล / รายละเอียด</label>

                                    <textarea
                                        name="reason"
                                        value={form.reason}
                                        onChange={handleChange}
                                        placeholder="ระบุเหตุผลหรือรายละเอียดเพิ่มเติม"
                                        rows="3"
                                    />
                                </div>

                            </div>

                            <div className="modal-footer">

                                <button
                                    className="cancel-button"
                                    onClick={() => setShowModal(false)}
                                >
                                    ยกเลิก
                                </button>

                                <button
                                    className="primary-button"
                                    onClick={createAppointment}
                                    disabled={creating}
                                >
                                    {creating ? "กำลังบันทึก..." : "บันทึกนัดหมาย"}
                                </button>

                            </div>
                        </div>
                    </div>
                )
            }
            {/* ================= ตารางเวลาแพทย์ ================= */}
            <header className="topbar">
                <div>
                    <h1>ตารางเวลาของแพทย์</h1>
                    <p>จัดการตารางเวลาการทำงานของแพทย์</p>
                </div>

                <button
                    className="primary-button"
                    onClick={() => {
                        setScheduleForm({
                            doctor_id: "",
                            department_id: "",
                            day_of_week: "",
                            start_time: "",
                            end_time: "",
                            break_start: "",
                            break_end: "",
                            is_available: true,
                        });

                        setShowScheduleModal(true);
                    }}
                >
                    + เพิ่มตารางเวลา
                </button>
            </header>

            <section className="card">
                <div className="section-title">
                    <div>
                        <h2>ตารางเวลาการทำงาน</h2>
                        <p>ข้อมูลวันและเวลาทำงานของแพทย์</p>
                    </div>
                </div>

                {/* Search */}
                <div className="schedule-toolbar">
                    <input
                        type="text"
                        placeholder="ค้นหาชื่อแพทย์..."
                        value={scheduleSearch}
                        onChange={(e) => {
                            setScheduleSearch(e.target.value);
                            setSchedulePage(1);
                        }}
                    />
                </div>

                {/* Table */}
                <div className="schedule-table">
                    <div className="schedule-header">
                        <span>แพทย์</span>
                        <span>แผนก</span>
                        <span>วัน</span>
                        <span>เวลาทำงาน</span>
                        <span>เวลาพัก</span>
                        <span>สถานะ</span>
                        <span></span>
                    </div>

                    {schedules.length > 0 ? (
                        schedules.map((schedule) => (
                            <div
                                className="schedule-row"
                                key={schedule.id}
                            >
                                {/* แพทย์ */}
                                <span>
                                    <strong>
                                        {schedule.doctor?.title_name}{" "}
                                        {schedule.doctor?.first_name}{" "}
                                        {schedule.doctor?.last_name}
                                    </strong>
                                </span>

                                {/* แผนก */}
                                <span>
                                    {schedule.department?.th_name || "-"}
                                </span>

                                {/* วัน */}
                                <span>
                                    {dayNames[schedule.day_of_week] || "-"}
                                </span>

                                {/* เวลาทำงาน */}
                                <span>
                                    {formatTime(schedule.start_time)} - {formatTime(schedule.end_time)}
                                </span>

                                {/* เวลาพัก */}
                                <span>
                                    {schedule.break_start && schedule.break_end
                                        ? `${formatTime(schedule.break_start)} - ${formatTime(schedule.break_end)}`
                                        : "-"}
                                </span>

                                {/* สถานะ */}
                                <span>
                                    <b
                                        className={
                                            schedule.is_available
                                                ? "status booked"
                                                : "status cancelled"
                                        }
                                    >
                                        {schedule.is_available
                                            ? "พร้อมให้บริการ"
                                            : "ไม่พร้อมให้บริการ"}
                                    </b>
                                </span>

                                {/* ลบ */}
                                <span>
                                    <button
                                        className="delete-button"
                                        onClick={() => handleDeleteSchedule(schedule.id)}
                                    >
                                        ลบ
                                    </button>
                                </span>
                            </div>
                        ))
                    ) : (
                        <div className="empty-state">
                            ไม่พบข้อมูลตารางเวลาของแพทย์
                        </div>
                    )}

                    {showScheduleModal && (
                        <div
                            className="modal-overlay"
                            onClick={() => setShowScheduleModal(false)}
                        >
                            <div
                                className="modal schedule-modal"
                                onClick={(e) => e.stopPropagation()}
                            >
                                {/* Header */}
                                <div className="modal-header">
                                    <div>
                                        <h2>เพิ่มตารางเวลาของแพทย์</h2>
                                        <p>กำหนดวันและเวลาทำงานของแพทย์</p>
                                    </div>

                                    <button
                                        type="button"
                                        className="modal-close"
                                        onClick={() => setShowScheduleModal(false)}
                                    >
                                        ×
                                    </button>
                                </div>

                                <form onSubmit={handleCreateSchedule}>
                                    {/* Doctor / Department */}
                                    <div className="form-section">
                                        <div className="form-section-title">
                                            <h3>ข้อมูลการทำงาน</h3>
                                            <p>เลือกแพทย์และแผนกที่ต้องการกำหนดตาราง</p>
                                        </div>

                                        <div className="form-group">
                                            <label>
                                                แพทย์ <span>*</span>
                                            </label>

                                            <select
                                                value={scheduleForm.doctor_id}
                                                onChange={(e) =>
                                                    setScheduleForm((prev) => ({
                                                        ...prev,
                                                        doctor_id: e.target.value,
                                                    }))
                                                }
                                            >
                                                <option value="">เลือกแพทย์</option>

                                                {doctors.map((doctor) => (
                                                    <option
                                                        key={doctor.id}
                                                        value={doctor.id}
                                                    >
                                                        {doctor.title_name}{" "}
                                                        {doctor.first_name}{" "}
                                                        {doctor.last_name}
                                                    </option>
                                                ))}
                                            </select>
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label>
                                                    แผนก <span>*</span>
                                                </label>

                                                <select
                                                    value={scheduleForm.department_id}
                                                    onChange={(e) =>
                                                        setScheduleForm((prev) => ({
                                                            ...prev,
                                                            department_id: e.target.value,
                                                        }))
                                                    }
                                                >
                                                    <option value="">เลือกแผนก</option>

                                                    {departments.map((department) => (
                                                        <option
                                                            key={department.id}
                                                            value={department.id}
                                                        >
                                                            {department.th_name} ({department.eng_name})
                                                        </option>
                                                    ))}
                                                </select>
                                            </div>

                                            <div className="form-group">
                                                <label>
                                                    วันทำงาน <span>*</span>
                                                </label>

                                                <select
                                                    value={scheduleForm.day_of_week}
                                                    onChange={(e) =>
                                                        setScheduleForm((prev) => ({
                                                            ...prev,
                                                            day_of_week: e.target.value,
                                                        }))
                                                    }
                                                >
                                                    <option value="">เลือกวัน</option>
                                                    <option value="0">อาทิตย์</option>
                                                    <option value="1">จันทร์</option>
                                                    <option value="2">อังคาร</option>
                                                    <option value="3">พุธ</option>
                                                    <option value="4">พฤหัสบดี</option>
                                                    <option value="5">ศุกร์</option>
                                                    <option value="6">เสาร์</option>
                                                </select>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Working Time */}
                                    <div className="form-section">
                                        <div className="form-section-title">
                                            <h3>เวลาทำงาน</h3>
                                            <p>กำหนดช่วงเวลาที่แพทย์ให้บริการ</p>
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label>
                                                    เวลาเริ่มงาน <span>*</span>
                                                </label>

                                                <div className="time-input">
                                                    <input
                                                        type="text"
                                                        name="start_time"
                                                        placeholder="09:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(e, "start_time")
                                                        }
                                                    />
                                                    <span>น.</span>
                                                </div>
                                            </div>

                                            <div className="form-group">
                                                <label>
                                                    เวลาเลิกงาน <span>*</span>
                                                </label>

                                                <div className="time-input">
                                                    {/* <input
                                                        type="text"
                                                        name="end_time"
                                                        onChange={handleAvailabilityTimeChange}
                                                        placeholder="08:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                    /> */}
                                                    <input
                                                        type="text"
                                                        name="end_time"
                                                        placeholder="17:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(e, "end_time")
                                                        }
                                                    />
                                                    <span>น.</span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Break Time */}
                                    <div className="form-section">
                                        <div className="form-section-title">
                                            <h3>เวลาพัก</h3>
                                            <p>สามารถเว้นว่างได้หากไม่มีเวลาพัก</p>
                                        </div>

                                        <div className="form-row">
                                            <div className="form-group">
                                                <label>เริ่มพัก</label>

                                                <div className="time-input">
                                                    {/* <input
                                                        type="text"
                                                        placeholder="12:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        value={scheduleForm.break_start}
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(
                                                                e,
                                                                "break_start"
                                                            )
                                                        }
                                                    /> */}

                                                    <input
                                                        type="text"
                                                        name="break_start"
                                                        placeholder="12:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(e, "break_start")
                                                        }
                                                    />

                                                    <span>น.</span>
                                                </div>
                                            </div>

                                            <div className="form-group">
                                                <label>สิ้นสุดพัก</label>

                                                <div className="time-input">
                                                    {/* <input
                                                        type="text"
                                                        placeholder="13:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        value={scheduleForm.break_end}
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(
                                                                e,
                                                                "break_end"
                                                            )
                                                        }
                                                    /> */}
                                                    <input
                                                        type="text"
                                                        name="break_end"
                                                        placeholder="13:00"
                                                        maxLength="5"
                                                        inputMode="numeric"
                                                        onChange={(e) =>
                                                            handleScheduleTimeChange(e, "break_end")
                                                        }
                                                    />
                                                    <span>น.</span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Availability */}
                                    <div className="availability-toggle">
                                        <div>
                                            <strong>เปิดให้จอง</strong>
                                            <p>อนุญาตให้ผู้ป่วยสามารถนัดหมายในช่วงเวลานี้</p>
                                        </div>

                                        <label className="switch">
                                            <input
                                                type="checkbox"
                                                checked={scheduleForm.is_available}
                                                onChange={(e) =>
                                                    setScheduleForm((prev) => ({
                                                        ...prev,
                                                        is_available: e.target.checked,
                                                    }))
                                                }
                                            />

                                            <span className="slider"></span>
                                        </label>
                                    </div>

                                    {/* Actions */}
                                    <div className="modal-actions">
                                        <button
                                            type="button"
                                            className="cancel-button"
                                            onClick={() =>
                                                setShowScheduleModal(false)
                                            }
                                        >
                                            ยกเลิก
                                        </button>

                                        <button
                                            type="submit"
                                            className="primary-button"
                                            disabled={creatingSchedule}
                                        >
                                            {creatingSchedule
                                                ? "กำลังบันทึก..."
                                                : "บันทึกตารางเวลา"}
                                        </button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    )}
                </div>

                {/* Pagination */}
                <div className="pagination">
                    <span className="pagination-info">
                        แสดง{" "}
                        {scheduleTotal === 0
                            ? 0
                            : (schedulePage - 1) * schedulePageSize + 1}
                        -
                        {Math.min(
                            schedulePage * schedulePageSize,
                            scheduleTotal
                        )}{" "}
                        จาก {scheduleTotal} รายการ
                    </span>

                    <div className="pagination-controls">
                        <select className="page-size-select"
                            value={schedulePageSize}
                            onChange={(e) => {
                                setSchedulePageSize(Number(e.target.value));
                                setSchedulePage(1);
                            }}
                        >
                            <option value={10}>10 รายการ</option>
                            <option value={20}>20 รายการ</option>
                            <option value={50}>50 รายการ</option>
                        </select>


                        <button
                            disabled={schedulePage <= 1}
                            onClick={() =>
                                setSchedulePage((prev) => prev - 1)
                            }
                        >
                            ก่อนหน้า
                        </button>

                        <span>
                            หน้า {schedulePage} /{" "}
                            {Math.max(
                                1,
                                Math.ceil(scheduleTotal / schedulePageSize)
                            )}
                        </span>
                        <button
                            disabled={
                                schedulePage >=
                                Math.ceil(scheduleTotal / schedulePageSize)
                            }
                            onClick={() =>
                                setSchedulePage((prev) => prev + 1)
                            }
                        >
                            ถัดไป
                        </button>
                    </div>
                </div>
            </section>

            <section className="card">
                <div className="section-title">
                    <div>
                        <h2>กฎการจัดตารางเวลา</h2>
                        <p>
                            ระบบจะใช้ตารางเวลาเหล่านี้ในการคำนวณ
                            ช่วงเวลาที่สามารถนัดหมายได้
                        </p>
                    </div>
                </div>

                <div className="rules">
                    <div className="rule">
                        <strong>เวลาทำงาน</strong>
                        <span>
                            ไม่สามารถสร้างนัดหมายนอกเวลาทำงานของแพทย์ได้
                        </span>
                    </div>

                    <div className="rule">
                        <strong>เวลาพัก</strong>
                        <span>
                            ไม่สามารถสร้างนัดหมายในช่วงเวลาพักของแพทย์ได้
                        </span>
                    </div>

                    <div className="rule">
                        <strong>ตารางเวลาที่ไม่พร้อมให้บริการ</strong>
                        <span>
                            ตารางเวลาที่ถูกกำหนดให้ไม่พร้อมให้บริการ
                            จะไม่สามารถใช้สำหรับการสร้างนัดหมายได้
                        </span>
                    </div>
                </div>
            </section>

        </main >
    );
}



export default AppointmentPage;