package app_test

import (
	"bytes"
	"calendar/app"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName       string
		fileBytes      []byte
		startDate      time.Time
		endDate        time.Time
		expectedOutput []byte
	}{
		{},//Add test cases
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			var out bytes.Buffer
			app.Run(tc.fileBytes, tc.startDate, tc.endDate, &out)
			if strings.TrimSpace(string(tc.expectedOutput)) != strings.TrimSpace(out.String()) {
				t.Fatalf("%s failed.\n Expected: %s\n Received: %s", tc.testName, strings.TrimSpace(string(tc.expectedOutput)), strings.TrimSpace(out.String()))
			}
		})
	}
}
