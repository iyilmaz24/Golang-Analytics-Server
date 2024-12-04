package helpers

import (
	"fmt"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/types"
)

func MergeDevices(old, new []types.Device) []types.Device {
	deviceMap := make(map[string]types.Device)

	for _, device := range old {
		key := fmt.Sprintf("%s|%s|%s", device.Type, device.OS, device.Browser)
		deviceMap[key] = device
	}

	for _, device := range new {
		key := fmt.Sprintf("%s|%s|%s", device.Type, device.OS, device.Browser)
		deviceMap[key] = device
	}

	mergedDevices := make([]types.Device, 0, len(deviceMap))
	for _, device := range deviceMap {
		mergedDevices = append(mergedDevices, device)
	}

	return mergedDevices
}