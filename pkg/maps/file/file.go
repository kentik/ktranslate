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

type tagfunc struct {
	c string
	f func(int64) string
}

type FileTagMapper struct {
	logger.ContextL
	tags  map[uint32]map[string][2]string
	funcs map[string]tagfunc
}

func NewFileTagMapper(log logger.Underlying, tagMapFilePath string) (*FileTagMapper, error) {
	f, err := os.Open(tagMapFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	ftm := FileTagMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileMapper"}, log),
	}

	tm := map[uint32]map[string][2]string{}
	funcs := map[string]tagfunc{}
	tmFound := 0
	for scanner.Scan() {
		pts := strings.SplitN(scanner.Text(), ",", 4)
		switch len(pts) {
		case 1:
			// Noop, just a blank line can skip.
		case 3: // its a function.
			switch pts[2] {
			case "to_hex":
				funcs[pts[0]] = tagfunc{
					c: kt.FixupName(pts[1]),
					f: func(in int64) string {
						return strconv.FormatInt(in, 16)
					},
				}
			default:
				ftm.Errorf("Invalid function %v, skipping", pts)
			}
		case 4: // its a fixed mapping.
			ida, err := strconv.Atoi(pts[2])
			if err != nil {
				continue
			}

			id := uint32(ida)
			if _, ok := tm[id]; !ok {
				tm[id] = map[string][2]string{kt.FixupName(pts[0]): [2]string{kt.FixupName(pts[1]), kt.FixupName(pts[3])}}
			} else {
				tm[id][kt.FixupName(pts[0])] = [2]string{kt.FixupName(pts[1]), kt.FixupName(pts[3])}
			}
			tmFound++
		default: // its a mistake.
			ftm.Errorf("Invalid line %v, skipping", pts)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	ftm.tags = tm
	ftm.funcs = funcs
	ftm.Infof("Loaded %d tag mappings and %d functions", tmFound, len(ftm.funcs))

	return &ftm, nil
}

func (ftm *FileTagMapper) LookupTagValue(cid kt.Cid, tagval uint32, colname string) (string, string, bool) {
	if tf, ok := ftm.funcs[colname]; ok {
		return tf.c, tf.f(int64(tagval)), ok
	}
	if tvv, ok := ftm.tags[tagval]; ok {
		if tv, ok := tvv[colname]; ok {
			return tv[0], tv[1], ok
		}
	}
	return "", "", false
}

func (ftm *FileTagMapper) LookupTagValueBig(cid kt.Cid, tagval int64, colname string) (string, string, bool) {
	if tf, ok := ftm.funcs[colname]; ok {
		return tf.c, tf.f(tagval), ok
	}
	if tvv, ok := ftm.tags[uint32(tagval)]; ok {
		if tv, ok := tvv[colname]; ok {
			return tv[0], tv[1], ok
		}
	}
	return "", "", false
}
