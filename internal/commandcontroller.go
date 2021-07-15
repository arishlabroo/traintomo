package internal

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"
)

type schedule struct {
	Name     string
	Schedule []string
}

var nameRgx = regexp.MustCompile("^[a-zA-Z0-9]{1,4}$")

func (s schedule) validateAndFormat() (bool, []string) {
	if !nameRgx.MatchString(s.Name) {
		return false, nil
	}

	if s.Schedule == nil { //empty slice is allowed
		return false, nil
	}

	formatted := make([]string, 0, len(s.Schedule))

	for _, v := range s.Schedule {
		t, err := time.Parse("3:04 PM", v)
		if err != nil {
			return false, nil
		}
		formatted = append(formatted, t.Format(timeFormat))
	}
	return true, formatted

}

type commandController struct {
	db DB
}

func newCommandController(db DB) commandController {
	return commandController{db: db}
}

func (s *commandController) handlePostSchedule() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var sch schedule
		err := decoder.Decode(&sch)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		valid, formatted := sch.validateAndFormat()
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.db.Set(sch.Name, formatted)

		w.WriteHeader(http.StatusOK)
	}
}
