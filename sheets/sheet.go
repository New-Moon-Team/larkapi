package sheets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	models "larkapi/models/cell"
	"net/http"
	"reflect"
)

func (sh Sheet) GetAllRaw() ([][]any, error) {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return nil, ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://open.larksuite.com/open-apis/sheets/v2/spreadsheets/%s/values/%s!A1:%s%d?valueRenderOption=ToString&dateTimeRenderOption=FormattedString",
			sh.spreadsheetToken,
			sh.ID,
			sh.Headers[len(sh.Headers)-1].Column,
			sh.RowCount,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = sh.service.tententAccessToken.Auth(req)
	if err != nil {
		return nil, err
	}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var readRes models.ReadSingleRangeResponse
	err = json.Unmarshal(b, &readRes)
	if err != nil {
		return nil, err
	}

	if readRes.Code != 0 {
		return nil, fmt.Errorf(readRes.Msg)
	}

	// var noEmpty [][]any
	// for _, r := range readRes.Data.ValueRange.Values {
	// 	empty := true

	// 	for _, c := range r {
	// 		if c != nil {
	// 			empty = false
	// 			break
	// 		}
	// 	}

	// 	if !empty {
	// 		noEmpty = append(noEmpty, r)
	// 	}
	// }
	// readRes.Data.ValueRange.Values = noEmpty

	return readRes.Data.ValueRange.Values, nil
}
func (sh Sheet) GetAll(result interface{}) error {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return err
		}
	}

	// Check if the result is a pointer to a slice
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result must be a pointer to a slice")
	}

	// Get the type of elements in the slice
	elementType := resultType.Elem().Elem()

	// Fetch all raw data from the sheet
	rawData, err := sh.GetAllRaw()
	if err != nil {
		return err
	}

	// Separate headers from the data
	headers := rawData[0]
	dataRows := rawData[1:]

	// Check if any field has the "_index" tag
	hasIndexTag := false
	for i := 0; i < elementType.NumField(); i++ {
		if elementType.Field(i).Tag.Get("sheet") == "_index" {
			hasIndexTag = true
			break
		}
	}

	if !hasIndexTag {
		return fmt.Errorf("none of the struct's fields has tag name '_index'")
	}

	// Iterate over each data row
	for rowIndex, dataRow := range dataRows {
		// Create a new element of the slice's element type
		newElement := reflect.New(elementType).Elem()

		// Check if all cells in the row are empty
		allEmpty := true
		for _, cellValue := range dataRow {
			if cellValue != nil {
				allEmpty = false
				break
			}
		}

		// If all cells are empty, skip this row
		if allEmpty {
			continue
		}

		// Iterate over each cell in the data row
		for cellIndex, cellValue := range dataRow {
			if cellValue == nil {
				continue
			}

			// Iterate over each field in the struct
			for fieldIndex := 0; fieldIndex < elementType.NumField(); fieldIndex++ {
				structField := elementType.Field(fieldIndex)

				// Get the value of the `sheet` tag
				sheetTag := structField.Tag.Get("sheet")

				// If the sheet tag matches the header, set the struct field value
				if sheetTag == headers[cellIndex] {
					// Convert the cell value to the type of the struct field
					convertedCellValue := reflect.ValueOf(cellValue).Convert(structField.Type)
					newElement.Field(fieldIndex).Set(convertedCellValue)
				}

				// If the sheet tag is "_index", set the row index as the value
				if sheetTag == "_index" {
					newElement.Field(fieldIndex).Set(reflect.ValueOf(rowIndex + 2))
				}
			}
		}

		// Append the new element to the result slice
		reflect.ValueOf(result).Elem().Set(reflect.Append(reflect.ValueOf(result).Elem(), newElement))
	}

	return nil
}

func (sh Sheet) SetMulti(data interface{}) error {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return err
		}
	}

	if reflect.TypeOf(data).Kind() != reflect.Slice && reflect.TypeOf(data).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a slice of struct")
	}

	var ranges []WriteValueRange
	for i := 0; i < reflect.ValueOf(data).Len(); i++ {
		d := reflect.ValueOf(data).Index(i).Interface()
		rgn, err := func(data any) ([]WriteValueRange, error) {
			dataType := reflect.TypeOf(data)
			if dataType.Kind() != reflect.Struct {
				return nil, fmt.Errorf("data must be a struct")
			}

			rowIndex := -1
			for i := 0; i < dataType.NumField(); i++ {
				field := dataType.Field(i)

				if field.Tag.Get("sheet") == "_index" {
					idx, ok := reflect.ValueOf(data).Field(i).Interface().(int)
					if ok {
						rowIndex = idx
					}
					break
				}
			}

			if rowIndex == -1 {
				return nil, fmt.Errorf("data must have a field with tag name '_index'")
			}

			ranges, err := sh.parseValueRange(data)
			if err != nil {
				return nil, err
			}

			return ranges, nil
		}(d)
		if err != nil {
			return err
		}
		ranges = append(ranges, rgn...)

	}

	return sh.writeRanges(ranges)
}

func (sh Sheet) Set(data interface{}) error {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return err
		}
	}

	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct")
	}

	rowIndex := -1
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		if field.Tag.Get("sheet") == "_index" {
			idx, ok := reflect.ValueOf(data).Field(i).Interface().(int)
			if ok {
				rowIndex = idx
			}
			break
		}
	}

	if rowIndex == -1 {
		return fmt.Errorf("data must have a field with tag name '_index'")
	}

	ranges, err := sh.parseValueRange(data)
	if err != nil {
		return err
	}

	return sh.writeRanges(ranges)
}

func (sh Sheet) writeRanges(ranges []WriteValueRange) error {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return err
		}
	}

	body := map[string]any{
		"valueRanges": ranges,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Printf("body: %s\n", b)
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(`https://open.larksuite.com/open-apis/sheets/v2/spreadsheets/%s/values_batch_update`, sh.spreadsheetToken),
		buf,
	)
	if err != nil {
		return err
	}
	req = req.WithContext(shs.ctx)
	err = sh.service.tententAccessToken.Auth(req)
	if err != nil {
		return err
	}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var writeRes models.WriteMultiRangeResponse
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Printf("res: %s\n", b)
	err = json.Unmarshal(b, &writeRes)
	if err != nil {
		return err
	}

	if writeRes.Code != 0 {
		return fmt.Errorf(writeRes.Msg)
	}

	return nil
}
func (sh *Sheet) loadHeaders() error {
	shs := sh.service
	if shs.IsAuthenticated() {
		if !shs.config.AutoRefreshToken {
			return ErrNotAuthenticated
		}

		err := shs.Login()
		if err != nil {
			return err
		}
	}

	var (
		firstCol = "A1"
		lastCol  = fmt.Sprintf("%s1", columnToLetter(sh.ColumnCount))
	)

	rgn := fmt.Sprintf("%s!%s:%s", sh.ID, firstCol, lastCol)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://open.larksuite.com/open-apis/sheets/v2/spreadsheets/%s/values/%s?valueRenderOption=ToString&dateTimeRenderOption=FormattedString",
			sh.spreadsheetToken,
			rgn,
		),
		nil,
	)
	if err != nil {
		return err
	}

	err = sh.service.tententAccessToken.Auth(req)
	if err != nil {
		return err
	}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var readRes models.ReadSingleRangeResponse
	err = json.Unmarshal(b, &readRes)
	if err != nil {
		return err
	}

	headersRow := readRes.Data.ValueRange.Values[0]
	var headers []SheetHeader

	for i, name := range headersRow {
		if name == nil {
			continue
		}

		h := SheetHeader{
			Name:   fmt.Sprintf("%v", name),
			Index:  i + 1,
			Column: columnToLetter(i + 1),
		}
		headers = append(headers, h)
	}
	sh.Headers = headers
	return nil
}
func (sh Sheet) parseValueRange(data interface{}) ([]WriteValueRange, error) {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data must be a struct")
	}

	rowIndex := -1
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		if field.Tag.Get("sheet") == "_index" {
			idx, ok := reflect.ValueOf(data).Field(i).Interface().(int)
			if ok {
				rowIndex = idx
			}
			break
		}
	}

	if rowIndex == -1 {
		return nil, fmt.Errorf("data must have a field with tag name '_index'")
	}

	var ranges []WriteValueRange
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tagName := field.Tag.Get("sheet")
		value := reflect.ValueOf(data).Field(i).Interface()

		for _, h := range sh.Headers {
			if h.Name == tagName {
				vrange := WriteValueRange{
					Range: fmt.Sprintf("%s!%s%d:%s%d", sh.ID, h.Column, rowIndex, h.Column, rowIndex),
					Values: [][]any{
						{value},
					},
				}
				ranges = append(ranges, vrange)
				break
			}
		}
	}
	return ranges, nil
}

type Sheet struct {
	ID               string         `json:"sheetId"`
	Title            string         `json:"title"`
	Index            int            `json:"index"`
	RowCount         int            `json:"rowCount"`
	ColumnCount      int            `json:"columnCount"`
	Headers          []SheetHeader  `json:"-"`
	spreadsheetToken string         `json:"-"`
	service          *SheetsService `json:"-"`
}

type SheetHeader struct {
	Name   string
	Index  int
	Column string
}

type WriteValueRange struct {
	Range  string  `json:"range"`
	Values [][]any `json:"values"`
}
