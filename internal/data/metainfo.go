package data

import (
	"errors"
	"fmt"
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
	Private     int64
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
	CreationDate *time.Time
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

	metainfo.CreationDate = parseCreationDate(value)

	metainfo.Comment = extractMapValue[string](value, "comment")

	metainfo.CreatedBy = extractMapValue[string](value, "created by")

	metainfo.Encoding = extractMapValue[string](value, "encoding")

	info, err := parseInfo(value)
	if err != nil {
		return nil, err
	}
	metainfo.Info = *info

	return &metainfo, nil
}

func parseAnnounceList(value map[string]any) ([][]string, error) {
	rawValue, ok := value["announce-list"]
	if !ok {
		return nil, nil
	}

	rawAnnounceList, ok := rawValue.([]any)
	if !ok {
		return nil, errors.New("malformed announce: expected list")
	}

	announceList := make([][]string, 0, len(rawAnnounceList))
	for i, announceSubList := range rawAnnounceList {
		a, ok := announceSubList.([]any)
		if !ok {
			return nil, fmt.Errorf("malformed announce item %d is not a list", i)
		}

		subList := make([]string, 0, len(a))
		for j, announceValues := range a {
			v, ok := announceValues.(string)
			if !ok {
				return nil, fmt.Errorf("malformed announce item %d of list %d is not a string", j, i)
			}
			subList = append(subList, v)
		}
		announceList = append(announceList, subList)
	}
	return announceList, nil
}

func parseCreationDate(value map[string]any) *time.Time {
	date, ok := value["creation date"].(int64)
	if !ok {
		return nil
	}

	creationDate := time.Unix(date, 0)

	return &creationDate
}

func extractMapValue[V any](value map[string]any, key string) *V {
	val, ok := value[key].(V)
	if !ok {
		return nil
	}

	return &val
}

func parseInfo(value map[string]any) (*info, error) {
	infoResult := info{}

	infoRaw, ok := value["info"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("must contain an info dictionary")
	}

	pieceLength := extractMapValue[int64](infoRaw, "piece length")
	if pieceLength == nil {
		return nil, fmt.Errorf("must contain piece length")
	}
	infoResult.PieceLength = *pieceLength

	pieces := extractMapValue[string](infoRaw, "pieces")
	if pieces == nil {
		return nil, fmt.Errorf("must contain pieces")
	}
	infoResult.Pieces = *pieces

	private := extractMapValue[int64](infoRaw, "private")
	if private != nil {
		infoResult.Private = *private
	}

	name := extractMapValue[string](infoRaw, "name")
	if name != nil {
		infoResult.Name = *name
	}

	infoResult.Length = extractMapValue[int64](infoRaw, "length")

	infoResult.Md5sum = extractMapValue[string](infoRaw, "md5sum")

	if infoResult.Length != nil {
		infoResult.FileMode = SingleFile
	} else {
		infoResult.FileMode = MultiFile
	}

	// multi file parsing..... \/

	return &infoResult, nil
}
