package deltachat

import (
	"encoding/json"
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

// UnmarshalJSON is the method that satisfies the Unmarshaller interface.
func (self *Timestamp) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	self.Time = time.Unix(timestamp, 0)
	return nil
}

// MarshalJSON turns Timestamp back into an int.
func (self *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", self.Time.Unix())), nil
}
