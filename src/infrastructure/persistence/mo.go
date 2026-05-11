package persistence

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) NewMO(o entity.MO) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) UpdateMO(o entity.MO) error {

	result := r.DB.Model(&o).
		Where("DATE(pxdate) = ? AND url_service_key = ? AND country = ? AND operator = ? AND service = ? AND adnet = ? AND pixel = ?", o.Pxdate, o.URLServiceKey, o.Country, o.Operator, o.Service, o.Adnet, o.Pixel).
		Updates(entity.MO{Msisdn: o.Msisdn, TrxId: o.TrxId, IsUsed: o.IsUsed, PixelUsedDate: o.PixelUsedDate, StatusPostback: o.StatusPostback, IsUnique: o.IsUnique, URLPostback: o.URLPostback, StatusURLPostback: o.StatusURLPostback, ReasonURLPostback: o.ReasonURLPostback})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}
