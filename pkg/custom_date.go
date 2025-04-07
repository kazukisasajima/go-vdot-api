package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

const layout = "2006-01-02"

// JSON 変換（フロントとのやり取り）
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes
	t, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(layout))
}

// ✅ GORM対応（保存用）
func (d DateOnly) Value() (driver.Value, error) {
	return d.Format(layout), nil
}

// ✅ GORM対応（読込用）
func (d *DateOnly) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot convert %v to DateOnly", value)
	}
	d.Time = t
	return nil
}
