package avro

import (
	"fmt"
	"sync"
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSeriToAvro(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	f, err := NewFormat(nil, kt.CompressionSnappy)
	assert.NoError(err)
	res, err := f.To(kt.InputTesting, serBuf)
	assert.NoError(err)
	assert.NotNil(res)
	out, err := f.From(res)
	assert.NoError(err, "%s", string(res.Body))
	assert.NotNil(out)
	assert.Equal(len(kt.InputTesting), len(out))
	for i, _ := range out {
		assert.Equal(kt.InputTesting[i].Timestamp, out[i]["timestamp"])
		customStr := out[i]["custom_str"].(map[string]interface{})
		customInt := out[i]["custom_int"].(map[string]interface{})
		for k, v := range customStr {
			assert.Equal(kt.InputTesting[i].CustomStr[k], v.(string))
		}
		for k, v := range customInt {
			assert.Equal(kt.InputTesting[i].CustomInt[k], v.(int32))
		}
	}

	// Now try again with a different length of input
	inputNew := make([]*kt.JCHF, 4)
	for i := 0; i < 4; i++ {
		inputNew[i] = kt.NewJCHF()
		inputNew[i].CompanyId = kt.Cid(i + 11)
	}

	res, err = f.To(inputNew, serBuf)
	assert.NoError(err)
	assert.NotNil(res)
	out, err = f.From(res)
	assert.NoError(err, "%s", string(res.Body))
	assert.NotNil(out)
	assert.Equal(len(inputNew), len(out))

	// Now try again once more with a different length of input
	inputNew = make([]*kt.JCHF, 1)
	for i := 0; i < 1; i++ {
		inputNew[i] = kt.NewJCHF()
		inputNew[i].CompanyId = kt.Cid(i + 15)
	}
	res, err = f.To(inputNew, serBuf)
	assert.NoError(err)
	assert.NotNil(res)
	out, err = f.From(res)
	assert.NoError(err, "%s", string(res.Body))
	assert.NotNil(out)
	assert.Equal(len(inputNew), len(out))
}

func TestSeriToAvroManyItterations(t *testing.T) {
	serBuf := make([]byte, 0)
	assert := assert.New(t)
	f, err := NewFormat(nil, kt.CompressionSnappy)
	assert.NoError(err)

	// Now try again with a different length of input
	inputLen := 1000
	rounds := 10
	var wg sync.WaitGroup
	wg.Add(rounds)

	for round := 1; round <= rounds; round++ {
		go func(rr int) {
			defer wg.Done()
			inputNew := make([]*kt.JCHF, inputLen*rr)
			for i := 0; i < len(inputNew); i++ {
				inputNew[i] = kt.NewJCHF()
				inputNew[i].CompanyId = 12
				inputNew[i].Timestamp = int64(i)
				inputNew[i].CustomStr = map[string]string{"fooaaa": fmt.Sprintf("%d", i)}
				inputNew[i].CustomInt = map[string]int32{"aaaa": int32(i)}
				inputNew[i].CustomBigInt = map[string]int64{"eeeed": int64(i)}
			}
			res, err := f.To(inputNew, serBuf)
			assert.NoError(err)
			assert.NotNil(res)
			out, err := f.From(res)
			assert.NoError(err, "%s", string(res.Body))
			assert.NotNil(out)
			assert.Equal(len(inputNew), len(out))
		}(round)
	}

	wg.Wait()
	wg.Add(rounds)

	for round := rounds; round > 0; round-- {
		go func(rr int) {
			defer wg.Done()
			inputNew := make([]*kt.JCHF, inputLen*rr)
			for i := 0; i < len(inputNew); i++ {
				inputNew[i] = kt.NewJCHF()
				inputNew[i].CompanyId = 12
				inputNew[i].Timestamp = int64(i)
				inputNew[i].CustomStr = map[string]string{"fooaaa": fmt.Sprintf("%d", i)}
				inputNew[i].CustomInt = map[string]int32{"aaaa": int32(i)}
				inputNew[i].CustomBigInt = map[string]int64{"eeeed": int64(i)}
			}
			res, err := f.To(inputNew, serBuf)
			assert.NoError(err)
			assert.NotNil(res)
			out, err := f.From(res)
			assert.NoError(err, "%s", string(res.Body))
			assert.NotNil(out)
			assert.Equal(len(inputNew), len(out))
		}(round)
	}

	wg.Wait()
}

func BenchmarkAvro(b *testing.B) {
	serBuf := make([]byte, 0)
	f, err := NewFormat(nil, kt.CompressionSnappy)
	assert.NoError(b, err)

	// Now try again with a different length of input
	inputLen := 1000

	for i := 0; i < b.N; i++ {
		inputNew := make([]*kt.JCHF, inputLen)
		for i := 0; i < len(inputNew); i++ {
			inputNew[i] = kt.NewJCHF()
			inputNew[i].CompanyId = 12
			inputNew[i].Timestamp = int64(i)
			inputNew[i].CustomStr = map[string]string{"fooaaa": fmt.Sprintf("%d", i)}
			inputNew[i].CustomInt = map[string]int32{"aaaa": int32(i)}
			inputNew[i].CustomBigInt = map[string]int64{"eeeed": int64(i)}
		}
		res, err := f.To(inputNew, serBuf)
		assert.NoError(b, err)
		assert.NotNil(b, res)
		_, err = f.From(res)
		assert.NoError(b, err, "%s", string(res.Body))
	}
}
