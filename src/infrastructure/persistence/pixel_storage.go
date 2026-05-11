package persistence

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetPx(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND date(pxdate) = CURRENT_DATE AND pixel = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetToken(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND DATE(pxdate) = CURRENT_DATE AND token = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetByAdnetCode(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND date(pxdate) = CURRENT_DATE AND pixel LIKE ? AND is_used = false AND is_unique = false", o.URLServiceKey, "%"+o.Pixel+"%").
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetPxByMsisdn(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND DATE(pxdate) = CURRENT_DATE AND msisdn = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) UpdatePixelById(o entity.PixelStorage) error {

	result := r.DB.Model(&o).Select("is_used", "pixel_used_date", "status_postback", "status_postback", "is_unique", "url_postback", "status_url_postback", "reason_url_postback", "reason_url_postback").Updates(o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SpecialGetPx(o entity.PixelStorage, si int, ei int) (entity.PixelStorage, bool) {

	start := o.Pxdate
	var (
		end          time.Time
		result_query []string
		rpl          *strings.Replacer
		tbl          string
	)

	for x := si; x >= ei; x-- {

		SQL := `SELECT * FROM {TBL} WHERE url_service_key = '{CAMPAIGNID}' AND date(pxdate) = '{DATE}' AND pixel = '{PIXEL}' AND is_used = false AND is_unique = false`

		end = start.AddDate(0, 0, x)

		if x == 0 {
			tbl = "pixel_storages"
		} else {
			tbl = "pixel_storages" + "_" + end.Format("20060102")
		}

		rpl = strings.NewReplacer(
			"{TBL}", tbl,
			"{DATE}", end.Format("2006-01-02"),
			"{CAMPAIGNID}", o.URLServiceKey,
			"{PIXEL}", o.Pixel,
		)
		SQL = rpl.Replace(SQL)

		result_query = append(result_query, SQL)
	}

	result := r.DB.Raw(strings.Join(result_query, ` union `) + " LIMIT 1").Scan(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) UpdatePixelBilled(o entity.PixelStorage) error {

	result := r.DB.Exec(
		`UPDATE pixel_storages 
		 SET updated_at = NOW(), 
		     m_status_time_charge = NOW(), 
		     m_status_charge = ?
		 WHERE DATE(pxdate) = CURRENT_DATE 
		   AND url_service_key = ? 
		   AND pixel = ? 
		   AND is_unique = false `,
		o.MStatusCharge,
		o.URLServiceKey,
		o.Pixel,
	)

	if result.Error != nil {
		r.Logs.Error(fmt.Sprintf("DB error update pixel billed: %#v", result.Error))
		return result.Error
	}

	r.Logs.Debug(fmt.Sprintf("update pixel billed affected rows: %d", result.RowsAffected))

	return nil
}