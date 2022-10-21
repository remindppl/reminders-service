package main

import (
	"context"
	"log"
	"sheets-to-db/dao"
	"sheets-to-db/sheets"
)

func main() {
	ctx := context.Background()

	s, err := sheets.New(ctx, "", "")
	if err != nil {
		panic(err)
	}
	rows, err := s.ReadSheetsData(ctx)
	if err != nil {
		panic(err)
	}

	db := dao.New(ctx)
	for _, row := range rows {
		cReq, err := dao.MakeFollowupRequest(row)
		if err != nil {
			log.Fatalf("Failed to convert %v", row)
		}
		db.Put(ctx, cReq)
	}

}
