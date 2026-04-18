package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeOnly time.Time

var timeFormats = []string{
	"15:04",              // Часы и минуты
	"15:04:05",           // Часы, минуты и секунды
	"15:04:05.000",       // Часы, минуты, секунды и миллисекунды
	"15:04:05.000000",    // Часы, минуты, секунды и микросекунды
	"15:04:05.000000000", // Часы, минуты, секунды и наносекунды
}

func ParseTimeFromString(strValue string) (*time.Time, error) {
	loc, strValueRest, err := ParseTimezoneExtended(strValue)
	if err != nil {
		return nil, err
	}

	strValueRest = strings.TrimSpace(strValueRest)

	var parsed time.Time
	for _, f := range timeFormats {
		parsed, err = time.ParseInLocation(f, strValueRest, loc)
		if err == nil {
			return &parsed, nil
		}
	}

	return nil, errors.New("cannot parse time from " + strValue)
}

// ParseTimezoneExtended пытается разобрать таймзону из строк формата +HHMM или +HH:MM
func ParseTimezoneExtended(strValue string) (*time.Location, string, error) {
	loc := time.Local
	s := strValue

	// Попытка извлечь таймзону в формате +HH:MM (6 символов)
	if len(s) >= 6 {
		tzCandidate := s[len(s)-6:]
		// Формат: _HH:MM или +HH:MM или -HH:MM
		if (tzCandidate[0] == '+' || tzCandidate[0] == '-' || tzCandidate[0] == ' ') &&
			isDigitString(tzCandidate[1:3]) && tzCandidate[3] == ':' &&
			isDigitString(tzCandidate[4:]) {

			offsetLoc, err := parseTimezoneWithColon(tzCandidate) // реализуем новый парсер
			if err == nil {
				loc = offsetLoc
				s = s[:len(s)-6]
				return loc, s, nil
			}
		}
	}

	// Если не получилось, пробуем старый формат +HHMM (5 символов)
	if len(s) >= 5 {
		tzCandidate := s[len(s)-5:]
		if (tzCandidate[0] == '+' || tzCandidate[0] == '-') &&
			isDigitString(tzCandidate[1:]) {

			offsetLoc, err := parseTimezone(tzCandidate)
			if err == nil {
				loc = offsetLoc
				s = s[:len(s)-5]
				return loc, s, nil
			}
		}
	}

	// Не нашли таймзону
	return loc, s, nil
}

// parseTimezoneWithColon парсит таймзону формата +HH:MM или -HH:MM
func parseTimezoneWithColon(tz string) (*time.Location, error) {
	sign := tz[0]
	hhStr := tz[1:3]
	mmStr := tz[4:6]

	hh, err := strconv.Atoi(hhStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hours from timezone %s", tz)
	}
	mm, err := strconv.Atoi(mmStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse minutes from timezone %s", tz)
	}

	offset := hh*3600 + mm*60
	if sign == '-' {
		offset = -offset
	}
	return time.FixedZone(fmt.Sprintf("UTC%s", tz), offset), nil
}

// Парсим таймзону вида +HHMM/-HHMM
func parseTimezone(tz string) (*time.Location, error) {
	if len(tz) != 5 {
		return nil, errors.New("unknown time zone format")
	}
	sign := tz[0]
	hhStr := tz[1:3]
	mmStr := tz[3:5]

	hh, err := strconv.Atoi(hhStr)
	if err != nil {
		return nil, fmt.Errorf("timezone hours parsing error: %w", err)
	}
	mm, err := strconv.Atoi(mmStr)
	if err != nil {
		return nil, fmt.Errorf("timezone minutes parsing error: %w", err)
	}

	offset := hh*3600 + mm*60
	if sign == '-' {
		offset = -offset
	}
	loc := time.FixedZone(fmt.Sprintf("UTC%s", tz), offset)
	return loc, nil
}

// String возвращает строковое представление TimeOnly
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal TimeOnly) String() string {
	return time.Time(thisVal).Format("15:04:05")
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *TimeOnly) ToString() string {
	return thisVal.String()
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal TimeOnly) Value() (driver.Value, error) {
	return thisVal, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		return errors.New("null value cannot be scanned into TimeOnly")
	}

	strValue, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	parsed, err := ParseTimeFromString(strValue)
	if err != nil {
		return err
	}
	*thisVal = TimeOnly(*parsed)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal TimeOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(thisVal).Format(time.TimeOnly))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *TimeOnly) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		return errors.New("empty time string in JSON")
	}
	s := strings.Trim(sd, "\"")
	val, err := ParseTimeFromString(s)
	if err != nil {
		return err
	}
	*thisVal = TimeOnly(*val)
	return nil
}
