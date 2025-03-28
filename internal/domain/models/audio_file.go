package models

import "time"

type AudioFile struct {
	id           int64
	link         string
	date_created time.Time
}
