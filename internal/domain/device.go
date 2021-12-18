package domain

import (
	"errors"
	"net"

	uaparser "github.com/mssola/user_agent"
)

var (
	// TODO: move to pkg httperror
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
		// TODO: return generic err
		return Device{}, ErrInvalidUserIP
	}

	ua := uaparser.New(userAgent)

	if ua.Bot() {
		// TODO: return generic err pkg/httperror
		return Device{}, ErrDeviceBot
	}

	return Device{
		UserAgent: userAgent,
		UserIP:    ip,
	}, nil
}
