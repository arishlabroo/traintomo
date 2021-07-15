package internal

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

type queryController struct {
	db DB
}

func newQueryController(db DB) queryController {
	return queryController{db: db}
}

func (s *queryController) handleGetNextConflict() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		keys, ok := r.URL.Query()["aftertime"]

		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requested := keys[0]
		t, err := time.Parse("3:04 PM", requested)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(struct{ Conflict string }{
			Conflict: s.getConflictAfter(t),
		})
	}
}

func (s *queryController) getConflictAfter(t time.Time) string {
	conflicts := s.getConflictingTimes()
	if len(conflicts) < 1 {
		return ""
	}

	for _, v := range conflicts {
		if v.After(t) {
			return v.Format(timeFormat)
		}
	}
	return conflicts[0].Format(timeFormat)

}

func (s *queryController) getConflictingTimes() []time.Time {
	trafficByTime := make(map[string]int)
	trains := s.db.Keys()
	for _, tr := range trains {
		if times, ok := s.db.Fetch(tr).([]string); ok {
			for _, t := range times {
				trafficByTime[t] = trafficByTime[t] + 1
			}
		}
	}

	conflicts := make([]time.Time, 0)
	for k, v := range trafficByTime {
		if v > 1 {
			if t, err := time.Parse(timeFormat, k); err == nil {
				conflicts = append(conflicts, t)
			}
		}
	}

	sort.SliceStable(conflicts, func(i, j int) bool {
		return conflicts[i].Before(conflicts[j])
	})

	return conflicts
}
