package file

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	tags string
)

func init() {
	flag.StringVar(&tags, "tag_map", "", "CSV file mapping tag ids to strings")
}

type FileTagMapper struct {
	logger.ContextL
	tags map[uint32][2]string
}

func NewFileTagMapper(log logger.Underlying, tagMapFilePath string) (*FileTagMapper, error) {
	f, err := os.Open(tagMapFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	tm := map[uint32][2]string{}
	for scanner.Scan() {
		pts := strings.SplitN(scanner.Text(), ",", 4)
		if len(pts) != 4 {
			continue
		}
		ida, err := strconv.Atoi(pts[2])
		if err != nil {
			continue
		}

		id := uint32(ida)
		tm[id] = [2]string{kt.FixupName(pts[1]), kt.FixupName(pts[3])}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	ftm := FileTagMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileMapper"}, log),
		tags:     tm,
	}
	ftm.Infof("Loaded %d tag mappings", len(tm))

	return &ftm, nil
}

func (ftm *FileTagMapper) LookupTagValue(cid kt.Cid, tagval uint32, colname string) (string, string, bool) {
	if tv, ok := ftm.tags[tagval]; ok {
		return tv[0], tv[1], ok
	}
	return "", "", false
}
