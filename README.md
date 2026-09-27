# H-Lab

ระบบ Hospital Appointment Booking สำหรับจัดการผู้ป่วย แพทย์ แผนก ตารางเวลา และการนัดหมาย

## Features

* จัดการข้อมูลผู้ป่วย
* จัดการข้อมูลแพทย์
* จัดการข้อมูลแผนก
* จัดการตารางเวลาของแพทย์
* ตรวจสอบเวลาว่างของแพทย์
* สร้างและจัดการนัดหมาย
* ตรวจสอบเวลานัดหมายซ้ำ
* รองรับ Docker

## Tech Stack

* Frontend: React
* Backend: Go + Fiber
* Database: MariaDB
* ORM: GORM
* Container: Docker / Docker Compose

## Installation & Run

### 1. Clone Repository

```bash
git clone <repository-url>
cd <project-folder>
```

### 2. Create `.env`

สร้างไฟล์ `.env` ที่ root ของโปรเจกต์

```env
DB_HOST=mariadb
DB_PORT=3306
DB_USER=hlabtest
DB_PASSWORD=hlab1234
DB_NAME=hlab-test
```

> ชื่อตัวแปรใน `.env` ต้องตรงกับที่ Backend ใช้งานจริง

### 3. Run with Docker

```bash
docker compose up --build
```

### 4. Access Application

Frontend:

```text
http://localhost:5173
```

Backend API:

```text
http://127.0.0.1:8080
```

### Stop Application

```bash
docker compose down
```
