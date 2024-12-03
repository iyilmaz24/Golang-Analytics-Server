package models

import (
	"database/sql"
	"errors"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/database"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/types"
)

type StatModel struct {
	DB *sql.DB
}

func (sm *StatModel) GetUserStats(id string) (*types.UserStat, error) {
	sqlQuery := database.GetUserStatsSQL()
	row := sm.DB.QueryRow(sqlQuery, id)

	s := &types.UserStat{}
	err := row.Scan(&s.Ip, &s.Location, &s.VD_WebApp, &s.FL_Portal, &s.NM_Portal, &s.TotalVisits, &s.Devices, &s.FirstAccess, &s.LastAccess)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	
	return s, nil
}

func (sm *StatModel) GetAppStats() (*types.AppStatsAggregate, error) {
	sqlQuery := database.GetAppStatsSQL()
	row := sm.DB.QueryRow(sqlQuery)

	s := &types.AppStatsAggregate{}
	err := row.Scan(&s.Name, &s.BaseUrl, &s.Stats, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

func (sm *StatModel) UpsertAppStats(s *types.AppStats) error {
	AppStatsAggregate, err := sm.GetAppStats()
	if err != nil {
		return err
	}
	statsJson := AppStatsAggregate.Stats

	statsJson.HomeViews += s.HomeViews
	statsJson.HomeReportViews += s.HomeReportViews
	statsJson.HomeToDashboardViews += s.HomeToDashboardViews
	statsJson.DashboardPageToDashboardViews += s.DashboardPageToDashboardViews
	statsJson.DashboardViews += s.DashboardViews
	statsJson.FaqsPageViews += s.FaqsPageViews
	statsJson.GalleryPageViews += s.GalleryPageViews
	statsJson.TeamPageViews += s.TeamPageViews

	sqlQuery := database.UpdateAppStatsSQL()
	_, err = sm.DB.Exec(sqlQuery, statsJson)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StatModel) UpsertUserStats(s *types.UserStat) error {
	// add logic here
	return nil
}

func (m *StatModel) CheckHealth() (*types.HealthCheck, error) {
	err := m.DB.Ping()
	status := "healthy"
	if err != nil {
		status = "unhealthy"
	}

	stats := m.DB.Stats()

	return &types.HealthCheck{
		Status:          status,
		OpenConnections: stats.OpenConnections,
		InUse:           stats.InUse,
		Idle:            stats.Idle,
	}, err
}

