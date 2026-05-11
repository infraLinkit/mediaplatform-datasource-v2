package persistence

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (m *BaseModel) CreateEmail(email *entity.Email) error {
	return m.DB.Create(email).Error
}

func (m *BaseModel) UpdateEmail(email *entity.Email) error {
	return m.DB.Updates(email).Error
}

func (m *BaseModel) DeleteEmail(id uint) error {
	return m.DB.Delete(&entity.Email{}, id).Error
}

func (r *BaseModel) GetEmail(o entity.GlobalRequestFromDataTable) ([]entity.Email, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Email{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("email_name ILIKE ?", "%"+search_value+"%").Or("email_purpose ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("email_name").Rows()
	defer rows.Close()

	var ss []entity.Email
	for rows.Next() {
		var s entity.Email
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) GetEmailByID(id uint) (entity.Email, error) {
	var email entity.Email
	if err := m.DB.First(&email, id).Error; err != nil {
		return entity.Email{}, err
	}
	return email, nil
}

func (m *BaseModel) CreateCountry(country *entity.Country) error {
	return m.DB.Create(country).Error
}

func (m *BaseModel) UpdateCountry(country *entity.Country) error {
	return m.DB.Updates(country).Error
}

func (m *BaseModel) DeleteCountry(id uint) error {
	return m.DB.Delete(&entity.Country{}, id).Error
}

func (r *BaseModel) GetCountry(o entity.GlobalRequestFromDataTable) ([]entity.Country, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Country{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("code ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Country
	for rows.Next() {
		var s entity.Country
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetCountryByCode(code string) (entity.Country, error) {

	var country entity.Country

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Country{})
	query = query.Where("code = ?", code).Take(&country)

	// Get the total count after applying filters

	return country, query.Error
}

func (r *BaseModel) GetContinent(o entity.GlobalRequestFromDataTable) ([]entity.Continent, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Country{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Continent
	for rows.Next() {
		var s entity.Continent
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateCompany(company *entity.Company) error {
	return m.DB.Create(company).Error
}

func (m *BaseModel) UpdateCompany(company *entity.Company) error {
	return m.DB.Updates(company).Error
}

func (m *BaseModel) DeleteCompany(id uint) error {
	return m.DB.Delete(&entity.Company{}, id).Error
}

func (r *BaseModel) GetCompany(o entity.GlobalRequestFromDataTableCompany) ([]entity.Company, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Company{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}
	orderBy := "id asc" // default

	if o.OrderColumn != "" && o.OrderDir != "" {
		orderBy = fmt.Sprintf("%s %s", o.OrderColumn, o.OrderDir)
	}
	query = query.Order(orderBy)

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Company
	for rows.Next() {
		var s entity.Company
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateDomain(domain *entity.Domain) error {
	return m.DB.Create(domain).Error
}

func (m *BaseModel) UpdateDomain(domain *entity.Domain) error {
	return m.DB.Updates(domain).Error
}

func (m *BaseModel) DeleteDomain(id uint) error {
	return m.DB.Delete(&entity.Domain{}, id).Error
}

func (r *BaseModel) GetDomain(o entity.GlobalRequestFromDataTable) ([]entity.Domain, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Domain{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("url ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("url").Rows()
	defer rows.Close()

	var ss []entity.Domain
	for rows.Next() {
		var s entity.Domain
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateOperator(operator *entity.Operator) error {
	return m.DB.Create(operator).Error
}

func (m *BaseModel) UpdateOperator(operator *entity.Operator) error {
	return m.DB.Updates(operator).Error
}

func (m *BaseModel) DeleteOperator(id uint) error {
	return m.DB.Delete(&entity.Operator{}, id).Error
}

func (r *BaseModel) GetOperator(o entity.GlobalRequestFromDataTable) ([]entity.Operator, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Operator{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Operator
	for rows.Next() {
		var s entity.Operator
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreatePartner(partner *entity.Partner) error {
	return m.DB.Create(partner).Error
}

func (m *BaseModel) UpdatePartner(partner *entity.Partner) error {
	return m.DB.Updates(partner).Error
}

func (m *BaseModel) DeletePartner(id uint) error {
	return m.DB.Delete(&entity.Partner{}, id).Error
}

func (r *BaseModel) GetPartner(o entity.GlobalRequestFromDataTable) ([]entity.Partner, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Partner{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").
			Or("operator ILIKE ?", "%"+search_value+"%").
			Or("aggregator ILIKE ?", "%"+search_value+"%").
			Or("country ILIKE ?", "%"+search_value+"%").
			Or("company ILIKE ?", "%"+search_value+"%").
			Or("client ILIKE ?", "%"+search_value+"%").
			Or("client_type ILIKE ?", "%"+search_value+"%").
			Or("url_postback ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Partner
	for rows.Next() {
		var s entity.Partner
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateService(service *entity.Service) error {
	return m.DB.Create(service).Error
}

func (m *BaseModel) UpdateService(service *entity.Service) error {
	return m.DB.Updates(service).Error
}

func (m *BaseModel) DeleteService(id uint) error {
	return m.DB.Delete(&entity.Service{}, id).Error
}

func (r *BaseModel) GetService(o entity.GlobalRequestFromDataTable) ([]entity.Service, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Service{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("service ILIKE ?", "%"+search_value+"%").
			Or("adn ILIKE ?", "%"+search_value+"%").
			Or("country ILIKE ?", "%"+search_value+"%").
			Or("operator ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("service").Rows()
	defer rows.Close()

	var ss []entity.Service
	for rows.Next() {
		var s entity.Service
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateAdnetList(adnet_list *entity.AdnetList) error {
	return m.DB.Create(adnet_list).Error
}

func (m *BaseModel) UpdateAdnetList(adnet_list *entity.AdnetList) error {
	return m.DB.Updates(adnet_list).Error
}

func (m *BaseModel) DeleteAdnetList(id uint) error {
	return m.DB.Delete(&entity.AdnetList{}, id).Error
}

func (r *BaseModel) GetAdnet(id string) (entity.AdnetList, error) {

	query := r.DB.Model(&entity.AdnetList{})
	query.Where("id = ?", id)

	rows, err := query.Rows()

	if err != nil {
		return entity.AdnetList{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var s entity.AdnetList
		r.DB.ScanRows(rows, &s)
		return s, nil
	}
	return entity.AdnetList{}, nil
}

func (r *BaseModel) GetAdnetList(o entity.GlobalRequestFromDataTable) ([]entity.AdnetList, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.AdnetList{})
	query.Select("*, api_url as api_url_before")
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").
			Or("code ILIKE ?", "%"+search_value+"%").
			Or("api_url ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()

	defer rows.Close()

	var ss []entity.AdnetList
	for rows.Next() {
		var s entity.AdnetList
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetAPIAdnetList(o entity.GlobalRequestFromDataTable) ([]entity.ApiPinReport, int64, error) {

	var (
		rows      *sql.Rows
		totalRows int64
	)

	// base query
	query := r.DB.Model(&entity.ApiPinReport{}).
		Select("adnet").
		Distinct("adnet")

	if o.Search != "" {
		searchValue := strings.TrimSpace(o.Search)
		query = query.Where(
			"adnet ILIKE ?",
			"%"+searchValue+"%",
		)
	}

	// count distinct adnet
	query.Count(&totalRows)

	// pagination
	if o.PageSize > 0 {
		query = query.Limit(o.PageSize)
		if o.Page > 0 {
			query = query.Offset((o.Page - 1) * o.PageSize)
		}
	}

	rows, err := query.
		Order("adnet").
		Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.ApiPinReport
	for rows.Next() {
		var s entity.ApiPinReport
		r.DB.ScanRows(rows, &s)
		result = append(result, s)
	}

	return result, totalRows, rows.Err()
}

func (r *BaseModel) GetAPIOperatorList(o entity.GlobalRequestFromDataTable) ([]entity.ApiPinReport, int64, error) {

	var (
		rows       *sql.Rows
		totalRows  int64
	)

	// base query
	query := r.DB.Model(&entity.ApiPinReport{}).
		Select("operator").
		Distinct("operator")

	if o.Search != "" {
		searchValue := strings.TrimSpace(o.Search)
		query = query.Where(
			"operator ILIKE ?",
			"%"+searchValue+"%",
		)
	}

	// count distinct operator
	query.Count(&totalRows)

	// pagination
	if o.PageSize > 0 {
		query = query.Limit(o.PageSize)
		if o.Page > 0 {
			query = query.Offset((o.Page - 1) * o.PageSize)
		}
	}

	rows, err := query.
		Order("operator").
		Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.ApiPinReport
	for rows.Next() {
		var s entity.ApiPinReport
		r.DB.ScanRows(rows, &s)
		result = append(result, s)
	}

	return result, totalRows, rows.Err()
}

func (r *BaseModel) GetAPIServiceList(o entity.GlobalRequestFromDataTable) ([]entity.ApiPinReport, int64, error) {

	var (
		rows       *sql.Rows
		totalRows  int64
	)

	// base query
	query := r.DB.Model(&entity.ApiPinReport{}).
		Select("service").
		Distinct("service")

	if o.Search != "" {
		searchValue := strings.TrimSpace(o.Search)
		query = query.Where(
			"service ILIKE ?",
			"%"+searchValue+"%",
		)
	}

	// count distinct service
	query.Count(&totalRows)

	// pagination
	if o.PageSize > 0 {
		query = query.Limit(o.PageSize)
		if o.Page > 0 {
			query = query.Offset((o.Page - 1) * o.PageSize)
		}
	}

	rows, err := query.
		Order("service").
		Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.ApiPinReport
	for rows.Next() {
		var s entity.ApiPinReport
		r.DB.ScanRows(rows, &s)
		result = append(result, s)
	}

	return result, totalRows, rows.Err()
}

func (r *BaseModel) GetAgency(o entity.GlobalRequestFromDataTable) ([]entity.Agency, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Agency{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("code ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Agency
	for rows.Next() {
		var s entity.Agency
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateAgency(agency *entity.Agency) error {
	return m.DB.Create(agency).Error
}

func (m *BaseModel) UpdateAgency(agency *entity.Agency) error {
	tx := m.DB.Begin()

	if err := tx.Model(&entity.Agency{}).
		Where("id = ?", agency.ID).
		Updates(agency).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("mainstream_groups").
		Where("agency = ?", agency.Name).
		Update("api_url", agency.APIURL).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (m *BaseModel) DeleteAgency(id uint) error {
	return m.DB.Delete(&entity.Agency{}, id).Error
}

func (r *BaseModel) GetChannel(o entity.GlobalRequestFromDataTable) ([]entity.Channel, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Channel{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Channel
	for rows.Next() {
		var s entity.Channel
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateChannel(channel *entity.Channel) error {
	return m.DB.Create(channel).Error
}

func (m *BaseModel) UpdateChannel(channel *entity.Channel) error {
	return m.DB.Updates(channel).Error
}

func (m *BaseModel) DeleteChannel(id uint) error {
	return m.DB.Delete(&entity.Channel{}, id).Error
}

func (r *BaseModel) GetMainstreamGroup(o entity.GlobalRequestFromDataTable) ([]entity.MainstreamGroup, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.MainstreamGroup{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("agency ILIKE ?", "%"+search_value+"%").Or("channel ILIKE ?", "%"+search_value+"%").Or("service ILIKE ?", "%"+search_value+"%").Or("unique_domain ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.MainstreamGroup
	for rows.Next() {
		var s entity.MainstreamGroup
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateMainstreamGroup(mainstreamGroup *entity.MainstreamGroup) error {
	return m.DB.Create(mainstreamGroup).Error
}

func (m *BaseModel) UpdateMainstreamGroup(mainstreamGroup *entity.MainstreamGroup) error {
	return m.DB.Updates(mainstreamGroup).Error
}

func (m *BaseModel) DeleteMainstreamGroup(id uint) error {
	return m.DB.Delete(&entity.MainstreamGroup{}, id).Error
}

func (m *BaseModel) FindMainstreamGroupByID(id uint) (*entity.MainstreamGroup, error) {
	var mainstreamGroup entity.MainstreamGroup
	if err := m.DB.First(&mainstreamGroup, id).Error; err != nil {
		return nil, err
	}
	return &mainstreamGroup, nil
}

func (r *BaseModel) GetDomainService(o entity.GlobalRequestFromDataTable) ([]entity.DomainService, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.DomainService{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.DomainService
	for rows.Next() {
		var s entity.DomainService
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateCompanyGroup(companyGroup *entity.CompanyGroup) error {
	return m.DB.Create(companyGroup).Error
}

func (m *BaseModel) UpdateCompanyGroup(companyGroup *entity.CompanyGroup) error {
	return m.DB.Updates(companyGroup).Error
}

func (m *BaseModel) DeleteCompanyGroup(id uint) error {
	return m.DB.Delete(&entity.CompanyGroup{}, id).Error
}

func (r *BaseModel) GetCompanyGroup(o entity.GlobalRequestFromDataTableCompany) ([]entity.CompanyGroup, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.CompanyGroup{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}
	orderBy := "id asc" // default

	if o.OrderColumn != "" && o.OrderDir != "" {
		orderBy = fmt.Sprintf("%s %s", o.OrderColumn, o.OrderDir)
	}
	query = query.Order(orderBy)

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.CompanyGroup
	for rows.Next() {
		var s entity.CompanyGroup
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}
