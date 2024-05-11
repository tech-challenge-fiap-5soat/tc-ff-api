package valueobject

import (
	"testing"

	orderStatus "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
	"github.com/stretchr/testify/assert"
)

func TestParseOrderStatus(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected orderStatus.OrderStatus
		err      bool
	}{
		{
			name:     "Valid order status",
			input:    "STARTED",
			expected: orderStatus.ORDER_STARTED,
			err:      false,
		},
		{
			name:     "Valid order status with lower case",
			input:    "started",
			expected: orderStatus.ORDER_STARTED,
			err:      false,
		},
		{
			name:     "Invalid order status",
			input:    "INVALID",
			expected: "",
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o, err := orderStatus.ParseOrderStatus(tc.input)

			if tc.err && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tc.err && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if o != tc.expected {
				t.Errorf("expected %s but got %s", tc.expected, o)
			}
		})
	}
}

func TestOrderStatusAsString(t *testing.T) {
	result := orderStatus.ORDER_STARTED.String()

	assert.Equal(t, "STARTED", result)
}
