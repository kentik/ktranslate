package file

import (
	"bufio"
	"compress/gzip"
	"context"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
)

var (
	tags       string
	tagsCity   string
	tagsRegion string
)

func init() {
	flag.StringVar(&tags, "tag_map", "", "CSV file mapping tag ids to strings")
	flag.StringVar(&tagsCity, "geo_city_map", "", "CSV file mapping geo city ids to strings")
	flag.StringVar(&tagsRegion, "geo_region_map", "", "CSV file mapping geo region ids to strings")
}

type tagfunc struct {
	c string
	f func(int64) string
}

type FileTagMapper struct {
	logger.ContextL
	tags  map[uint32]map[string][2]string
	funcs map[string]tagfunc
	kvs   map[uint32]string
}

func NewFileTagMapper(log logger.Underlying, tagMapFilePath string) (*FileTagMapper, error) {
	ftm := FileTagMapper{
		ContextL: logger.NewContextLFromUnderlying(logger.SContext{S: "fileMapper"}, log),
		tags:     map[uint32]map[string][2]string{},
		funcs:    map[string]tagfunc{},
		kvs:      map[uint32]string{},
	}

	err := ftm.loadPath(tagMapFilePath)
	if err != nil {
		return nil, err
	}

	return &ftm, nil
}

func (ftm *FileTagMapper) loadPath(tagMapFilePath string) error {
	f, err := os.Open(tagMapFilePath)
	if err != nil {
		return err
	}
	var zr *gzip.Reader
	defer func() {
		if strings.HasSuffix(tagMapFilePath, ".gz") && zr != nil {
			zr.Close()
		}
		f.Close()
	}()
	var scanner *bufio.Scanner

	if strings.HasSuffix(tagMapFilePath, ".gz") {
		zr, err = gzip.NewReader(f)
		if err != nil {
			return err
		}
		scanner = bufio.NewScanner(zr)
	} else {
		scanner = bufio.NewScanner(f)
	}

	tm := map[uint32]map[string][2]string{}
	funcs := map[string]tagfunc{}
	kvs := map[uint32]string{}
	tmFound := 0
	for scanner.Scan() {
		pts := strings.SplitN(scanner.Text(), ",", 4)
		switch len(pts) {
		case 1:
			// Noop, just a blank line can skip.
		case 2:
			// Put into the basic kv map
			ida, err := strconv.Atoi(pts[0])
			if err != nil {
				continue
			}
			kvs[uint32(ida)] = pts[1]
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
				// Treat as a kv.
				ida, err := strconv.Atoi(pts[0])
				if err != nil {
					continue
				}
				kvs[uint32(ida)] = strings.Join(pts[1:], ",")
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
		return err
	}

	ftm.tags = tm
	ftm.funcs = funcs
	ftm.kvs = kvs
	ftm.Infof("Loaded %d tag mappings, %d functions and %d kvs from %s", tmFound, len(ftm.funcs), len(kvs), tagMapFilePath)

	return nil
}

func (ftm *FileTagMapper) LookupKV(k uint32) string {
	return ftm.kvs[k]
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

func (ftm *FileTagMapper) Run(ctx context.Context) {}
