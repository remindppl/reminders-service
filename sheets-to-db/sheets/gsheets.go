package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"sheets-to-db/common"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	scope = []string{
		"https://www.googleapis.com/auth/drive.metadata.readonly",
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	}

	env_token_path     = "TOKEN_PATH"
	env_svc_creds_path = "SVC_CREDS_PATH"
)

type Row []string
type Rows []Row

type Sheets interface {
	ReadSheetsData(context.Context) (Rows, error)
}

type GSheets struct {
	sheetId    string
	sheetRange string
}

func New(ctx context.Context, sheetId, sheetRange string) (Sheets, error) {
	if err := common.HasEnvironmentVars([]string{env_token_path, env_svc_creds_path}); err != nil {
		return nil, err
	}
	return &GSheets{
		sheetId:    sheetId,
		sheetRange: sheetRange,
	}, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	tokFile, _ := os.LookupEnv(env_token_path)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load token-file: %q. Err: %w", tokFile, err)
	}
	return config.Client(context.Background(), tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read provided credentials file: %q, %w", file, err)
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getSheetsClient(ctx context.Context) (*sheets.Service, error) {
	credsFile, _ := os.LookupEnv(env_svc_creds_path)
	b, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %q, Err: %v", credsFile, err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, err := getClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new http.Client. Err: %w", err)
	}
	return sheets.NewService(ctx, option.WithHTTPClient(client))
}

func (s *GSheets) ReadSheetsData(ctx context.Context) (Rows, error) {
	svc, err := getSheetsClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := svc.Spreadsheets.Values.Get(s.sheetId, s.sheetRange).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet: %q, range: %q. Err: %w", s.sheetId, s.sheetRange, err)
	}

	if len(resp.Values) == 0 {
		log.Printf("No Values found for given sheetID: %q, range: %q", s.sheetId, s.sheetRange)
		return nil, nil
	}

	var rows []Row
	for _, rv := range resp.Values {
		var row Row
		for _, cv := range rv {
			cellVal, ok := cv.(string)
			if !ok {
				cellVal = fmt.Sprint(cv)
			}
			row = append(row, cellVal)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
