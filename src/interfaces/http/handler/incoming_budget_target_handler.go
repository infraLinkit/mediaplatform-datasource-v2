package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func GetLastDayOfMonth(year, month int) time.Time {
	// 1. Create a date for the 1st of the next month.
	// Go automatically handles the rollover (e.g., month 12 -> 13 becomes next year).
	firstOfNextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	// 2. Subtract one day to get the last day of the previous month.
	return firstOfNextMonth.AddDate(0, 0, -1)
}

func (h *IncomingHandler) EditTargetBudget(c *fiber.Ctx) error {

	data := make(map[string]string)

	// Parse the body into the map
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	country := data["country"]
	operator := data["operator"]
	partner := data["partner"]
	service := data["service"]
	adnet := data["adnet"]
	year := data["year"]
	month := data["month"]
	budget, _ := strconv.ParseFloat(data["budget"], 64)

	_year, _ := strconv.Atoi(year)
	_month, _ := strconv.Atoi(month)

	start := GetLastDayOfMonth(_year, _month)
	end := start

	level := data["level"]

	if level == "country" {
		operator = ""
		partner = ""
		service = ""
		adnet = ""
	} else if level == "operator" {
		partner = ""
		service = ""
		adnet = ""
	} else if level == "partner" {
		service = ""
		adnet = ""
	} else if level == "service" {
		adnet = ""
	}

	service_list, _ := h.DS.GetTargetBudgetList(country, start, end, operator,
		partner, service, adnet)

	total_budget := 0.00
	//total_data := float64(len(service_list))

	avg_buget := 0.00

	TargetBudget := entity.TargetBudget{}
	TargetBudget.Country = country
	TargetBudget.Year, _ = strconv.Atoi(year)
	TargetBudget.Month, _ = strconv.Atoi(month)

	if level == "country" {

		total_data := 0
		for index, value := range service_list {
			total_budget += value.Budget
			fmt.Printf("Index: %d, Value: %+v\n", index, value)
			if value.Budget == 0 {
				total_data += 1
			}
		}

		if total_budget > budget {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "NOK",
				"error":  "Input budget < total budget !",
			})
		}

		avg_buget = (budget - total_budget) / float64(total_data)
		service_data := []entity.TargetBudgetDetail{}
		TargetBudget.Budget = budget
		for _, value := range service_list {
			if value.Budget == 0 {
				value.Budget = avg_buget
				value.BudgetPerDay = 0.0
				service_data = append(service_data, value)
			}
		}
		h.DS.AddTargetBudget(TargetBudget, service_data)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK",
			"error":  "",
		})
	} else {
		// GET CURRENT BUDGET
		current_budgets, is_exist := h.DS.GetTargetBudget(country, start, end, "", "", "", "")
		current_budget := entity.TargetBudget{}
		if !is_exist {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "NOK",
				"error":  "Main Budget not set up yeat !",
			})
		} else {
			current_budget = current_budgets[0]
		}

		service_data := []entity.TargetBudgetDetail{}
		total_data := 0.0
		total_budget_exclude := 0.0

		// GET ALL TOTAL BUDGET PER SERVICE INCLUDE NEW DATA
		var is_valid bool

		for _, value := range service_list {
			// GET TOTAL BUDGET EXCLUDE UPDATE BUDGET

			if level == "operator" {
				is_valid = operator != value.Operator
			} else if level == "partner" {
				is_valid = operator != value.Operator && partner != value.Partner
			} else if level == "service" {
				is_valid = operator != value.Operator && partner != value.Partner && service != value.Service
			} else if level == "adnet" {
				is_valid = operator != value.Operator && partner != value.Partner && service != value.Service && adnet != value.Adnet
			} else {
				is_valid = false
			}

			if is_valid {
				total_budget_exclude += float64(value.Budget)
			} else if value.Budget == 0 {
				total_data += 1
			}

		}

		if total_budget_exclude+budget > current_budget.Budget {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "NOK",
				"error":  "Input budget < total budget (operator) !!",
			})
		}

		// UPDATE BUDGET PER OPERATOR
		intYear, _ := strconv.Atoi(year)
		intMonth, _ := strconv.Atoi(month)
		if level == "adnet" {
			service_data = append(service_data, entity.TargetBudgetDetail{
				Country:  country,
				Year:     intYear,
				Month:    intMonth,
				Operator: operator,
				Partner:  partner,
				Service:  service,
				Adnet:    adnet,
			})
		} else {
			for _, value := range service_list {

				if level == "operator" {
					is_valid = operator != value.Operator
				} else if level == "partner" {
					is_valid = operator != value.Operator && partner != value.Partner
				} else if level == "service" {
					is_valid = operator != value.Operator && partner != value.Partner && service != value.Service
				} else if level == "adnet" {
					is_valid = operator != value.Operator && partner != value.Partner && service != value.Service && adnet != value.Adnet
				} else {
					is_valid = false
				}

				if is_valid {
					if value.Budget == 0 {
						avg_buget = budget / float64(total_data)
						value.Budget = avg_buget
						value.BudgetPerDay = 0.0
						service_data = append(service_data, value)
					}
				}
			}
		}

		h.DS.AddTargetBudget(current_budget, service_data)

	}

	fmt.Println("XXXX")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
		"error":  "",
	})
}
