package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) GetIPRanges(o entity.GlobalRequestFromDataTable) ([]entity.IPRange, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	query := r.DB.Model(&entity.IPRange{})

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("upload_date").Rows()
	defer rows.Close()

	var ss []entity.IPRange
	for rows.Next() {
		var s entity.IPRange
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetIPRangeFiles(o entity.GlobalRequestFromDataTable) ([]struct{ IPType, Month string }, int64, error) {
	var results []struct {
		IPType string
		Month  string
	}
	var total int64
	sub := r.DB.Table("ip_range_csv_rows").
		Select("ip_type, to_char(upload_date, 'YYYY-MM') as month").
		Group("ip_type, month")

	r.DB.Table("(?) as sub", sub).Count(&total)

	query := r.DB.Table("(?) as sub", sub).
		Order("month desc, ip_type")
	if o.PageSize > 0 {
		query = query.Limit(o.PageSize)
	}
	if o.Page > 0 {
		query = query.Offset((o.Page - 1) * o.PageSize)
	}
	err := query.Scan(&results).Error
	return results, total, err
}

func (h *BaseModel) GetDataIPSafe(key string) (map[string]map[string][]string, error) {
	ctx := context.Background()
	result := make(map[string]map[string][]string)

	raw, err := h.R.Conn().Do(ctx,
		h.R.Conn().B().JsonGet().Key(key).Path("$").Build()).ToString()

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "ERR") {
			return result, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *BaseModel) SetDataIPSafe(key string, data map[string]map[string][]string) error {
	ctx := context.Background()

	exists, err := h.R.Conn().Do(ctx, h.R.Conn().B().Exists().Key(key).Build()).ToInt64()
	if err != nil {
		return err
	}

	if exists == 0 {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return h.R.Conn().Do(ctx,
			h.R.Conn().B().JsonSet().Key(key).Path("$").Value(string(jsonBytes)).Build()).Error()
	}

	for isp, ipTypeMap := range data {
		safeISP := strings.ReplaceAll(isp, " ", "-")
		ispPath := fmt.Sprintf("$.%s", safeISP)

		val, err := h.R.Conn().Do(ctx,
			h.R.Conn().B().JsonGet().Key(key).Path(ispPath).Build()).ToString()

		if err != nil || val == "null" || val == "[]" {
			jsonVal, err := json.Marshal(ipTypeMap)
			if err != nil {
				return fmt.Errorf("failed to marshal ipTypeMap for %s: %w", isp, err)
			}

			if err := h.R.Conn().Do(ctx,
				h.R.Conn().B().JsonSet().Key(key).Path(fmt.Sprintf("$.%s", isp)).Value(string(jsonVal)).Build()).Error(); err != nil {
				return fmt.Errorf("failed to insert new ISP %s: %w", isp, err)
			}

			fmt.Printf("Added new ISP: %s\n", isp)
			continue
		}

		for ipType, newNetworks := range ipTypeMap {
			path := fmt.Sprintf("$.%s.%s", safeISP, ipType)

			var existing []string
			rawStr, err := h.R.Conn().Do(ctx,
				h.R.Conn().B().JsonGet().Key(key).Path(path).Build()).ToString()

			if err != nil {
				if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "ERR") {
					existing = []string{}
				} else {
					return err
				}
			} else {
				var nested [][]string
				if err := json.Unmarshal([]byte(rawStr), &nested); err == nil {
					for _, group := range nested {
						existing = append(existing, group...)
					}
				} else {
					if err := json.Unmarshal([]byte(rawStr), &existing); err != nil {
						var inner string
						if err := json.Unmarshal([]byte(rawStr), &inner); err != nil {
							return fmt.Errorf("unmarshal failed at %s: %w", path, err)
						}
						if err := json.Unmarshal([]byte(inner), &existing); err != nil {
							return fmt.Errorf("unmarshal inner failed at %s: %w", path, err)
						}
					}
				}
			}

			existingSet := make(map[string]bool)
			for _, n := range existing {
				n = strings.TrimSpace(n)
				if n != "" {
					existingSet[n] = true
				}
			}
			for _, n := range newNetworks {
				n = strings.TrimSpace(n)
				if n != "" && !existingSet[n] {
					existing = append(existing, n)
					existingSet[n] = true
				}
			}

			jsonVal, _ := json.Marshal(existing)
			if err := h.R.Conn().Do(ctx,
				h.R.Conn().B().JsonSet().Key(key).Path(path).Value(string(jsonVal)).Build()).Error(); err != nil {
				return err
			}
			fmt.Printf("Merged %s for ISP %s (%d items)\n", ipType, isp, len(existing))
		}
	}

	return nil
}