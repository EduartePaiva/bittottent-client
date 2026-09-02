package data

import "time"

type FileMode string

const (
	SingleFile FileMode = "SINGLE_FILE"
	MultiFile  FileMode = "MULTI_FILE"
)

type info struct {
	PieceLength int64
	Pieces      string
	Private     *int64
	FileMode    FileMode
	Name        string
}

type Metainfo struct {
	Info         map[string]any
	Announce     string
	AnnounceList [][]string
	CreationData *time.Time
	Comment      *string
	CreatedBy    *string
	Encoding     *string
}
