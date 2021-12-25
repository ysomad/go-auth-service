package domain

// Device represents data transfer object with user device data
type Device struct {
	UserAgent string
	IP        string
}

func NewDevice(userAgent string, ip string) Device {
	return Device{
		UserAgent: userAgent,
		IP:        ip,
	}
}
