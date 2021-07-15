package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"

	mock_internal "github.com/arishlabroo/traintomo/internal/mocks"
)

func TestServer_FailsWhenGettingPostEndpoint(t *testing.T) {
	s := NewServer(NewInMemoryDB())
	req := httptest.NewRequest("GET", "/postschedule", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Error("failed request with", w.Result().StatusCode)
	}
}

func TestServer_FailsWhenPostingToGetEndpoint(t *testing.T) {
	s := NewServer(NewInMemoryDB())
	req := httptest.NewRequest("POST", "/getnextconflict?aftertime=03%3A04+PM", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Error("failed request with", w.Result().StatusCode)
	}
}

func TestServer_SavesFormattedValues(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_internal.NewMockDB(ctrl)
	s := NewServer(db)

	db.EXPECT().Set("ABCD", []string{"03:05 PM", "12:05 PM", "03:05 PM", "03:00 PM"}).Times(1)
	b, _ := json.Marshal(schedule{
		Name:     "ABCD",
		Schedule: []string{"3:05 PM", "12:05 PM", "03:05 PM", "03:00 PM"},
	})
	req := httptest.NewRequest("POST", "/postschedule", bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
}

func TestServer_HandlesPostSchedule(t *testing.T) {
	tests := []struct {
		testName string
		sch      schedule
		status   int
	}{
		{
			testName: "valid body",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusOK,
		},
		{
			testName: "numbers in name",
			sch: schedule{
				Name:     "AB22",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusOK,
		},
		{
			testName: "only numbers in name",
			sch: schedule{
				Name:     "4422",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusOK,
		},
		{
			testName: "short name",
			sch: schedule{
				Name:     "A",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusOK,
		},
		{
			testName: "empty name",
			sch: schedule{
				Name:     "",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "no name",
			sch: schedule{
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "only whitespace name",
			sch: schedule{
				Name:     "   ",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "whitespace in name",
			sch: schedule{
				Name:     "AB D",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "whitespace prefixed name",
			sch: schedule{
				Name:     " AB",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "special characters in name",
			sch: schedule{
				Name:     "(AB)",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "too long name",
			sch: schedule{
				Name:     "ABCDE",
				Schedule: []string{"1:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "invalid time",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{"16:20"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "invalid time out of range",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{"13:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "invalid time out of range 24hr",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{"25:20 PM"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "valid and invalid time",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{"1:20 PM", "apple"},
			},
			status: http.StatusBadRequest,
		},
		{
			testName: "empty schedule",
			sch: schedule{
				Name:     "ABCD",
				Schedule: []string{},
			},
			status: http.StatusOK,
		},
		{
			testName: "no schedule",
			sch: schedule{
				Name: "ABCD",
			},
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock_internal.NewMockDB(ctrl)
			s := NewServer(db)
			if tt.status == http.StatusOK {
				db.EXPECT().Set(tt.sch.Name, gomock.Any()).Times(1)
			}

			b, _ := json.Marshal(tt.sch)
			req := httptest.NewRequest("POST", "/postschedule", bytes.NewBuffer(b))
			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.status {
				t.Error("failed request")
			}
		})
	}
}

func TestServer_GetConflict(t *testing.T) {

	tests := []struct {
		testName    string
		schedules   []schedule
		requestTime string
		status      int
		conflict    string
	}{
		{
			testName: "no_request",
			status:   http.StatusBadRequest,
		},
		{
			testName:    "bad_request",
			requestTime: "25:45 LM",
			status:      http.StatusBadRequest,
		},
		{
			testName:    "nothing_in_store",
			requestTime: "5:45 AM",
			status:      http.StatusOK,
			conflict:    "",
		},
		{
			testName: "has_conflict",
			schedules: []schedule{
				{Name: "ABC", Schedule: []string{"12:05 PM"}},
				{Name: "XYZ", Schedule: []string{"12:05 PM"}},
			},
			requestTime: "07:04 PM",
			status:      http.StatusOK,
			conflict:    "12:05 PM",
		},
		{
			testName: "has_multiple_conflict",
			schedules: []schedule{
				{Name: "ABC", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
			},
			requestTime: "01:34 PM",
			status:      http.StatusOK,
			conflict:    "02:05 PM",
		},
		{
			testName: "no_conflict_after_request",
			schedules: []schedule{
				{Name: "ABC", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
			},
			requestTime: "03:34 PM",
			status:      http.StatusOK,
			conflict:    "12:05 PM",
		},
		{
			testName: "has_conflict_multiple_trains",
			schedules: []schedule{
				{Name: "ABC", Schedule: []string{"07:05 AM", "12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"07:05 AM"}},
			},
			requestTime: "11:34 AM",
			status:      http.StatusOK,
			conflict:    "12:05 PM",
		},
		{
			testName: "has_conflict_after_request_multiple_trains",
			schedules: []schedule{
				{Name: "ABC", Schedule: []string{"07:05 AM", "12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"12:05 PM", "01:05 PM", "02:05 PM"}},
				{Name: "XYZ", Schedule: []string{"07:05 AM"}},
			},
			requestTime: "03:34 PM",
			status:      http.StatusOK,
			conflict:    "07:05 AM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock_internal.NewMockDB(ctrl)
			s := NewServer(db)
			if tt.schedules != nil && len(tt.schedules) > 0 {
				trains := make([]string, 0, len(tt.schedules))
				for _, v := range tt.schedules {
					db.EXPECT().Fetch(v.Name).Return(v.Schedule).Times(1)
					trains = append(trains, v.Name)
				}
				db.EXPECT().Keys().Return(trains).Times(1)
			} else {
				db.EXPECT().Keys().Return([]string{}).AnyTimes()
			}

			req := httptest.NewRequest("GET", fmt.Sprintf("/getnextconflict?aftertime=%s", url.QueryEscape(tt.requestTime)), nil)
			w := httptest.NewRecorder()
			s.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.status {
				t.Error("failed test request with", w.Result().StatusCode)
			}

			if tt.status == http.StatusOK {

				bodyBytes, err := ioutil.ReadAll(w.Result().Body)
				if err != nil {
					t.Error("failed to parse body")
				}

				var res struct{ Conflict string }
				err = json.Unmarshal(bodyBytes, &res)
				if err != nil {
					t.Error("failed to parse body json")
				}
				if res.Conflict != tt.conflict {
					t.Error("Incorrect conflict found", res.Conflict)
				}

				w.Result().Body.Close()
			}
		})
	}

}
