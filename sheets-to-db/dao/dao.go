package dao

import (
	"time"
)

const (
	FIELD_ORDER_REQUEST_ID   = 0
	FIELD_ORDER_NAME         = 2
	FIELD_ORDER_LOCATION     = 3
	FIELD_ORDER_PHONE_NUM    = 4
	FIELD_ORDER_CROP         = 6
	FIELD_ORDER_CROP_VARIETY = 7
	FIELD_ORDER_SOW_DATE     = 8
)

var (
	inputTimeLayout  = "01/02/2006 15:04:05"
	outputTimeLayout = "2006-01-02T15:04:05"
)

type FollowUpRequest struct {
	RequestID   string    `json:"request_id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	PhoneNum    string    `json:"phone_num"`
	Crop        string    `json:"crop"`
	CropVariety string    `json:"crop_variety"`
	SowDate     time.Time `json:"sow_date"`
}

func MakeFollowupRequest(values []string) (*FollowUpRequest, error) {
	var err error
	fr := FollowUpRequest{}

	fr.RequestID, err = makeRequestID(values[FIELD_ORDER_REQUEST_ID])
	if err != nil {
		return nil, err
	}
	fr.Name = values[FIELD_ORDER_NAME]
	fr.Location = values[FIELD_ORDER_LOCATION]
	fr.PhoneNum = values[FIELD_ORDER_PHONE_NUM]
	fr.Crop = values[FIELD_ORDER_CROP]
	fr.CropVariety = values[FIELD_ORDER_CROP_VARIETY]
	fr.SowDate, err = getSowTimestamp(values[FIELD_ORDER_SOW_DATE])
	if err != nil {
		return nil, err
	}
	return &fr, nil
}

func makeRequestID(ts string) (string, error) {
	t, err := time.Parse(inputTimeLayout, ts)
	if err != nil {
		return "", err
	}
	return t.Format(outputTimeLayout), nil
}

func getSowTimestamp(ts string) (time.Time, error) {
	return time.Parse("01/02/2006", ts)
}
