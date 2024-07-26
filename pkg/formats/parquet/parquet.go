package parquet

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/rollup"

	"github.com/xitongsys/parquet-go-source/buffer"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	WriteParallelism = 4 // This is the default from examples.
)

type ParquetFormat struct {
	logger.ContextL
	compression  kt.Compression
	doGz         bool
	lastMetadata map[string]*kt.LastMetadata
	invalids     map[string]bool
	mux          sync.RWMutex
}

type ParquetMetric struct {
	Host       string            `parquet:"name=host, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	SourceType string            `parquet:"name=sourcetype, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Timestamp  int64             `parquet:"name=timestamp, type=INT64"`
	Strings    map[string]string `parquet:"name=string_map, type=MAP, convertedtype=MAP, keytype=BYTE_ARRAY, keyconvertedtype=UTF8, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
	Ints       map[string]int64  `parquet:"name=int_map, type=MAP, convertedtype=MAP, keytype=BYTE_ARRAY, keyconvertedtype=UTF8, valuetype=INT64"`
}

func NewFormat(log logger.Underlying, compression kt.Compression) (*ParquetFormat, error) {
	jf := &ParquetFormat{
		compression:  compression,
		ContextL:     logger.NewContextLFromUnderlying(logger.SContext{S: "parquetFormat"}, log),
		doGz:         false,
		invalids:     map[string]bool{},
		lastMetadata: map[string]*kt.LastMetadata{},
	}

	switch compression {
	case kt.CompressionNone:
	case kt.CompressionGzip:
	case kt.CompressionSnappy:
	default:
		return nil, fmt.Errorf("Invalid compression (%s): format parquet only supports none|gzip|snappy", compression)
	}

	return jf, nil
}

func (f *ParquetFormat) To(msgs []*kt.JCHF, serBuf []byte) (*kt.Output, error) {
	ms := make([]ParquetMetric, 0, len(msgs)*4)
	ct := time.Now().Unix()
	for _, m := range msgs {
		ms = append(ms, f.toParquetMetric(m, ct)...)
	}

	if len(ms) == 0 {
		return nil, nil
	}

	fw := buffer.NewBufferFileFromBytesZeroAlloc(serBuf)

	//write
	pw, err := writer.NewParquetWriter(fw, new(ParquetMetric), WriteParallelism)
	if err != nil {
		return nil, err
	}

	// Defaults from the example. No idea what these do.
	// https://github.com/xitongsys/parquet-go/blob/master/example/local_flat.go
	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.PageSize = 8 * 1024              //8K

	switch f.compression {
	case kt.CompressionNone:
		pw.CompressionType = parquet.CompressionCodec_UNCOMPRESSED
	case kt.CompressionGzip:
		pw.CompressionType = parquet.CompressionCodec_GZIP
	case kt.CompressionSnappy:
		pw.CompressionType = parquet.CompressionCodec_SNAPPY
	}
	for _, m := range ms {
		if err = pw.Write(m); err != nil {
			f.Errorf("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		return nil, err
	}
	fw.Close()

	return kt.NewOutput(fw.Bytes()), nil
}

func (f *ParquetFormat) From(raw *kt.Output) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)

	fr := buffer.NewBufferFileFromBytes(raw.Body)
	pr, err := reader.NewParquetReader(fr, new(ParquetMetric), WriteParallelism)
	if err != nil {
		return nil, err
	}
	num := int(pr.GetNumRows())
	read := 0

	process := func(ms []ParquetMetric) {
		for _, m := range ms {
			val := map[string]interface{}{
				"host":       m.Host,
				"sourcetype": m.SourceType,
				"timestamp":  m.Timestamp,
			}
			for k, v := range m.Strings {
				val[k] = v
			}
			for k, v := range m.Ints {
				val[k] = v
			}
			values = append(values, val)
		}
	}

	for i := 0; i < num/10; i++ {
		ms := make([]ParquetMetric, 10) //read 10 rows at a time.
		if err = pr.Read(&ms); err != nil {
			f.Errorf("Read error", err)
			continue
		}
		read += len(ms)
		process(ms)
	}

	//  Now get the rest.
	f.Infof("XXX %v %v", num, read)
	for i := 0; i < (num - read); i++ {
		ms := make([]ParquetMetric, 1)
		if err = pr.Read(&ms); err != nil {
			f.Errorf("Read error", err)
			continue
		}
		process(ms)
	}

	pr.ReadStop()
	fr.Close()

	return values, nil
}

func (f *ParquetFormat) Rollup(rolls []rollup.Rollup) (*kt.Output, error) {
	ct := time.Now().Unix()
	ms := f.toParquetMetricRollup(rolls, ct)

	if len(ms) == 0 {
		return nil, nil
	}

	fw := buffer.NewBufferFile()

	//write
	pw, err := writer.NewParquetWriter(fw, new(ParquetMetric), WriteParallelism)
	if err != nil {
		return nil, err
	}

	// Defaults from the example. No idea what these do.
	// https://github.com/xitongsys/parquet-go/blob/master/example/local_flat.go
	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.PageSize = 8 * 1024              //8K

	switch f.compression {
	case kt.CompressionNone:
		pw.CompressionType = parquet.CompressionCodec_UNCOMPRESSED
	case kt.CompressionGzip:
		pw.CompressionType = parquet.CompressionCodec_GZIP
	case kt.CompressionSnappy:
		pw.CompressionType = parquet.CompressionCodec_SNAPPY
	}
	for _, m := range ms {
		if err = pw.Write(m); err != nil {
			f.Errorf("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		return nil, err
	}
	fw.Close()

	return kt.NewOutput(fw.Bytes()), nil
}

func (f *ParquetFormat) toParquetMetricRollup(in []rollup.Rollup, ts int64) []ParquetMetric {
	ms := make([]ParquetMetric, 0, len(in))
	for _, roll := range in {
		dims := roll.GetDims()
		intm := map[string]int64{
			"count":    int64(roll.Count),
			"sum":      int64(roll.Metric),
			"min":      int64(roll.Min),
			"max":      int64(roll.Max),
			"interval": roll.Interval.Microseconds(),
		}
		strs := map[string]string{
			"name": "kentik.rollup." + roll.Name,
		}
		host := ""
		bad := false
		for i, pt := range strings.Split(roll.Dimension, roll.KeyJoin) {
			strs[dims[i]] = pt
			if pt == "0" || pt == "" {
				bad = true
			}
			if dims[i] == "device_name" {
				host = pt
			}
		}
		if !bad {
			ms = append(ms, ParquetMetric{
				Host:       host,
				SourceType: string(kt.ProviderRouter),
				Timestamp:  ts,
				Strings:    strs,
				Ints:       intm,
			})
		}
	}
	return ms
}

func (f *ParquetFormat) toParquetMetric(in *kt.JCHF, ts int64) []ParquetMetric {
	strs := map[string]string{}
	intm := map[string]int64{}
	host := in.DeviceName
	if host == "" {
		host = strconv.Itoa(int(in.DeviceId))
	}

	for k, v := range in.Flatten() {
		switch tv := v.(type) {
		case string:
			if tv != "" {
				strs[k] = tv
			}
		case int:
			if tv != 0 {
				intm[k] = int64(tv)
			}
		case int32:
			if tv != 0 {
				intm[k] = int64(tv)
			}
		case int64:
			if k == "timestamp" {
				strs["timestamp"] = time.Unix(tv, 0).Format(time.RFC3339)
			} else if tv != 0 {
				intm[k] = tv
			}
		}
	}
	return []ParquetMetric{
		ParquetMetric{
			Host:       host,
			SourceType: string(in.Provider),
			Timestamp:  ts,
			Strings:    strs,
			Ints:       intm,
		},
	}
}
