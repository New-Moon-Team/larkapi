package sheets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	models "larkapi/models/spreadsheet"
	"net/http"
	"time"

	"golang.org/x/sync/semaphore"
)

func (shs *SheetsService) Login() error {
	b, err := json.Marshal(map[string]any{
		"app_id":     shs.config.AppId,
		"app_secret": shs.config.AppSecret,
	})
	if err != nil {
		return err
	}
	// fmt.Printf("body: %s\n", b)

	bf := bytes.NewReader(b)

	req, err := http.NewRequest(
		http.MethodPost,
		`https://open.larksuite.com/open-apis/auth/v3/tenant_access_token/internal`,
		bf,
	)
	if err != nil {
		return err
	}
	req = req.WithContext(shs.ctx)

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrLoginFailed
	}

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return ErrUnknownAuthResponse
	}

	var tat TenantAccessToken
	err = json.Unmarshal(b, &tat)
	if err != nil {
		return ErrUnknownAuthResponse
	}
	if tat.Code != 0 {
		return ErrLoginFailed
	}

	tat.OptainedTime = time.Now()

	shs.tententAccessToken = &tat
	// fmt.Printf("token object: %+v\n", tat)

	// log.Println("login ok")
	return nil
}
func (shs *SheetsService) GetSpreadSheet(spreadsheetId string) (*SpreadSheet, error) {
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
		fmt.Sprintf("https://open.larksuite.com/open-apis/sheets/v2/spreadsheets/%s/metainfo", spreadsheetId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	shs.tententAccessToken.Auth(req)
	req = req.WithContext(shs.ctx)

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
	// fmt.Printf("body: %s\n", b)

	var meta models.SpreadsheetMetaResponse
	err = json.Unmarshal(b, &meta)
	if err != nil {
		return nil, err
	}
	if meta.Code != 0 {
		return nil, ErrSheetNotFound
	}

	var ss SpreadSheet
	d := meta.Data
	p := d.Properties

	ss.SheetCount = p.SheetCount
	ss.Title = p.Title
	ss.Token = d.SpreadsheetToken

	for _, dsh := range d.Sheets {
		var sh Sheet
		b, err := json.Marshal(dsh)
		if err != nil {
			continue
		}

		err = json.Unmarshal(b, &sh)
		if err != nil {
			continue
		}

		sh.spreadsheetToken = ss.Token
		sh.service = shs
		ss.Sheets = append(ss.Sheets, sh)
	}

	nconn := int64(len(ss.Sheets))
	sem := semaphore.NewWeighted(nconn)
	for i := range ss.Sheets {
		sem.Acquire(shs.ctx, 1)
		go func(i int) {
			defer sem.Release(1)

			sh := &ss.Sheets[i]
			sh.loadHeaders()
		}(i)
	}
	sem.Acquire(shs.ctx, nconn)

	return &ss, nil
}
func (tat *TenantAccessToken) Auth(req *http.Request) error {
	// fmt.Println("auth")
	if tat.IsExpired() {
		return ErrNotAuthenticated
	}

	// fmt.Printf("token %s\n", tat.Token)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tat.Token))

	return nil
}
func (shs SheetsService) IsAuthenticated() bool {
	if shs.tententAccessToken == nil {
		return false
	}

	return !shs.tententAccessToken.IsExpired()
}
func (tat *TenantAccessToken) IsExpired() bool {
	return time.Now().After(tat.OptainedTime.Add(time.Duration(tat.Expire) * time.Second))
}

func NewService(ctx context.Context, config SheetsServiceConfig) (ISheetsService, error) {
	srv := &SheetsService{
		ctx:    ctx,
		config: config,
	}

	err := srv.Login()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

type ISheetsService interface {
	Login() error
	GetSpreadSheet(spreadsheetId string) (*SpreadSheet, error)
}
type SheetsService struct {
	ctx                context.Context
	config             SheetsServiceConfig
	tententAccessToken *TenantAccessToken
}

type SheetsServiceConfig struct {
	AutoRefreshToken bool
	AppId            string
	AppSecret        string
}
type TenantAccessToken struct {
	Code         int       `json:"code"`
	Msg          string    `json:"msg"`
	Token        string    `json:"tenant_access_token"`
	Expire       int       `json:"expire"`
	OptainedTime time.Time `json:"-"`
}
