package model

const (
	ExamStatusPrivate  = 0 // Chỉ mình user tạo mới thấy
	ExamStatusPending  = 1 // Đang chờ admin duyệt
	ExamStatusApproved = 2 // Đã được duyệt, public cho mọi người
	ExamStatusRejected = 3 // Bị từ chối bởi admin
)
