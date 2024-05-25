package cell

type ReadSingleRangeResponse struct {
	Code int `json:"code"`
	Data struct {
		Revision         int    `json:"revision"`
		SpreadsheetToken string `json:"spreadsheetToken"`
		ValueRange       struct {
			MajorDimension string  `json:"majorDimension"`
			Range          string  `json:"range"`
			Revision       int     `json:"revision"`
			Values         [][]any `json:"values"`
		} `json:"valueRange"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type FindCellsResponse struct {
	Code int `json:"code"`
	Data struct {
		FindResult struct {
			MatchedCells        []string `json:"matched_cells"`
			MatchedFormulaCells []any    `json:"matched_formula_cells"`
			RowsCount           int      `json:"rows_count"`
		} `json:"find_result"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type WriteSingleRangeResponse struct {
	Code int `json:"code"`
	Data struct {
		Revision         int    `json:"revision"`
		SpreadsheetToken string `json:"spreadsheetToken"`
		UpdatedCells     int    `json:"updatedCells"`
		UpdatedColumns   int    `json:"updatedColumns"`
		UpdatedRange     string `json:"updatedRange"`
		UpdatedRows      int    `json:"updatedRows"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type WriteMultiRangeResponse struct {
	Code int `json:"code"`
	Data struct {
		Responses []struct {
			SpreadsheetToken string `json:"spreadsheetToken"`
			UpdatedCells     int    `json:"updatedCells"`
			UpdatedColumns   int    `json:"updatedColumns"`
			UpdatedRange     string `json:"updatedRange"`
			UpdatedRows      int    `json:"updatedRows"`
		} `json:"responses"`
		Revision         int    `json:"revision"`
		SpreadsheetToken string `json:"spreadsheetToken"`
	} `json:"data"`
	Msg string `json:"msg"`
}
