package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type AvgInt64 int64

func (ai AvgInt64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.3f", float64(ai)/1e9)), nil
}

type Stats struct {
	VMCount            int      `json:"vm_count"`
	RequestCount       int64    `json:"request_count"`
	AverageRequestTime AvgInt64 `json:"average_request_time"`
	mutex              sync.RWMutex
}

func (s *Stats) Wrap(path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		d := time.Since(start)
		s.updateStats(path, d)
	})
}

func (s *Stats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonfast.NewEncoder(w).Encode(s)
}

func (s *Stats) updateStats(path string, d time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldAvg := s.AverageRequestTime
	s.RequestCount += 1
	s.AverageRequestTime = AvgInt64(int64(s.AverageRequestTime)*(s.RequestCount-1)/s.RequestCount + d.Nanoseconds()/s.RequestCount)

	log.Infof("new request '%s' with duration: %d ns; average: %d ns(was %d ns)",
		path, int64(d), s.AverageRequestTime, oldAvg)
}
