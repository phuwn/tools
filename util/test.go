package util

import (
	"time"

	"bou.ke/monkey"
)

// MockRuntimeFunc - mock runtime function for phuwn/tools testing
func MockRuntimeFunc() {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2020, 9, 20, 17, 0, 58, 651387237, time.UTC)
	})
}
