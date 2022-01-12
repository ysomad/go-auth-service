package domain

// Device represents data transfer object with user device data
type Device struct {
	UserAgent string
	IP        string
}

func NewDevice(ua string, ip string) Device {
	return Device{
		UserAgent: ua,
		IP:        ip,
	}
}
