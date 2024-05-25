package sheets

func (ss SpreadSheet) GetSheet(sheetName string) (*Sheet, error) {
	for _, sheet := range ss.Sheets {
		if sheet.Title == sheetName {
			return &sheet, nil
		}
	}
	return nil, ErrSheetNotFound
}

type SpreadSheet struct {
	Title      string  `json:"title"`
	SheetCount int     `json:"sheetCount"`
	Token      string  `json:"spreadsheetToken"`
	Sheets     []Sheet `json:"sheets"`
}
