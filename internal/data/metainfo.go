package data

import (
	"errors"
	"time"
)

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
	Length      *int64
	Md5sum      *string
	Files       []struct {
		Length int64
		Md5sum *string
		Path   []string
	}
}

type Metainfo struct {
	Info         info
	Announce     string
	AnnounceList [][]string
	CreationData *time.Time
	Comment      *string
	CreatedBy    *string
	Encoding     *string
}

func ParseMetainfo(value map[string]any) (*Metainfo, error) {
	metainfo := Metainfo{}

	announce, ok := value["announce"].(string)
	if !ok {
		return nil, errors.New("Missing announce")
	}
	metainfo.Announce = announce

	announceList, err := parseAnnounceList(value)
	if err != nil {
		return nil, err
	}
	metainfo.AnnounceList = announceList

	return &metainfo, nil
}

func parseAnnounceList(value map[string]any) ([][]string, error) {
	rawAnnounceList, ok := value["announce-list"].([]any)
	if !ok {
		return nil, nil
	}

	al := make([][]string, 0)
	for _, announceSubList := range rawAnnounceList {
		switch a := announceSubList.(type) {
		case []any:
			subList := make([]string, 0)
			for _, announceValues := range a {
				switch v := announceValues.(type) {
				case string:
					subList = append(subList, v)
				default:
					return nil, errors.New("error malformed announce")
				}
			}
			al = append(al, subList)
		default:
			return nil, errors.New("error malformed announce")
		}
	}

	return al, nil
}
