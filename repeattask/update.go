package repeattask

import (
	"time"
)

func UpdateDate(date, repeat string) (string, error) {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().UTC().Location())

	switch repeat[0] {
	case 'd':
		next, err := RepeatD(now, date, repeat)
		if err != nil {
			return "", err
		}
		return next, nil
	case 'y':
		next, err := RepeatY(now, date, repeat)
		if err != nil {
			return "", err
		}
		return next, nil
	case 'w':
		next, err := RepeatW(now, date, repeat)
		if err != nil {
			return "", err
		}
		return next[0], nil
	}

	return "", nil
}
