package types

import "time"

type AppStats struct {
	HomeViews                     int
	HomeReportViews               int
	HomeToDashboardViews          int
	DashboardPageToDashboardViews int
	DashboardViews                int
	FaqsPageViews                 int
	GalleryPageViews              int
	TeamPageViews                 int
}

type AppStatsAggregate struct {
	Name      string
	BaseUrl   string
	Stats     AppStats
	CreatedAt time.Time
	UpdatedAt time.Time
}