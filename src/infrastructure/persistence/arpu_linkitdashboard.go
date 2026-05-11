package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"encoding/base64"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

// GetDataArpu method untuk BaseModel
func (r *BaseModel) GetDataArpu(fe entity.ArpuParams) (result entity.ARPUResponse, err error) {
	// Menggunakan environment variable APIARPU dari config
	baseURL := r.Config.APIARPU
	if baseURL == "" {
		return entity.ARPUResponse{}, errors.New("APIARPU environment variable is empty")
	}
	fmt.Println("APIARPU: ", baseURL)

	// Membuat URL untuk request API
	apiURL, err := url.Parse(baseURL + "/api/v4/arpu/arpu90")
	if err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Menambahkan query parameters
	query := url.Values{}
	if fe.From != "" {
		query.Set("from", fe.From)
	}
	if fe.To != "" {
		query.Set("to", fe.To)
	}
	if fe.Country != "" {
		query.Set("country", fe.Country)
	}
	if fe.Operator != "" {
		query.Set("operator", fe.Operator)
	}
	if fe.Service != "" {
		query.Set("service", fe.Service)
	}

	apiURL.RawQuery = query.Encode()

	// Membuat HTTP request
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	encUsername := r.Config.ARPUUsername
	encPassword := r.Config.ARPUPassword

	fmt.Println("ARPUUsername: ", encUsername, "ARPUPassword: ", encPassword)

	username, err := decryptEnv(encUsername)
	if err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to decrypt ARPUUsername: %v", err)
	}
	password, err := decryptEnv(encPassword)
	if err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to decrypt ARPUPassword: %v", err)
	}
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	fmt.Println("auth: ", auth, username, password)

	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Accept", "application/json")

	// Melakukan request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Membaca response
	var apiResponse entity.ARPUResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return entity.ARPUResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}

	// Validasi response
	if apiResponse.Status != 200 {
		return entity.ARPUResponse{}, fmt.Errorf("API returned status %d: %s", apiResponse.Status, apiResponse.Message)
	}

	r.Logs.Info(fmt.Sprintf("Successfully retrieved ARPU data for %s/%s/%s", fe.Country, fe.Operator, fe.Service))
	return apiResponse, nil
}

func decryptEnv(enc string) (string, error) {
	// contoh: jika hanya base64
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *BaseModel) SendWakiCallback(curdate_time time.Time) error {
	var summaries []entity.SummaryCampaign

	baseURL := r.Config.APILINKITDashboard
	if baseURL == "" {
		return errors.New("APILINKITDashboard environment variable is empty")
	}

	// Ambil semua summary_campaigns untuk hari ini yg mo_received > 0
	if err := r.DB.
		Where("summary_date = ? AND mo_received > 0 AND deleted_at IS NULL", curdate_time).
		Find(&summaries).Error; err != nil {
		return err
	}

	for _, sc := range summaries {
		// Bangun query URL
		q := url.Values{
			"date":           {curdate_time.Format("2006-01-02")},
			"campaign_id":    {sc.URLServiceKey},
			"publisher":      {sc.Adnet},
			"adnet":          {sc.Adnet},
			"operator":       {sc.Partner},
			"adn":            {sc.ShortCode},
			"client":         {sc.Partner},
			"aggregator":     {sc.Aggregator},
			"country":        {sc.Country},
			"service":        {sc.Service},
			"mo_received":    {strconv.Itoa(sc.MoReceived)},
			"mo_postback":    {strconv.Itoa(sc.Postback)},
			"total_mo":       {strconv.Itoa(sc.MoReceived)},
			"total_postback": {strconv.Itoa(sc.Postback)},
			"landing":        {strconv.Itoa(sc.Traffic)},
			"cr_mo_received": {strconv.FormatFloat(sc.CrMO, 'f', 2, 64)},
			"cr_mo_postback": {strconv.FormatFloat(sc.CrPostback, 'f', 2, 64)},
			"url_campaign":   {sc.URLAfter},
			"url_service":    {sc.URLBefore},
			"sbaf":           {strconv.FormatFloat(sc.SBAF, 'f', 2, 64)},
			"saaf":           {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"spending":       {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"campaign":       {sc.CampaignObjective},
			"payout":         {strconv.FormatFloat(sc.PO, 'f', 2, 64)},
			"price_per_mo":   {strconv.FormatFloat(sc.PricePerMO, 'f', 2, 64)},
		}

		// Gabungkan URL dan query param
		fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())

		var strlog string

		external.Block{
			Try: func() {

				if _, _, _, strlog, err := external.Get(fullURL, map[string]string{}, external.PHttp{
					Timeout:            5,
					KeepAlive:          3,
					IsDisableKeepAlive: true,
					MaxIdleConns:       5,
					IdleConnTimeout:    5,
					DisableCompression: false,
				}); err == nil {
					r.Logs.Info(strlog)
				} else {
					r.Logs.Error(strlog)
				}
			},
			Catch: func(e external.Exception) {

				r.Logs.Error(
					fmt.Sprintf("an error occured : %#v, data : %#v", e, strlog),
				)

			},
			Finally: func() {
			},
		}.Do()

		r.Logs.Info(fmt.Sprintf("✅ Sent to LinkIT: %s", fullURL))

		time.Sleep(time.Duration(5) * time.Second)
	}

	return nil
}

func (r *BaseModel) FetchAndUpdateARPUData(curdate_time time.Time) {
	// Step 1: Ambil kombinasi unik dari DB
	var summaries []struct {
		Country     string
		Operator    string
		Service     string
		SummaryDate time.Time
	}

	r.DB.Model(&entity.SummaryCampaign{}).
		// Distinct("country", "partner AS operator", "service").
		// Where("deleted_at IS NULL").
		Where("mo_received > 0").
		Where("summary_date = ? ", curdate_time).
		Scan(&summaries)

	for _, item := range summaries {
		/* currentYear := time.Now().Year()
		from := time.Date(currentYear-1, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		to := time.Now().AddDate(0, 0, -1).Format("2006-01-02") */

		// Bangun URL
		fullURL := fmt.Sprintf(
			"%s/api/v4/arpu/arpu90?from=%s&to=%s&country=%s&operator=%s&service=%s&to_renewal=%s",
			r.Config.APIARPU,
			curdate_time.Format("2006-01-02"),
			curdate_time.Format("2006-01-02"),
			url.QueryEscape(item.Country),
			url.QueryEscape(item.Operator),
			url.QueryEscape(item.Service),
			curdate_time.Format("2006-01-02"),
		)

		// 🔐 Ambil kredensial ARPU API
		encUsername := r.Config.ARPUUsername
		encPassword := r.Config.ARPUPassword

		username, _ := decryptEnv(encUsername)
		password, _ := decryptEnv(encPassword)

		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

		var (
			strlog string
			resp   []byte
			err    error
		)

		headerReq := map[string]string{
			"Accept":        "application/json",
			"Authorization": "Basic " + auth,
		}

		external.Block{
			Try: func() {

				if resp, _, _, strlog, err = external.Get(fullURL, headerReq, external.PHttp{
					Timeout:            5,
					KeepAlive:          3,
					IsDisableKeepAlive: true,
					MaxIdleConns:       5,
					IdleConnTimeout:    5,
					DisableCompression: false,
				}); err == nil {
					r.Logs.Info(strlog)
				} else {
					r.Logs.Error(strlog)
				}
			},
			Catch: func(e external.Exception) {

				r.Logs.Error(
					fmt.Sprintf("an error occured : %#v, data : %#v\n", e, strlog),
				)

			},
			Finally: func() {
			},
		}.Do()

		var arpuResp entity.ARPUResponse
		if err = json.Unmarshal(resp, &arpuResp); err != nil {
			r.Logs.Error(fmt.Sprintf("Error parse response, err : %#v, data : %#v\n", err, fullURL))
		}

		// 🔄 Loop hasil ARPU
		for _, d := range arpuResp.Data.Data {
			err := r.DB.Model(&entity.SummaryCampaign{}).
				Where("LOWER(adnet) = LOWER(?) AND LOWER(country) = LOWER(?) AND LOWER(partner) = LOWER(?) AND LOWER(service) = LOWER(?) AND summary_date = ?",
					d.Adnet, item.Country, item.Operator, item.Service, curdate_time).
				Updates(map[string]interface{}{
					"roi": d.Arpu90USDNet,
				}).Error
			if err != nil {
				r.Logs.Info(fmt.Sprintf("No Match arpu_update on adnet %s: %#v\n", d.Adnet, err))
			} else {
				r.Logs.Info(fmt.Sprintf("✅ ROI updated for adnet %s => %.2f\n", d.Adnet, d.Arpu90USDNet))
			}
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

func (r *BaseModel) SuccesRateLinkit(curdate_time time.Time) (entity.SuccessRateResponse, error) {
	var result entity.SuccessRateResponse

	var summaries []struct {
		Country     string
		Operator    string
		Service     string
		SummaryDate time.Time
	}

	// Ambil data summary campaign unik untuk hari ini
	if err := r.DB.Model(&entity.SummaryCampaign{}).
		Distinct("country", "partner AS operator", "service", "summary_date").
		Where("deleted_at IS NULL").
		Where("summary_date = ?", curdate_time).
		Where("mo_received > 0").
		Scan(&summaries).Error; err != nil {
		r.Logs.Error(fmt.Printf(" Failed to fetch summary data: %#v\n", err))
		return result, err
	}

	/* ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() */

	// Iterasi setiap kombinasi operator/service/date
	for _, item := range summaries {

		/* if i > 0 {
			<-ticker.C
		} */

		var (
			strlog      string
			resp        []byte
			status_code int
			err         error
		)

		fullURL := fmt.Sprintf(
			"%s/success-rate?operator=%s&service=%s&date=%s",
			r.Config.APILINKITDashboard,
			url.QueryEscape(strings.ToLower(item.Operator)),
			url.QueryEscape(strings.ToLower(item.Service)),
			url.QueryEscape(item.SummaryDate.Format("2006-01-02")),
		)

		external.Block{
			Try: func() {

				if resp, _, status_code, strlog, err = external.Post(fullURL, map[string]string{}, []byte(""), external.PHttp{
					Timeout:            5,
					KeepAlive:          3,
					IsDisableKeepAlive: true,
					MaxIdleConns:       5,
					IdleConnTimeout:    5,
					DisableCompression: false,
				}); err == nil {
					r.Logs.Info(strlog)
				} else {
					r.Logs.Error(strlog)
				}
			},
			Catch: func(e external.Exception) {

				r.Logs.Error(
					fmt.Sprintf("an error occured : %#v, data : %#v\n", e, strlog),
				)

			},
			Finally: func() {
			},
		}.Do()

		// Coba decode ke response normal
		var successRate entity.SuccessRateResponse
		if err := json.Unmarshal(resp, &successRate); err != nil {
			r.Logs.Error(fmt.Printf(" Failed to decode as SuccessRateResponse for %s/%s: %v", item.Operator, item.Service, err))
			continue
		}

		if status_code != 200 {
			// Coba decode pesan error
			var errorMsg struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(resp, &errorMsg); err != nil {
				r.Logs.Error(fmt.Sprintf("success_rate for %s/%s, code=%d but failed to parse message. Raw: %s", item.Operator, item.Service, successRate.Code, string(resp)))
			} else {
				r.Logs.Info(fmt.Printf("success_rate for %s/%s: %s", item.Operator, item.Service, errorMsg.Message))
			}
			continue
		}

		// Bersihkan "8.24%" jadi 8.24 float
		cleanRate := strings.TrimSuffix(successRate.Data.SuccessRate, "%")
		rateFloat, err := strconv.ParseFloat(cleanRate, 64)
		if err != nil {
			r.Logs.Error(fmt.Printf(" Failed to parse success rate '%s' for %s: %v", cleanRate, successRate.Data.Operator, err))
			continue
		}

		// Update successrate_fp di database
		err = r.DB.Model(&entity.SummaryCampaign{}).
			Where("LOWER(partner) = LOWER(?) AND LOWER(service) = LOWER(?) AND summary_date = ?",
				successRate.Data.Operator,
				successRate.Data.Service,
				curdate_time,
			).
			Updates(map[string]interface{}{
				"success_fp": rateFloat,
			}).Error

		if err != nil {
			r.Logs.Error(fmt.Printf("No Match success_rate for operator=%s service=%s: %v",
				successRate.Data.Operator, successRate.Data.Service, err))
		} else {
			r.Logs.Info(fmt.Printf("✅ successrate_fp updated: operator=%s service=%s => %.2f%%",
				successRate.Data.Operator, successRate.Data.Service, rateFloat))
		}

		time.Sleep(time.Duration(5) * time.Second)

		// if i < len(summaries)-1 {
		// 	time.Sleep(5 * time.Minute)
		// }
	}

	return result, nil
}
