package entity

import (
	"errors"
	"net"

	uaparser "github.com/mssola/user_agent"
)

var (
	ErrDeviceBot     = errors.New("device is a bot")
	ErrInvalidUserIP = errors.New("invalid user ip")
)

// Device represents data transfer object with user device data
type Device struct {
	UserAgent string
	UserIP    string
}

func NewDevice(userAgent string, ip string) (Device, error) {
	if net.ParseIP(ip) == nil {
		return Device{}, ErrInvalidUserIP
	}

	ua := uaparser.New(userAgent)

	if ua.Bot() {
		return Device{}, ErrDeviceBot
	}

	return Device{
		UserAgent: ua.UA(),
		UserIP:    ip,
	}, nil
}
