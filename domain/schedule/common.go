package schedule

import (
	"strings"
	"time"
)

const hourMinuteFormat = "15:04"
const dateTimeFormat = "2006-01-02T15:04:05"
const dateTimeFormatWithTimeStamp = "2006-01-02T15:04:05-07:00"

type JaneTime struct {
	time.Time
}

func NewJaneTime(input time.Time) JaneTime {
	return JaneTime{
		Time: time.Date(
			input.Year(),
			input.Month(),
			input.Day(),
			input.Hour(),
			input.Minute(),
			0,
			0,
			time.Local,
		),
	}
}

func (janeTime *JaneTime) UnmarshalJSON(bytes []byte) error {
	timeString := strings.Trim(string(bytes), "\"")
	if timeString == "null" {
		janeTime.Time = time.Time{}
		return nil
	}
	parsedTime, err := time.Parse(dateTimeFormat, timeString)
	janeTime.Time = parsedTime
	return err
}

func (janeTime JaneTime) MarshalJSON() ([]byte, error) {
	timeString := janeTime.Format(dateTimeFormatWithTimeStamp)
	return []byte("\"" + timeString + "\""), nil
}
