package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) GetLastCampaignId(tbl string) int {

	var result int
	row := r.DB.Table(tbl).
		Select("COALESCE(MAX(id), 0) + 0").Row()
	row.Scan(&result)

	return result
}

func (r *BaseModel) GetCampaignByCampaignId(o entity.Campaign) (entity.Campaign, bool) {

	result := r.DB.Model(&o).
		Where("campaign_id = ?", o.CampaignId).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign id not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) NewCampaign(o entity.Campaign) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) NewCampaignDetail(o entity.CampaignDetail) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) ResetCappingCampaign(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).
		Where("is_active = ?", o.IsActive).
		Updates(entity.CampaignDetail{CounterMORatio: 0, StatusCapping: false})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetCampaignByCampaignDetailId(o entity.CampaignDetail) (entity.CampaignDetail, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign existed or data found %#v", o))
		return o, true
	}
}

func (r *BaseModel) UpdateCampaign(o entity.Campaign) error {

	result := r.DB.Model(&o).
		Where("campaign_id = ?", o.CampaignId).
		Updates(entity.Campaign{Name: o.Name, CampaignObjective: o.CampaignObjective, Country: o.Country, Advertiser: o.Advertiser})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec("UPDATE campaign_details SET is_active = ? WHERE id = ?", o.IsActive, o.ID)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SaveCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Save(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) DelCampaign(o entity.Campaign) error {

	result := r.DB.
		Unscoped().
		Where("campaign_id = ?", o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) CountCampaignDetailByCampaignID(o entity.CampaignDetail) (int, error) {
	var count int64
	err := r.DB.Model(&entity.CampaignDetail{}).Where("campaign_id = ?", o.CampaignId).Count(&count).Error
	return int(count), err
}

func (r *BaseModel) DelCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.
		Unscoped().
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) EditSettingCampaignDetail(o entity.CampaignDetail) error {
	updates := map[string]interface{}{
		"po":            o.PO,
		"mo_capping":    o.MOCapping,
		"ratio_send":    o.RatioSend,
		"ratio_receive": o.RatioReceive,
		"last_update":   o.LastUpdate,
	}

	if o.StatusCapping {
		updates["status_capping"] = false
	}

	result := r.DB.Model(&entity.CampaignDetail{}).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?",
			o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(updates)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateCampaignPO(o entity.CampaignDetail) error {
	/* return r.DB.Model(&entity.CampaignDetail{}).
	Where("url_service_key=? AND country=? AND operator=? AND partner=? AND service=? AND adnet=? AND campaign_id=?",
		o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
	Update("po", o.PO).Error */

	/* SQL := fmt.Sprintf(`UPDATE campaign_details SET po = '%s' WHERE url_service_key = '%s'`, o.PO, o.URLServiceKey)
	result := r.DB.Exec(SQL)

	r.Logs.Info(fmt.Sprintf("SQL : %s, affected: %d, is error : %#v", SQL, result.RowsAffected, result.Error)) */

	result := r.DB.Exec(`UPDATE campaign_details SET po = ? WHERE url_service_key = ?`, o.PO, o.URLServiceKey)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateCampaignRatio(o entity.CampaignDetail) error {
	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET ratio_send = ?, ratio_receive = ?
		WHERE url_service_key = ?
	`,
		o.RatioSend,
		o.RatioReceive,
		o.URLServiceKey,
	)

	r.Logs.Debug(fmt.Sprintf("UpdateCampaignRatio affected: %d, error: %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateCampaignMOCapping(o entity.CampaignDetail) error {
	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET mo_capping = ?, status_capping = ?, last_update = ?
		WHERE url_service_key = ?
	`,
		o.MOCapping,
		false,
		o.LastUpdate,
		o.URLServiceKey,
	)

	r.Logs.Debug(fmt.Sprintf("UpdateCampaignMOCapping affected: %d, error: %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateStatusCampaignDetail(o entity.CampaignDetail) error {
	result := r.DB.Model(&o).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?",
			o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.CampaignDetail{IsActive: o.IsActive})

	r.Logs.Debug(fmt.Sprintf("Query Affected: %d, Error: %v", result.RowsAffected, result.Error))

	if result.RowsAffected == 0 {
		r.Logs.Debug("No rows updated. Check if the WHERE condition matches any records.")
	}

	return result.Error
}

func (r *BaseModel) GetCampaignDetail(o entity.CampaignDetail) (entity.CampaignDetail, bool) {

	result := r.DB.Model(&o).First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign detail existed or data found %#v", o))
		return o, true
	}
}

func (r *BaseModel) CounterRatioById(o entity.CampaignDetail) error {

	result := r.DB.Exec("UPDATE campaign_details SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = ?", o.ID)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateStatusCounterById(o entity.CampaignDetail) error {

	//result := r.DB.Model(&o).Where("id = ?", o.ID).Updates(entity.CampaignDetail{CounterMOCapping: o.CounterMOCapping, StatusCapping: o.StatusCapping, CounterMORatio: o.CounterMORatio, StatusRatio: o.StatusRatio, LastUpdate: o.LastUpdate})

	//result := r.DB.Model(&o).Select("counter_mo_capping", "status_capping", "counter_mo_ratio", "status_ratio", "last_update").Updates(o)

	result := r.DB.Exec(fmt.Sprintf("UPDATE campaign_details SET counter_mo_capping = %d, status_capping = %t, counter_mo_ratio = %d, status_ratio = %t WHERE id = %d", o.CounterMOCapping, o.StatusCapping, o.CounterMORatio, o.StatusRatio, o.ID))

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetCampaignDetailByStatus(o entity.CampaignDetail, useStatus bool) ([]entity.CampaignDetail, error) {

	rows, _ := r.DB.Model(&entity.CampaignDetail{}).Where("is_active = ?", o.IsActive).Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignDetail
	)

	for rows.Next() {

		var s entity.CampaignDetail

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	return ss, rows.Err()
}

func (r *BaseModel) UpdateCPAReport(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).Updates(entity.CampaignDetail{CostPerConversion: o.CostPerConversion, AgencyFee: o.AgencyFee})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateCampaignMonitoringBudget(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).Where("country = ? AND operator = ?", o.Country, o.Operator).Updates(entity.CampaignDetail{TargetDailyBudget: o.TargetDailyBudget})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateKeyMainstreamCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET status_submit_key_mainstream = ?, key_mainstream = ? 
		WHERE url_service_key = ? AND campaign_id = ?`,
		o.StatusSubmitKeyMainstream, o.KeyMainstream, o.URLServiceKey, o.CampaignId,
	)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateGoogleSheetCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET google_sheet = ? 
		WHERE url_service_key = ? AND campaign_id = ?`,
		o.GoogleSheet, o.URLServiceKey, o.CampaignId,
	)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateGoogleSheetBillableCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET google_sheet_billable = ? 
		WHERE url_service_key = ? AND campaign_id = ?`,
		o.GoogleSheetBillable, o.URLServiceKey, o.CampaignId,
	)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetCampaignDetailByStatusAndCapped(o entity.CampaignDetail, useStatus bool) ([]entity.CampaignDetail, error) {

	rows, _ := r.DB.Model(&entity.CampaignDetail{}).Where("is_active = ? AND status_capping = true", o.IsActive).Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignDetail
	)

	for rows.Next() {

		var s entity.CampaignDetail

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	return ss, rows.Err()
}

func (r *BaseModel) ResetCappingCampaignByCapped(o entity.CampaignDetail) error {

	result := r.DB.Exec(fmt.Sprintf("UPDATE campaign_details SET counter_mo_capping = 0, counter_mo_capping_service = 0, status_capping = false, status_capping_service = false, counter_mo_ratio = 0, status_ratio = false AND is_active = true WHERE id = %d", o.ID))

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetUrlServiceKeyByService(o entity.CampaignDetail) ([]string, error) {
	var keys []string

	err := r.DB.
		Table("campaign_details").
		Select("url_service_key").
		Where(`
			country = ? AND
			operator = ? AND
			partner = ? AND
			service = ?
		`, o.Country, o.Operator, o.Partner, o.Service).
		Pluck("url_service_key", &keys).Error

	if err != nil {
		return nil, err
	}

	for i := range keys {
		keys[i] = keys[i] + "-configIdx"
	}

	return keys, nil
}

func (r *BaseModel) UpdateMOCappingS2S(o entity.CampaignDetail) error {
	result := r.DB.Exec(`
		UPDATE campaign_details cd
		SET mo_capping_service = ?, status_capping = false
		FROM campaigns c
		WHERE c.campaign_id = cd.campaign_id
			AND cd.country = ? 
			AND cd.operator = ? 
			AND cd.partner = ? 
			AND cd.service = ? 
			AND c.campaign_objective IN ('CPA', 'CPI', 'CPM', 'CPC', 'SINGLE URL S2S')
	`, o.MOCappingService, o.Country, o.Operator, o.Partner, o.Service)

	r.Logs.Debug(fmt.Sprintf("UpdateMOCounterService - affected: %d, error: %v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateMOCounterServiceS2S(o entity.CampaignDetail) error {
	var result *gorm.DB

	if o.CounterMOCappingService >= o.MOCappingService {
		result = r.DB.Exec(`
			UPDATE campaign_details cd
			SET counter_mo_capping_service = ?, status_capping = true
			FROM campaigns c
			WHERE c.campaign_id = cd.campaign_id
			AND cd.country = ? AND cd.operator = ? AND cd.partner = ? AND cd.service = ?
			AND c.campaign_objective IN ('CPA', 'CPI', 'CPM', 'CPC', 'SINGLE URL S2S')
		`, o.CounterMOCappingService, o.Country, o.Operator, o.Partner, o.Service)
	} else {
		result = r.DB.Exec(`
			UPDATE campaign_details cd
			SET counter_mo_capping_service = ?
			FROM campaigns c
			WHERE c.campaign_id = cd.campaign_id
			AND cd.country = ? AND cd.operator = ? AND cd.partner = ? AND cd.service = ?
			AND c.campaign_objective IN ('CPA', 'CPI', 'CPM', 'CPC', 'SINGLE URL S2S')
		`, o.CounterMOCappingService, o.Country, o.Operator, o.Partner, o.Service)
	}

	r.Logs.Debug(fmt.Sprintf("UpdateMOCounterService - affected: %d, error: %v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetMOCappingServiceAndCounter(o entity.CampaignDetail) (int, int, error) {
	var moCappingService int
	var counterMOCappingService int

	err := r.DB.Raw(`
		SELECT mo_capping_service, counter_mo_capping_service 
		FROM campaign_details 
		WHERE country = ? AND operator = ? AND partner = ? AND service = ?
		LIMIT 1
	`, o.Country, o.Operator, o.Partner, o.Service).Row().Scan(&moCappingService, &counterMOCappingService)

	if err != nil {
		r.Logs.Error(fmt.Sprintf("Failed to fetch MO Capping data: %v", err))
		return 0, 0, err
	}

	r.Logs.Debug(fmt.Sprintf("Fetched mo_capping_service: %d, counter_mo_capping_service: %d", moCappingService, counterMOCappingService))
	return moCappingService, counterMOCappingService, nil
}

func (r *BaseModel) GetCampaignDetailsWithObjectiveByCampaignId(campaignId string) ([]struct {
	entity.CampaignDetail
	CampaignObjective string
}, error) {
	var results []struct {
		entity.CampaignDetail
		CampaignObjective string
	}
	err := r.DB.Table("campaign_details").
		Select("campaign_details.*, campaigns.campaign_objective").
		Joins("JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id").
		Where("campaign_details.campaign_id = ?", campaignId).
		Scan(&results).Error
	return results, err
}

func (r *BaseModel) UpdateStatusCounterMOCampaignDetail(o entity.CampaignDetail) error {

	var result *gorm.DB

	if o.IsActive {

		var moReceived int64

		today := time.Now().Format("2006-01-02")

		err := r.DB.
			Table("inc_summary_campaigns").
			Select("COALESCE(mo_received, 0)").
			Where("url_service_key = ? AND summary_date = ?",
				o.URLServiceKey, today).
			Scan(&moReceived).Error

		if err != nil {
			return err
		}

		result = r.DB.Exec(`
			UPDATE campaign_details 
			SET is_active = ?, 
				status_capping = false,
			    counter_mo_capping = ?, 
			    counter_mo_capping_service = 0
			WHERE url_service_key = ?`,
			o.IsActive, moReceived, o.URLServiceKey,
		)

	} else {

		result = r.DB.Exec(`
			UPDATE campaign_details 
			SET is_active = ?
			WHERE url_service_key = ?`,
			o.IsActive, o.URLServiceKey,
		)
	}

	r.Logs.Debug(fmt.Sprintf("affected: %d, error: %#v",
		result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetAPICampaignDetails(country, operator, service, adnet string) ([]entity.InquiryAPICampIDResult, error) {
	var results []entity.InquiryAPICampIDResult
	err := r.DB.Raw(`
		SELECT
			cd.url_service_key,
			cd.country,
			cd.operator,
			cd.partner,
			cd.aggregator,
			cd.adnet,
			cd.service
		FROM campaign_details cd
		JOIN campaigns c ON c.campaign_id = cd.campaign_id
		WHERE cd.country = ?
		  AND cd.operator = ?
		  AND cd.service = ?
		  AND cd.adnet = ?
		  AND c.campaign_objective = 'API'
	`, country, operator, service, adnet).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}