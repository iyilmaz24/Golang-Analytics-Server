package types

import "time"

// devices := []Device{
// 	{Type: "Desktop", OS: "Windows", Browser: "Chrome"},
// 	{Type: "Mobile", OS: "iOS", Browser: "Safari"},
// }
type Device struct {
	Type    string
	OS     	string
	Browser string
}

// stat := Stat{
// 	Ip:            "192.168.1.1",
// 	Location:      "New York, USA",
// 	VDWebAppVisits: 10,
// 	FLPortalVisits: 5,
// 	NMPortalVisits: 3,
// 	Devices:       devices,
// 	FirstAccess:   time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
// 	LastAccess:    time.Now(),
// 	TotalVisits:   18,
// }
type UserStat struct {
	Ip  string
	Location string
	VD_WebApp int
	FL_Portal int
	NM_Portal int
	TotalVisits int
	Devices []Device
	FirstAccess time.Time
	LastAccess time.Time
}