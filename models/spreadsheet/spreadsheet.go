package spreadsheet

type SpreadsheetMetaResponse struct {
	Code int `json:"code"`
	Data struct {
		Properties struct {
			OwnerUser  int64  `json:"ownerUser"`
			Revision   int    `json:"revision"`
			SheetCount int    `json:"sheetCount"`
			Title      string `json:"title"`
		} `json:"properties"`
		Sheets []struct {
			ColumnCount    int    `json:"columnCount"`
			FrozenColCount int    `json:"frozenColCount"`
			FrozenRowCount int    `json:"frozenRowCount"`
			Index          int    `json:"index"`
			RowCount       int    `json:"rowCount"`
			SheetID        string `json:"sheetId"`
			Title          string `json:"title"`
		} `json:"sheets"`
		SpreadsheetToken string `json:"spreadsheetToken"`
	} `json:"data"`
	Msg string `json:"msg"`
}
