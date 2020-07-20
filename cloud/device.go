package cloud

import (
	"context"
	"time"
)

// Device represents device data.
type Device struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	TemperatureOffset int32        `json:"temperature_offset"`
	HumidityOffset    int32        `json:"humidity_offset"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
	FirmwareVersion   string       `json:"firmware_version"`
	MacAddress        string       `json:"mac_address"`
	SerialNumber      string       `json:"serial_number"`
	NewestEvents      NewestEvents `json:"newest_events"`
}

// NewestEvents represents event data.
type NewestEvents struct {
	Temperature SensorValue `json:"te"`
	Humidity    SensorValue `json:"hu"`
	Illuminance SensorValue `json:"il"`
	Motion      SensorValue `json:"mo"`
}

// SensorValue represents sensor data.
type SensorValue struct {
	Value     float64   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}

// Devices provides interface of /devices end-point.
type Devices interface {
	GetDevices(ctx context.Context) ([]*Device, error)
}

type devices struct {
	cli *Client
}

// GetDevices provides implementation of GET /devices
func (api *devices) GetDevices(ctx context.Context) ([]*Device, error) {
	var d []*Device
	if err := api.cli.Get(ctx, "devices", nil, &d); err != nil {
		return nil, err
	}
	return d, nil
}
