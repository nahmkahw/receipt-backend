// STATUS 4xx (Client Error) 5xx (Server Error)
package errs

import "net/http"

type ErrorHandler struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ทำ receiver function conform ตาม error interface ของ godror
// ซึ่งจะทำให้เรายังสามาถร retrun error ได้ตามปกติเพราะ error ที่เราสร้างมานี้มีโครงสร้างเดียวกันกับ Error ของ go
func (e ErrorHandler) Error() string {
	return e.Message
}

// สร้าง Code และ Message ใช้เอง 
func NewMessageAndStatusCode(status_code int, message string) error {
	return ErrorHandler{
		Code: status_code,
		Message: message,
	}
}

// STATUS 200 OK — การส่ง request สำเร็จ
func NewSuccessStatus() error {
	return ErrorHandler{
		Code: http.StatusOK,
		Message: "กระบวนการที่ร้องขอสำเร็จ",
	}
}

// STATUS 201 Created — client create ข้อมูลลง data base สำเร็จ
func NewCreatedStatus() error {
	return ErrorHandler{
		Code: http.StatusCreated,
		Message: "กระบวนการที่ร้องขอสำเร็จ",
	}
}

// STATUS 204: No Content — server ประมวลผลเรียบร้อยแล้ว แต่ไม่มีเนื้อหาส่งคืน
func NewNoContentStatus() error {
	return ErrorHandler{
		Code: http.StatusNoContent,
		Message: "ไม่พบข้อมูล",
	}
}

// STATUS 400 Bad Request — client ส่ง body request ไม่ถูกต้อง
func NewBadRequestError() error {
	return ErrorHandler{
		Code: http.StatusBadRequest,
		Message: "ส่ง body request ไม่ถูกต้อง",
	}
}

// STATUS 401 Unauthorized — client ยัง ไม่ได้ระบุตัวตน หรือไม่มี header (สำหรับยังไม่ได้ Login)
func NewUnauthorizedError() error {
	return ErrorHandler{
		Code: http.StatusUnauthorized,
		Message: "กรุณา Sign-in เพื่อเข้าสู่ระบบ",
	}
}

// STATUS 402 Payment Required — มีการเรียกชำระเงิน (ใช้ในอนาคต)
func NewPaymentRequiredError() error {
	return ErrorHandler{
		Code: http.StatusPaymentRequired,
		Message: "กรุณาทำการชำระเงิน",
	}
}

// STATUS 403 Forbidden — client ระบุตัวตนแล้วแต่ไม่มีสิทธิ์เข้าถึงส่วนนี้
func NewForbiddenError() error {
	return ErrorHandler{
		Code: http.StatusForbidden,
		Message: "คุณไม่ได้รับสิทธิ์เข้าถึงส่วนนี้",
	}
}

// STATUS 404 Not Found — ไม่พบหน้าที่ร้องขอ
func NewNotFoundError() error {
	return ErrorHandler{
		Code: http.StatusNotFound,
		Message: "ไม่พบข้อมูลที่ร้องขอ",
	}
}

// STATUS 422 Unprocessable Entity -เกิดจาก request ที่ส่งเข้ามามีรูปแบบที่ถูกต้อง แต่ข้อมูลที่ส่งเข้ามาไม่ถูกต้อง หรือไม่ครบตามความต้องการของเซิฟเวอร์
func NewUnprocessableEntityError() error {
	return ErrorHandler{
		Code: http.StatusUnprocessableEntity,
		Message: "ข้อมูลไม่ครบถ้วน กรุณากรอกข้อมูลให้ครบถ้วน",
	}
}

// STATUS 500 มีข้อผิดพลาดบางอย่างภายใน server โดยไม่ทราบสาเหตุ
func NewInternalServerError() error {
	return ErrorHandler{
		Code: http.StatusInternalServerError,
		Message: "ขณะนี้มีผู้เข้าใช้งานระบบเป็นจำนวนมาก กรุณารอสักครู่ ขออภัยในความไม่สะดวก",
	}
}

// STATUS 501 Not Implemented -server ไม่เข้าใจ request หรือไม่สามารถทำงานตามคำสั่งได้
func NewNotImplementedError() error {
	return ErrorHandler{
		Code: http.StatusNotImplemented,
		Message: "ขณะนี้มีผู้เข้าใช้งานระบบเป็นจำนวนมาก กรุณารอสักครู่ ขออภัยในความไม่สะดวก",
	}
}

// STATUS 502 Bad Gateway —server เป็น Gateway หรือ Proxy ได้รับ response ผิด
func NewBadGatewayError() error {
	return ErrorHandler{
		Code: http.StatusBadGateway,
		Message: "Gateway หรือ Proxy ได้รับ response ไม่ถูกต้อง",
	}
}

// STATUS 503 Service Unavailable — ใช้งานเกินพิกัด(ล่ม) หรือกำลังปรับปรุง server
func NewServiceUnavailableError() error {
	return ErrorHandler{
		Code: http.StatusServiceUnavailable,
		Message: "ขณะนี้มีผู้เข้าใช้งานระบบเป็นจำนวนมาก กรุณารอสักครู่ ขออภัยในความไม่สะดวก",
	}
}

// STATUS 504 Gateway Timeout — server ไม่ได้รับตอบสนองจาก server อื่น จนหมดเวลากันก่อน
func NewGatewayTimeoutError() error {
	return ErrorHandler{
		Code: http.StatusGatewayTimeout,
		Message: "ขณะนี้มีผู้เข้าใช้งานระบบเป็นจำนวนมาก กรุณารอสักครู่ ขออภัยในความไม่สะดวก",
	}
}