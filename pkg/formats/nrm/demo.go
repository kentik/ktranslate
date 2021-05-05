package nrm

import (
	"math/rand"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
)

var (
	VALID_GEO = map[string]bool{
		"US": true,
		"VN": true,
		"AU": true,
	}
)

type Demozer struct {
	logger.ContextL
	waves     map[string]*SineWave
	period    uint32
	amplitude uint32
}

func NewDemozer(log logger.ContextL, period uint32, amplitude uint32) *Demozer {
	return &Demozer{
		ContextL:  log,
		waves:     map[string]*SineWave{},
		period:    period,
		amplitude: amplitude,
	}
}

// Move wave a random value from 0 forward.
func (d *Demozer) NewSineWaveOffset() *SineWave {
	w := NewSineWave(d, d.period, d.amplitude)
	w.PeriodVar = 40 // Change the period of each wave by +- 20 each time.
	offset := rand.Int31n(int32(d.period * 4))
	for i := 0; i < int(offset); i++ {
		<-w.Output
	}
	return w
}

// Update some of these values to show a change over time.
func (d *Demozer) demoize(ms []NRMetric) {
	adjusts := map[string]float64{}

	// For all of the metrics, get the idenifying value.
	for i, m := range ms {
		// if we have a wave for this value, adjust by the next tick.
		// Otherwise, start a wave for it.
		if vpc, ok := m.Attributes["vpc_identification"].(string); ok {
			if vpc != "" {
				key := m.Name + vpc
				adjust, ok := adjusts[key]
				if !ok {
					if _, ok := d.waves[key]; !ok { // 1 wave per vpc + metric + country
						d.waves[key] = d.NewSineWaveOffset()
					}
					adjust = <-d.waves[key].Output
					adjusts[key] = adjust
				}

				if vals, ok := m.Value.(map[string]uint64); ok {
					fv := float64(vals["sum"])
					adj := fv * adjust
					if fv+adj > 0 {
						vals["sum"] = uint64(fv + adj)
					} else {
						vals["sum"] = uint64((adj + fv) * -1)
					}
					ms[i].Value = vals
				}
			}
		}
	}
}

type SineWave struct {
	logger.ContextL
	Period    uint32
	Amplitude uint32
	PeriodVar int32

	// Generated data is written to this channel.
	Output chan float64
}

func NewSineWave(log logger.ContextL, period uint32, amplitude uint32) *SineWave {
	sw := &SineWave{
		ContextL:  log,
		Period:    period,
		Amplitude: amplitude,
		Output:    make(chan float64),
	}
	go sw.Generate()
	return sw
}

func (wave *SineWave) Generate() {
	var period uint32 = wave.Period
	var step float64 = float64(wave.Amplitude) / float64(period)

	currentValue := float64(0)
	sign := float64(1)
	for {
		for i := uint32(0); i < period/2; i++ {
			wave.Output <- currentValue * sign
			currentValue += step
		}
		for i := uint32(0); i < period/2; i++ {
			wave.Output <- currentValue * sign
			currentValue -= step
		}

		if sign < 0 { // If we are shifting our period, do so here afer 1 revolution
			if wave.PeriodVar > 0 {
				adj := rand.Int31n(wave.PeriodVar)
				op := period
				nv := int32(wave.Period) + (adj - (wave.PeriodVar / 2))
				if nv > 0 {
					period = uint32(nv)
					wave.Infof("New Period: %d, from %d", period, op)
				}
			}
		}

		sign = sign * -1
	}
}
