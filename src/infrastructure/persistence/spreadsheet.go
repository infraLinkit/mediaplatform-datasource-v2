package persistence

import (
	"fmt"
	"regexp"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"google.golang.org/api/sheets/v4"
)

type StatusData struct {
	Status       string
	StatusCode   string
	StatusDetail string
}

func (r *BaseModel) UpdateGoogleSheetPixel(GS *sheets.Service, ps entity.PixelStorage, conversion_name string, desc string) {
	sheetId, err := GetSpreadsheetID(ps.GoogleSheet)
	if err != nil {
		r.Logs.Info(fmt.Sprintf("Google sheet link not valid for campaign ID:  %#v\n", ps.CampaignId))
		r.Logs.Info(fmt.Sprintf("Google sheet link :  %#v ", ps.GoogleSheet))
		return
	}

	/* prop, err := GS.Spreadsheets.Get(sheetId).Fields("properties.title").Context(context.Background()).Do()
	if err != nil {
		Logs.Error(fmt.Sprintf("Failed to read title: %#v\n", err))
	} */

	resp, err := GS.Spreadsheets.Values.Get(sheetId, "Sheet1!A1:D1").Do()
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Failed to read sheet: %#v\n", err))
		return
	}

	if resp == nil {
		r.Logs.Error("Google sheet response nil")
		return
	}

	var header *sheets.ValueRange

	if len(resp.Values) == 0 && desc == "1" {
		header = &sheets.ValueRange{
			Range: "Sheet1!A1:D7",
			Values: [][]interface{}{
				{"#### INSTRUCTIONS ####"},
				{"# IMPORTANT: Remember to set the TimeZone value in the \"parameters\" row and/or in your Conversion Time column"},
				{"# For instructions on how to setup your data, visit http://goo.gl/T1C5Ov"},
				{},
				{"#### TEMPLATE ####"},
				{"Parameters: TimeZone=+0700"},
				{"Google Click ID", "Conversion Name", "Conversion Time", "MSISDN"},
			},
		}

		_, err := GS.Spreadsheets.Values.Update(sheetId, "Sheet1!A1:D7", header).ValueInputOption("RAW").Do()
		if err != nil {
			r.Logs.Error(fmt.Sprintf("Failed to insert header: %#v\n", err))
			return
		}

	} else {

		header = &sheets.ValueRange{
			Range: "Sheet1!A1:D1",
			Values: [][]interface{}{
				{"Google Click ID", "Conversion Name", "Conversion Time", "MSISDN"},
			},
		}

		_, err := GS.Spreadsheets.Values.Update(sheetId, "Sheet1!A1:D1", header).ValueInputOption("RAW").Do()
		if err != nil {
			r.Logs.Error(fmt.Sprintf("Failed to insert header: %#v\n", err))
			return
		}
	}

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			ps.Pixel,
			conversion_name,
			ps.PixelUsedDate,
			ps.Msisdn,
		}},
	}
	_, err = GS.Spreadsheets.Values.Append(sheetId, "Sheet1!A:D", values).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Google sheet input failed error:  %#v\n", err))
	}
}

func GetSpreadsheetID(url string) (string, error) {
	re := regexp.MustCompile(`https://docs\.google\.com/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("spreadsheet ID not found in URL")
	}

	return matches[1], nil
}
