package main

import (
	"context"
	"log"
	"sheets-to-db/dao"
	"sheets-to-db/sheets"
)

var (
	// TODO: Get these values either through env vars or through CLI args.
	sheetID    = "1Fklhwg08szFk09SonD0UeCd6_8p__AAOMciyN4Nrmkw"
	sheetRange = "Form Responses 1!A2:K"
)

func main() {
	ctx := context.Background()

	db, err := dao.New(ctx, "followup_requests")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Fatalf("failed to close db connection. Err: %v", err)
		}
	}()

	s, err := sheets.New(ctx, sheetID, sheetRange)
	if err != nil {
		panic(err)
	}
	rows, err := s.ReadSheetsData(ctx)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		cReq, err := dao.MakeFollowupRequest(row)
		if err != nil {
			log.Fatalf("Failed to convert %v. Err: %v", row, err)
		}
		if err := db.Put(ctx, cReq); err != nil {
			log.Fatalf("Failed to add record to DB. Err: %v", err)
		}
	}
}
