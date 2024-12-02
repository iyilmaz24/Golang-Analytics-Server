package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/database"
)

type AppStats struct {
	HomeViews 		 int
	HomeReportViews     int
	HomeToDashboardViews int
	DashboardPageToDashboardViews int
	DashboardPageViews  int
	FaqsPageViews       int
	GalleryPageViews    int
	TeamPageViews       int
}

type AppStatsAggregate struct {          
	Name      string
	BaseUrl   string         
	Stats     AppStats
	CreatedAt time.Time
	UpdatedAt time.Time      
}

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

type StatModel struct {
	DB *sql.DB
}

type HealthCheck struct {
	Status          string `json:"status"`
	OpenConnections int    `json:"open_connections"`
	InUse           int    `json:"in_use"`
	Idle            int    `json:"idle"`
}

func (sm *StatModel) GetUserStats(id string) (*UserStat, error) {
	sqlQuery := database.GetUserStatsRowByKey()
	row := sm.DB.QueryRow(sqlQuery, id)

	s := &UserStat{}
	err := row.Scan(&s.Ip, &s.Location, &s.VD_WebApp, &s.FL_Portal, &s.NM_Portal, &s.TotalVisits, &s.Devices, &s.FirstAccess, &s.LastAccess)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	
	return s, nil
}

func (sm *StatModel) GetAppStats() (*AppStatsAggregate, error) {
	sqlQuery := database.GetAppStats()
	row := sm.DB.QueryRow(sqlQuery)

	s := &AppStatsAggregate{}
	err := row.Scan(&s.Name, &s.BaseUrl, &s.Stats, &s.CreatedAt, &s.UpdatedAt)

	// break down the Stats struct/JSON into individual fields

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	
	return s, nil
}

func (sm *AppStats) Upsert(s *AppStats) error {
	// add logic here
	return nil
}

func (sm *UserStat) Upsert(s *UserStat) error {
	// add logic here
	return nil
}

func (m *StatModel) CheckHealth() (*HealthCheck, error) {
	err := m.DB.Ping()
	status := "healthy"
	if err != nil {
		status = "unhealthy"
	}

	stats := m.DB.Stats()

	return &HealthCheck{
		Status:          status,
		OpenConnections: stats.OpenConnections,
		InUse:           stats.InUse,
		Idle:            stats.Idle,
	}, err
}

