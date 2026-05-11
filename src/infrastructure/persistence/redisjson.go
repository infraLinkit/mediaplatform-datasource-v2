package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/redis/rueidis"
)

func (h *BaseModel) GetDataConfig(key string, path string) (*entity.DataConfig, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.DataConfig
		tempCfg [][]*entity.DataConfig // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {

			if len(tempCfg[0]) > 0 {
				h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
				dcfg = tempCfg[0][0]
				return dcfg, nil
			} else {
				err = errors.New("key is empty or not found")
				h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
				return dcfg, err
			}

		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) GetDataConfigCounter(key string, path string) (*entity.DataCounter, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.DataCounter
		tempCfg [][]*entity.DataCounter // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data counter key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) IncrCounterData(key string, path string, val float64) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonNumincrby().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Increament error counter key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Increament success counter key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) AppendCounterData(key string, path string, o entity.DataCounterDetail) {

	// Get Config Data Landing
	ctx := context.Background()

	b, _ := json.Marshal(o)
	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonArrappend().Key(key).Path(path).Value(string(b)).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) AppendData(key string, path string, o []byte) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonArrappend().Key(key).Path(path).Value(string(o)).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetCounterData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) IndexRedis(key string, field string) {

	h.R.Conn().B().FtCreate().Index(key).Prefix(1)

	//`$.campaign_id AS campaign_id TEXT $.pixel AS pixel TEXT $.user_agent AS user_agent TEXT $.os AS os TEXT $.browser AS browser TEXT $.ips AS ips TEXT $.user_is_rejected AS user_is_rejected TAG $.user_is_duplicated AS user_is_duplicated TAG $.refferal_url AS refferal_url TEXT $.handset_code AS handset_code TEXT $.handset_type AS handset_type TEXT $.pixel_is_used AS pixel_is_used TAG`

	//Create Indexing key
	result := h.R.Conn().Do(context.Background(), h.R.Conn().B().FtCreate().Index(key).OnJson().Schema().FieldName(field).Text().Build())
	isIdxCreated, err := result.AsBool()

	if err != nil {
		h.Logs.Info(fmt.Sprintf("[v] Created FT idx key ( %s ) : %t, %#v...\n\r", key, isIdxCreated, result))
	} else {

		if !isIdxCreated {
			h.Logs.Info(fmt.Sprintf("[x] Failed created FT idx key ( %s ) : %#v...\n\r", key, isIdxCreated))
		} else {
			h.Logs.Info(fmt.Sprintf("[v] Success created FT idx key ( %s ) : %t, %#v...\n\r", key, isIdxCreated, result))
		}

	}

}

func (h *BaseModel) GetData(i []interface{}, key string, path string) []interface{} {

	// Get Config Data Landing
	ctx := context.Background()

	rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonGet().Key(key).Path(path).Build()), &i)

	return i
}

func (h *BaseModel) SetData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetExpireData(key string, expr int64) {

	//Set key expire and delete automatically
	result := h.R.Conn().Do(context.Background(), h.R.Conn().B().Expire().Key(key).Seconds(expr).Build())

	isExpiredCreated, err := result.AsBool()

	if err != nil {
		h.Logs.Info(fmt.Sprintf("[x] Error Set duration expire key ( %s ) : %t, %#v...\n\r", key, isExpiredCreated, result))
	} else {

		if !isExpiredCreated {
			h.Logs.Info(fmt.Sprintf("[x] Failed Set duration expire key ( %s ) : %#v...\n\r", key, isExpiredCreated))
		} else {
			h.Logs.Info(fmt.Sprintf("[v] Success Set duration expire key ( %s ) : %t, %#v...\n\r", key, isExpiredCreated, result))
		}

	}
}

func (h *BaseModel) DelData(key string, path string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonDel().Key(key).Path(path).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Del data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Del data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) GetAlertData(key string, path string) (*entity.AlertData, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.AlertData
		tempCfg [][]*entity.AlertData // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data counter key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) GetDataSummary(key string, path string) (*entity.Summary, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.Summary
		tempCfg [][]*entity.Summary // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) RGetApiPinReport(key string, path string) ([]entity.ApiPinReport, bool) {

	var (
		isEmpty bool
		p       []entity.ApiPinReport
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var pinreport [][]entity.ApiPinReport
		v.DecodeJSON(&pinreport)

		if len(pinreport) > 0 {
			isEmpty = false
			p = pinreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetApiPinPerformanceReport(key string, path string) ([]entity.ApiPinPerformance, bool) {

	var (
		isEmpty bool
		p       []entity.ApiPinPerformance
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var pinperformancereport [][]entity.ApiPinPerformance
		v.DecodeJSON(&pinperformancereport)

		if len(pinperformancereport) > 0 {
			isEmpty = false
			p = pinperformancereport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetDisplayCPAReport(key string, path string) ([]entity.SummaryCampaign, bool) {

	var (
		isEmpty bool
		p       []entity.SummaryCampaign
	)

	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var displaycpareport [][]entity.SummaryCampaign
		v.DecodeJSON(&displaycpareport)

		if len(displaycpareport) > 0 {
			isEmpty = false
			p = displaycpareport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	fmt.Println("---- Query from Redis / not from DB ----")

	return p, isEmpty
}

func (h *BaseModel) RGetConversionLogReport(key string, path string) ([]entity.PixelStorage, bool) {

	var (
		isEmpty bool
		p       []entity.PixelStorage
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var conversionLogReport [][]entity.PixelStorage
		v.DecodeJSON(&conversionLogReport)

		if len(conversionLogReport) > 0 {
			isEmpty = false
			p = conversionLogReport[0]

			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetArpuReport(key string, path string) (entity.ARPUResponse, bool) {
	var (
		isEmpty = true
		p       entity.ARPUResponse
	)

	ctx := context.Background()

	data, err := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")
	if err != nil {
		h.Logs.Error(fmt.Sprintf("Failed to get key from Redis: %v", err))
		return p, true
	}

	for _, v := range data {
		raw, err := v.ToString()
		if err != nil || raw == "" || raw == "null" {
			h.Logs.Debug(fmt.Sprintf("Redis key (%s) is null or empty", key))
			continue
		}

		var arpuArray []entity.ARPUResponse
		if err := json.Unmarshal([]byte(raw), &arpuArray); err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to decode ARPU array for key (%s): %v", key, err))
			continue
		}

		if len(arpuArray) > 0 {
			p = arpuArray[0]
			isEmpty = false

			count := 0
			if p.Data != nil {
				count = len(p.Data.Data)
			}
			h.Logs.Debug(fmt.Sprintf("Parsed JSON key (%s), total items: %d", key, count))
		}
	}
	if !isEmpty {
		fmt.Println("---- Query from Redis / not from DB ----")
	}

	return p, isEmpty
}

func (h *BaseModel) RGetDisplayCostReport(key string, path string) ([]entity.CostReport, bool) {
	var (
		isEmpty bool
		p       []entity.CostReport
	)
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var costreport [][]entity.CostReport
		v.DecodeJSON(&costreport)

		if len(costreport) > 0 {
			isEmpty = false
			p = costreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}
	return p, isEmpty
}

func (h *BaseModel) RGetDisplayCostReportDetail(key string, path string) ([]entity.CostReport, bool) {
	var (
		isEmpty bool
		p       []entity.CostReport
	)
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var displaycostreport [][]entity.CostReport
		v.DecodeJSON(&displaycostreport)

		if len(displaycostreport) > 0 {
			isEmpty = false
			p = displaycostreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetCampaignManagement(key string, path string) ([]entity.CampaignManagementData, bool) {

	var (
		isEmpty bool
		p       []entity.CampaignManagementData
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var campaignmanagement [][]entity.CampaignManagementData
		v.DecodeJSON(&campaignmanagement)

		if len(campaignmanagement) > 0 {
			isEmpty = false
			p = campaignmanagement[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetMenu(key string, path string) ([]entity.Menu, bool) {

	var (
		isEmpty bool
		p       []entity.Menu
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var menu [][]entity.Menu
		v.DecodeJSON(&menu)

		if len(menu) > 0 {
			isEmpty = false
			p = menu[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetAlertReportAll(key string, path string) ([]entity.SummaryAll, bool) {
	var (
		isEmpty bool
		p       []entity.SummaryAll
	)

	ctx := context.Background()
	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var summaryMo [][]entity.SummaryAll
		v.DecodeJSON(&summaryMo)
		if len(summaryMo) > 0 {
			isEmpty = false
			p = summaryMo[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}
	return p, isEmpty
}

func (h *BaseModel) RGetRole(key string, path string) ([]entity.RoleManagementData, bool) {

	var (
		isEmpty bool
		p       []entity.RoleManagementData
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var role [][]entity.RoleManagementData
		v.DecodeJSON(&role)

		if len(role) > 0 {
			isEmpty = false
			p = role[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetUser(key string, path string) ([]entity.UserManagementData, bool) {

	var (
		isEmpty bool
		p       []entity.UserManagementData
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var user [][]entity.UserManagementData
		v.DecodeJSON(&user)

		if len(user) > 0 {
			isEmpty = false
			p = user[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetUserApprovalRequest(key string, path string) ([]entity.UserApprovalRequestData, bool) {

	var (
		isEmpty bool
		p       []entity.UserApprovalRequestData
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var user [][]entity.UserApprovalRequestData
		v.DecodeJSON(&user)

		if len(user) > 0 {
			isEmpty = false
			p = user[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetDisplayMainstreamReport(key string, path string) ([]entity.SummaryCampaign, bool) {
	var (
		isempty bool
		p       []entity.SummaryCampaign
	)

	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var displaymainstreamreport [][]entity.SummaryCampaign
		v.DecodeJSON(&displaymainstreamreport)

		if len(displaymainstreamreport) > 0 {
			isempty = false
			p = displaymainstreamreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isempty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	fmt.Println("---- Query from redis ----")

	return p, isempty
}

func (h *BaseModel) ScanKeys(pattern string) ([]string, error) {
	ctx := context.Background()
	var (
		cursor uint64 = 0
		keys   []string
	)
	for {
		result, err := h.R.Conn().Do(ctx, h.R.Conn().B().Scan().Cursor(cursor).Match(pattern).Count(100).Build()).AsScanEntry()
		if err != nil {
			h.Logs.Error(fmt.Sprintf("ScanKeys error: %v", err))
			return nil, err
		}
		keys = append(keys, result.Elements...)
		cursor = result.Cursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (h *BaseModel) GetDataJSON(key string) (map[string]interface{}, error) {
	ctx := context.Background()
	result, err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonGet().Key(key).Path("$").Build()).ToString()
	if err != nil {
		h.Logs.Error(fmt.Sprintf("GetDataJSON error: %v", err))
		return nil, err
	}
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		h.Logs.Error(fmt.Sprintf("GetDataJSON unmarshal error: %v", err))
		return nil, err
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return nil, errors.New("no data found or empty JSON")
}

func (h *BaseModel) GetDomainServices(key string, path string) ([]entity.DomainServices, error) {
	ctx := context.Background()

	var tempCfg [][]map[string][]entity.DomainServices
	var result []entity.DomainServices

	if err := rueidis.DecodeSliceOfJSON(
		h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()),
		&tempCfg,
	); err != nil {
		h.Logs.Warn(fmt.Sprintf("Key (%s) not found, creating empty list ...", key))
		h.R.Conn().Do(ctx, h.R.Conn().B().JsonSet().Key(key).Path("$").Value(`{"data":[]}`).Build())
		return result, nil
	}

	if len(tempCfg) > 0 && len(tempCfg[0]) > 0 {
		result = tempCfg[0][0]["data"]
	}

	if result == nil {
		result = []entity.DomainServices{}
	}

	return result, nil
}
