package models

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/iyilmaz24/Go-Analytics-Server/internal/database"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/helpers"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/database/types"
	geo "github.com/iyilmaz24/Go-Analytics-Server/internal/services"
)

type StatModel struct {
	DB *sql.DB
	Geo *geo.Geo
}

func (sm *StatModel) GetUserStats(id string) (*types.UserStat, error) {
	sqlQuery := database.GetUserStatsSQL()
	row := sm.DB.QueryRow(sqlQuery, id)

	var devicesJsonb []byte // stores raw JSONB data from Postgres
	s := &types.UserStat{}
	err := row.Scan(&s.Ip, &s.Location, &s.Region, &s.VD_WebApp, &s.FL_Portal, &s.NM_Portal, &s.TotalVisits, &devicesJsonb, &s.FirstAccess, &s.LastAccess)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	if len(devicesJsonb) > 0 {
		err = json.Unmarshal(devicesJsonb, &s.Devices) // unmarshal JSONB data from Postgres to Go struct
		if err != nil {
			s.Devices = []types.Device{}
		}
	} else {
		s.Devices = []types.Device{} // default if no devices/JSONB found in Postgres DB row
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

func (sm *StatModel) UpdateAppStats(s *types.AppStats) error {
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

	anonId := helpers.GetAnonymousID(s.Ip) // generate anonymous ID from IP address

	var location string
	region := s.Region
	newDevices := s.Devices
	var devicesJsonb []byte // JSONB for Postgres

	user, err := sm.GetUserStats(anonId) // check if user exists in DB

	if err != nil { 
		if errors.Is(err, ErrNoRecord) { // user does not exist
			if anonId != "invalid-ip" {
				location = sm.Geo.GetGeoLocation(s.Ip) // create new user location
			} else {
				location = "N/A"
			}
		} else { // other error
			return err
		}
	} else { // user exists
		location = user.Location // use existing location
		region = user.Region // use existing region

		newDevices = helpers.MergeDevices(s.Devices, user.Devices) // combine existing and new devices into []Device
	}

	devicesJsonb, err = json.Marshal(newDevices)
	if err != nil {
		devicesJsonb = []byte("[]") // default to empty JSONB if marshalling fails
	}

	sqlQuery := database.UpsertUserStatsSQL()
	_, err = sm.DB.Exec(sqlQuery, anonId, location, region, s.VD_WebApp, s.FL_Portal, s.NM_Portal, s.TotalVisits, devicesJsonb, s.FirstAccess, s.LastAccess)
	if err != nil {
		return err
	}
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

