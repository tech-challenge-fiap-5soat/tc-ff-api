package valueobject

import (
	"encoding/json"
	"testing"
	"time"

	vo "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

func TestCustomTime_UnmarshalJSON(t *testing.T) {
	t.Run("Valid JSON String", func(t *testing.T) {
		inputJSON := []byte(`"2024-01-20 21:16:09"`)
		expectedTime, _ := time.Parse(time.DateTime, "2024-01-20 21:16:09")

		var ct vo.CustomTime
		err := json.Unmarshal(inputJSON, &ct)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if ct.Time != expectedTime {
			t.Errorf("Expected time %v, got %v", expectedTime, ct.Time)
		}
	})

	t.Run("Null JSON String", func(t *testing.T) {
		inputJSON := []byte("null")

		var ct vo.CustomTime
		err := json.Unmarshal(inputJSON, &ct)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !ct.IsZero() {
			t.Errorf("Expected zero time, got %v", ct.Time)
		}
	})
}

func TestCustomTime_MarshalJSON(t *testing.T) {
	t.Run("Non-Zero Time", func(t *testing.T) {
		ct := vo.CustomTime{time.Date(2024, 01, 20, 15, 04, 05, 0, time.UTC)}
		expectedJSON := []byte(`"2024-01-20T15:04:05Z"`)

		resultJSON, err := json.Marshal(ct)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if string(resultJSON) != string(expectedJSON) {
			t.Errorf("Expected JSON %s, got %s", expectedJSON, resultJSON)
		}
	})

	t.Run("Zero Time", func(t *testing.T) {
		ct := vo.CustomTime{}
		expectedJSON := []byte("null")

		resultJSON, err := ct.MarshalJSON()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if string(resultJSON) != string(expectedJSON) {
			t.Errorf("Expected JSON %s, got %s", expectedJSON, resultJSON)
		}
	})
}

func TestCustomTime_IsSet(t *testing.T) {
	t.Run("Non-Zero Time", func(t *testing.T) {
		ct := vo.CustomTime{time.Now()}
		if !ct.IsSet() {
			t.Errorf("Expected IsSet to be true, got false")
		}
	})

	t.Run("Zero Time", func(t *testing.T) {
		ct := vo.CustomTime{}
		if ct.IsSet() {
			t.Errorf("Expected IsSet to be false, got true")
		}
	})
}
