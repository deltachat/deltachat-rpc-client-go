package deltachat

import (
	"encoding/json"
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

// UnmarshalJSON parses a Delta Chat timestamp into a Timestamp type.
func (self *Timestamp) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	self.Time = time.Unix(timestamp, 0)
	return nil
}

// MarshalJSON turns Timestamp back into the format expected by Delta Chat core.
func (self Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", self.Time.Unix())), nil
}
