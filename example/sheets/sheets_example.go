package main

import (
	"context"
	"fmt"
	"os"

	"github.com/New-Moon-Team/larkapi/sheets"
)

var (
	appId         string
	appSecret     string
	spreadsheetId = "RD2NsBMGrhXc02tKhUVlDoKIgle"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := sheets.NewService(ctx, sheets.SheetsServiceConfig{
		AutoRefreshToken: true,
		AppId:            appId,
		AppSecret:        appSecret,
	})
	if err != nil {
		panic(err)
	}

	ss, err := srv.GetSpreadSheet(spreadsheetId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", ss)

	sh, err := ss.GetSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", sh)

	rows, err := sh.GetAllRaw()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", rows)

	var data []sheetData
	err = sh.GetAll(&data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)

	data[0].Age = 50
	data[0].Foot = "foot"
	err = sh.SetMulti(data)
	if err != nil {
		panic(err)
	}
}

func init() {
	appId = os.Getenv("APP_ID")
	appSecret = os.Getenv("APP_SECRET")
}

type sheetData struct {
	Index int    `sheet:"_index"`
	Name  string `sheet:"Tên"`
	Age   int    `sheet:"Tuổi"`
	Foot  string `sheet:"Foot"`
}
