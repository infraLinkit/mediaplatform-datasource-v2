package cmd

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/spf13/cobra"
)

var migrateEntities = []interface{}{
	&entity.TargetBudget{}, &entity.TargetBudgetDetail{}, &entity.SummaryDashboard{},
	&entity.Campaign{}, &entity.CampaignDetail{}, &entity.MO{}, &entity.PixelStorage{},
	&entity.ClickStorage{}, &entity.Postback{}, &entity.SummaryCampaign{},
	&entity.SummaryCampaignBilling{}, &entity.DataClicked{}, &entity.DataLanding{},
	&entity.DataRedirect{}, &entity.DataTraffic{}, &entity.ApiPinReport{},
	&entity.ApiPinPerformance{}, &entity.Menu{}, &entity.Country{}, &entity.Continent{},
	&entity.Company{}, &entity.CompanyGroup{}, &entity.Domain{}, &entity.Operator{},
	&entity.Partner{}, &entity.Service{}, &entity.AdnetList{}, &entity.SummaryMo{},
	&entity.SummaryCr{}, &entity.SummaryCapping{}, &entity.SummaryRatio{}, &entity.Role{},
	&entity.Permission{}, &entity.User{}, &entity.CcEmail{}, &entity.Email{},
	&entity.DetailUser{}, &entity.UserAdnet{}, &entity.LpDesignType{}, &entity.Agency{},
	&entity.Channel{}, &entity.MainstreamGroup{}, &entity.SummaryLanding{},
	&entity.IPRangeCsvRow{}, &entity.IPRange{}, &entity.IncSummaryCampaign{},
	&entity.IncSummaryCampaignHour{}, &entity.SummaryTraffic{}, &entity.BudgetIO{},
	&entity.SummaryBudgetIO{}, &entity.UserCompany{}, &entity.DomainService{},
	&entity.HistoryCappingKey{}, &entity.OperatorAlias{},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run DB schema migration (AutoMigrate) and exit",
	Long: `Run GORM AutoMigrate untuk semua entity yang terdaftar di migrateEntities.
Dipakai sebagai standalone process (init container / CI step / manual run) supaya
schema migration terpisah dari app server lifecycle.

Env:
  AUTO_MIGRATE_TIMEOUT_MIN  context timeout (default 5 menit)`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.InitCfg()
		c, err := cfg.Initiate("migrate")
		if err != nil {
			log.Fatalf("init setup failed: %v", err)
		}

		timeoutMin, _ := strconv.Atoi(os.Getenv("AUTO_MIGRATE_TIMEOUT_MIN"))
		if timeoutMin <= 0 {
			timeoutMin = 5
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMin)*time.Minute)
		defer cancel()

		start := time.Now()
		c.Logs.Infof("[migrate] start AutoMigrate, timeout=%dm, entities=%d", timeoutMin, len(migrateEntities))
		if err := c.DB.WithContext(ctx).AutoMigrate(migrateEntities...); err != nil {
			log.Fatalf("[migrate] AutoMigrate failed after %s: %v", time.Since(start), err)
		}
		c.Logs.Infof("[migrate] AutoMigrate done in %s", time.Since(start))
	},
}
