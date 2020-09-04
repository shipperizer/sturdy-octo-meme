package status

import (
	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/stat"
)

const (
	Multiplier   = 10
	MaxDeviation = 10000
)

// Service provides exported method GetDecision. All other methods are unexported.
type Service struct {
	lags     []float64
	maxLag   float64
	lagCount int
}

// StoreLag stores the lag metric into memory
func (s *Service) StoreLag(lag int64) {
	log.Debug("Storing lag: ", lag)

	// add last lag
	s.lags = append(s.lags, float64(lag))

	// keep maxLags values
	s.lags = s.lags[1 : s.lagCount+1]
}

// ElaborateHealthStatus evaluates if there was any change in the last maxLags and if:
// * the lag keeps going higher without ever decreasing it will return false
// * if the lag is above maxLag will return false
// in any other occasion it will return true
func (s *Service) ElaborateHealthStatus() bool {
	if s.lags[s.lagCount-1] >= s.maxLag {
		log.Error("Lag is bigger than ", s.maxLag)
		return false
	}

	// check last 10% of lags
	meanRecent, devRecent := stat.MeanStdDev(
		s.lags[s.lagCount-int(s.lagCount/10):],
		nil,
	)

	// if mean of last 10% is 0 all is good
	if meanRecent == 0 {
		return true
	}

	log.Debug("Last 10% mean and stdev: ", meanRecent, devRecent)

	// if stalled recently
	if meanRecent > 0 && devRecent == 0 {
		return false
	}

	meanTotal, devTotal := stat.MeanStdDev(s.lags, nil)

	log.Debug("Max Allowed Deviation: ", MaxDeviation)
	log.Debug("Multiplier: ", Multiplier)
	log.Debug("Last 10% mean and stdev: ", meanRecent, devRecent)
	log.Debug("Total mean and stdev: ", meanTotal, devTotal)

	// if there was a spike in events
	if meanRecent > Multiplier*meanTotal {
		// check deviation is not bigger than MaxDeviation
		// or that the recent deviation is not too much bigger than the total deviation
		// TODO @shipperizer evaluation improve conditions
		if devRecent > MaxDeviation || devRecent > Multiplier*devTotal {
			return false
		}
	}

	return true
}

// NewService creates a new service
func NewService(lagCount, maxLag int) ServiceInterface {
	return &Service{
		lags:     make([]float64, lagCount),
		maxLag:   float64(maxLag),
		lagCount: lagCount,
	}
}
