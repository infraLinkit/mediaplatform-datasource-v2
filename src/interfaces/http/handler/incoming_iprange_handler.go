package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm/clause"
)

func (h *IncomingHandler) DisplayIPRanges(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil {
		pageSize = 10
	}
	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.GlobalRequestFromDataTable{
		Page:     page,
		Action:   m["action"],
		Draw:     draw,
		PageSize: pageSize,
		Search:   m["search[value]"],
	}

	var (
		errResponse  error
		total_data   int64
		iprange_list []entity.IPRange
	)

	iprange_list, total_data, errResponse = h.DS.GetIPRanges(fe)

	r := entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}

	if errResponse == nil {

		r = entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            iprange_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) UploadIPRanges(c *fiber.Ctx) error {
	ipType := c.FormValue("ip_type")
	uploadDate := c.FormValue("upload_date")

	if ipType == "" || uploadDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ip_type and upload_date are required",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open file"})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	_, _ = reader.Read()

	redisData := make(map[string]map[string][]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 2 {
			continue
		}

		network := strings.TrimSpace(record[0])
		isp := strings.TrimSpace(record[1])
		if network == "" || isp == "" {
			continue
		}

		var countryCode, mobileCountryCode string
		err = h.DB.Raw(`
			SELECT c.code, c.mobile_country_code
			FROM operators o
			JOIN countries c ON o.country = c.code
			WHERE UPPER(o.name) = UPPER(?)
			LIMIT 1
		`, isp).Row().Scan(&countryCode, &mobileCountryCode)
		if err != nil || countryCode == "" || mobileCountryCode == "" {
			continue
		}

		h.DB.Exec(`
			INSERT INTO ip_ranges (network, isp, mobile_country_code, ip_type, upload_date)
			VALUES (?, ?, ?, ?, ?)
		`, network, isp, mobileCountryCode, ipType, uploadDate)

		key := countryCode + ":" + mobileCountryCode
		ispUpper := strings.ToUpper(isp)

		if _, ok := redisData[key]; !ok {
			redisData[key] = make(map[string][]string)
		}
		redisData[key][ispUpper] = append(redisData[key][ispUpper], network)
	}

	for key, value := range redisData {
		jsonString, err := json.Marshal(value)
		if err != nil {
			continue
		}
		h.DS.SetData(key, "$", string(jsonString))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "CSV processed and saved"})
}

func (h *IncomingHandler) UploadIPRangeRows(c *fiber.Ctx) error {
	ipType := c.FormValue("ip_type")
	uploadDate := c.FormValue("upload_date")

	if ipType == "" || uploadDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ip_type and upload_date are required",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open file"})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	_, err = reader.Read()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid csv header"})
	}

	const batchSize = 500
	var batch []entity.IPRangeCsvRow
	var insertedRows int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 6 {
			continue
		}

		network := strings.TrimSpace(record[0])
		isp := strings.TrimSpace(record[1])
		mobileCountryCode := strings.TrimSpace(record[5])
		if network == "" || isp == "" {
			continue
		}

		uploadDateParsed, err := time.Parse("2006-01-02", uploadDate)
		if err != nil {
			continue
		}

		batch = append(batch, entity.IPRangeCsvRow{
			IPType:            ipType,
			UploadDate:        uploadDateParsed,
			Network:           network,
			ISP:               isp,
			MobileCountryCode: mobileCountryCode,
		})
		insertedRows++

		if len(batch) >= batchSize {
			result := h.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "ip_type"}, {Name: "upload_date"}, {Name: "network"}, {Name: "isp"}},
				DoNothing: true,
			}).Create(&batch)
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		result := h.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "ip_type"}, {Name: "upload_date"}, {Name: "network"}, {Name: "isp"}},
			DoNothing: true,
		}).Create(&batch)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "CSV rows saved (only unique rows inserted)",
		"rows":    insertedRows,
	})
}

func (h *IncomingHandler) ImplementIPRange(c *fiber.Ctx) error {
	var body entity.ImplementIPRangeRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid json body",
		})
	}

	ipType := body.IPType
	month := body.UploadMonth

	if ipType == "" || month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ip_type and upload_month are required (format: yyyy-mm)",
		})
	}

	if len(month) != 7 || month[4] != '-' {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "upload_month format must be yyyy-mm",
		})
	}

	var csvRows []struct {
		Network    string
		ISP        string
		UploadDate time.Time
	}
	if err := h.DB.Model(&entity.IPRangeCsvRow{}).
		Select("network", "isp", "upload_date").
		Where("ip_type = ? AND to_char(upload_date, 'YYYY-MM') = ?", ipType, month).
		Scan(&csvRows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch ip range data"})
	}

	if len(csvRows) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no data found for given ip_type and month"})
	}

	var operators []struct {
		Name              string
		CountryCode       string
		MobileCountryCode string
	}
	if err := h.DB.Table("operators").
		Select("UPPER(operators.name) as name, countries.code as country_code, countries.mobile_country_code").
		Joins("JOIN countries ON countries.code = operators.country").
		Scan(&operators).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch operator data"})
	}

	// Build map: map[country_code][operator_name] = mobile_country_code
	operatorMap := make(map[string]map[string]string)
	for _, op := range operators {
		if _, ok := operatorMap[op.CountryCode]; !ok {
			operatorMap[op.CountryCode] = make(map[string]string)
		}
		operatorMap[op.CountryCode][op.Name] = op.MobileCountryCode
	}

	ipRanges := make([]entity.IPRange, 0)
	redisData := make(map[string]map[string]map[string][]string)

	for _, row := range csvRows {
		network := strings.TrimSpace(row.Network)
		isp := strings.TrimSpace(row.ISP)
		if network == "" || isp == "" {
			continue
		}

		upperISP := strings.ToUpper(isp)
		matched := false

		// Loop through all operator map to find matching name (substring)
		for countryCode, opMap := range operatorMap {
			for opName, mobileCode := range opMap {
				if strings.Contains(upperISP, opName) {
					ipRanges = append(ipRanges, entity.IPRange{
						Network:           network,
						ISP:               opName,
						MobileCountryCode: mobileCode,
						Country:           countryCode,
						IPType:            ipType,
						UploadDate:        row.UploadDate,
					})

					key := countryCode + ":" + mobileCode
					safeOpName := strings.ReplaceAll(opName, " ", "-")
					if _, ok := redisData[key]; !ok {
						redisData[key] = make(map[string]map[string][]string)
					}
					if _, ok := redisData[key][safeOpName]; !ok {
						redisData[key][safeOpName] = make(map[string][]string)
					}
					redisData[key][safeOpName][ipType] = append(redisData[key][safeOpName][ipType], network)

					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
	}

	if len(ipRanges) > 0 {
		err := h.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "network"}, {Name: "isp"}, {Name: "ip_type"}, {Name: "upload_date"}},
			DoNothing: true,
		}).Create(&ipRanges).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to insert ip_ranges"})
		}
	}

	for key, value := range redisData {
		if err := h.DS.SetDataIPSafe(key, value); err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to save to redis for key %s: %v", key, err))
			continue
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully implemented IP ranges for month: " + month,
	})
}

func (h *IncomingHandler) GetIPRangeFiles(c *fiber.Ctx) error {
	m := c.Queries()
	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil || pageSize <= 0 {
		pageSize = 10
	}
	draw, _ := strconv.Atoi(m["draw"])

	fe := entity.GlobalRequestFromDataTable{
		Page:     page,
		Draw:     draw,
		PageSize: pageSize,
		Search:   m["search[value]"],
	}

	rawResults, total, err := h.DS.GetIPRangeFiles(fe)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch ip range files",
		})
	}

	var results []entity.ResultIPRange
	for _, r := range rawResults {
		year := ""
		monthName := ""
		if len(r.Month) == 7 {
			year = r.Month[:4]
			monthNum := r.Month[5:]
			t, err := time.Parse("01", monthNum)
			if err == nil {
				monthName = t.Format("January")
			} else {
				monthName = monthNum
			}
		}
		filename := "GeoIP2-ISP-Blocks-" + r.IPType + "-" + year + "-" + monthName + ".csv"
		results = append(results, entity.ResultIPRange{
			IPType:   r.IPType,
			Month:    r.Month,
			Filename: filename,
		})
	}

	resp := entity.GlobalResponseWithDataTable{
		Draw:            draw,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            results,
		RecordsTotal:    int(total),
		RecordsFiltered: int(total),
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *IncomingHandler) DownloadIPRangeCSV(c *fiber.Ctx) error {

	var req entity.DownloadReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid json body"})
	}

	ipType := req.IPType
	monthParam := req.Month

	if ipType == "" || monthParam == "" || len(monthParam) != 7 || monthParam[4] != '-' {
		return c.Status(400).JSON(fiber.Map{"error": "ip_type and month (format: yyyy-mm) are required"})
	}

	year := monthParam[:4]
	monthNum := monthParam[5:]
	t, err := time.Parse("01", monthNum)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid month format, must be yyyy-mm"})
	}
	monthName := t.Format("January")
	filename := "GeoIP2-ISP-Blocks-" + ipType + "-" + year + "-" + monthName + ".csv"

	var rows []struct {
		Network           string
		ISP               string
		MobileCountryCode string
		IPType            string
		Country           string
	}
	if err := h.DB.Table("ip_ranges").
		Select("network, isp, mobile_country_code, ip_type, country").
		Where("ip_type = ? AND to_char(upload_date, 'YYYY-MM') = ?", ipType, monthParam).
		Order("network, isp").
		Scan(&rows).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch data"})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	writer := csv.NewWriter(c)
	defer writer.Flush()

	writer.Write([]string{"network", "isp", "mobile_country_code", "ip_type", "country"})
	for _, row := range rows {
		writer.Write([]string{
			row.Network,
			row.ISP,
			row.MobileCountryCode,
			row.IPType,
			row.Country,
		})
	}
	return nil
}
