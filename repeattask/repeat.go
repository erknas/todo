package repeattask

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zeze322/todo/lib"
)

func RepeatD(now time.Time, date string, repeat string) (string, error) {
	t, err := time.Parse(lib.Layout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date")
	}

	repeat = strings.ReplaceAll(repeat, " ", "")

	days, err := strconv.Atoi(repeat[1:])
	if err != nil {
		return "", fmt.Errorf("bad day value")
	}

	if days < 1 || days > 400 {
		return "", fmt.Errorf("invalid day change")
	}

	if t == now && days == 1 {
		return now.Format(lib.Layout), nil
	}

	t = t.AddDate(0, 0, days)

	if t.Before(now) {
		for t.Before(now) {
			t = t.AddDate(0, 0, days)
		}
	}

	return t.Format(lib.Layout), nil
}

func RepeatY(now time.Time, date string, repeat string) (string, error) {
	t, err := time.Parse(lib.Layout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date")
	}

	next := t.AddDate(1, 0, 0)

	if next.Before(now) {
		return t.AddDate(now.Year()-t.Year(), 0, 0).Format(lib.Layout), nil
	}

	return next.Format(lib.Layout), nil
}

func RepeatW(now time.Time, date string, repeat string) ([]string, error) {
	t, err := time.Parse(lib.Layout, date)
	if err != nil {
		return nil, fmt.Errorf("invalid date")
	}

	repeat = strings.ReplaceAll(repeat, " ", "")
	days := strings.Split(repeat[1:], ",")

	next := make([]string, 0, len(days))

	for _, d := range days {
		day, err := strconv.Atoi(d)
		if err != nil {
			return nil, fmt.Errorf("bad day value")
		}

		if day < 1 || day > 7 {
			return nil, fmt.Errorf("invalid day")
		}

		if t.Before(now) {
			if day-int(now.Weekday()) > 0 {
				next = append(next, now.AddDate(0, 0, day-int(now.Weekday())).Format(lib.Layout))
			}
			if day-int(now.Weekday()) <= 0 {
				next = append(next, now.AddDate(0, 0, day-int(now.Weekday())+7).Format(lib.Layout))
			}
		}

		if !t.Before(now) {
			if day-int(t.Weekday()) > 0 {
				next = append(next, t.AddDate(0, 0, day-int(t.Weekday())).Format(lib.Layout))
			}
			if day-int(t.Weekday()) <= 0 {
				next = append(next, t.AddDate(0, 0, day-int(t.Weekday())+7).Format(lib.Layout))
			}
		}

	}

	sort.Strings(next)

	return next, nil
}

func RepeatM(now time.Time, date string, repeat string) ([]string, error) {
	t, err := time.Parse(lib.Layout, date)
	if err != nil {
		return nil, fmt.Errorf("invalid date")
	}

	daysMonths := strings.Split(repeat[2:], " ")

	if len(daysMonths) == 1 {
		days := strings.Split(daysMonths[0], ",")

		next := make([]string, 0, len(days))

		for _, d := range days {
			day, err := strconv.Atoi(d)
			if err != nil {
				return nil, fmt.Errorf("invalid day value")
			}

			if day > 31 || day < -2 {
				return nil, fmt.Errorf("invalid day")
			}

			if day == -1 {
				day = 0
			} else if day == -2 {
				day = -1
			}

			if t.Before(now) {
				if day < now.Day() {
					next = append(next, time.Date(now.Year(), now.Month()+1, day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
				}

				if day > now.Day() {
					if lastDayOfMonth(now) < day {
						next = append(next, time.Date(now.Year(), now.Month()+1, day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					} else {
						next = append(next, time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					}
				}
			}

			if !t.Before(now) {
				if day < t.Day() {
					next = append(next, time.Date(t.Year(), t.Month()+1, day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
				}

				if day > t.Day() {
					if lastDayOfMonth(t) < day {
						next = append(next, time.Date(t.Year(), t.Month()+1, day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					} else {
						next = append(next, time.Date(t.Year(), t.Month(), day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					}
				}

				sort.Strings(next)

				return next, nil
			}
		}

		return nil, nil
	}

	if len(daysMonths) == 2 {
		days := strings.Split(daysMonths[0], ",")
		months := strings.Split(daysMonths[1], ",")

		next := make([]string, 0, len(days))

		for _, d := range days {
			day, err := strconv.Atoi(d)
			if err != nil {
				return nil, fmt.Errorf("invalid day value")
			}

			if day > 31 || day < -2 {
				return nil, fmt.Errorf("invalid day")
			}

			if day == -1 {
				day = 0
			} else if day == -2 {
				day = -1
			}

			for _, m := range months {
				month, err := strconv.Atoi(m)
				if err != nil {
					return nil, fmt.Errorf("invalid month value")
				}

				if month > 12 || month < 0 {
					return nil, fmt.Errorf("invalid month")
				}

				t = time.Date(t.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)

				if t.Before(now) {
					if month <= int(now.Month()) {
						next = append(next, time.Date(now.Year()+1, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					} else {
						next = append(next, time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC).Format(lib.Layout))
					}
				} else {
					next = append(next, t.Format(lib.Layout))
				}
			}
		}

		sort.Strings(next)

		return next, nil
	}

	return nil, nil
}

func lastDayOfMonth(date time.Time) int {
	year, month, _ := date.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	return firstDay.AddDate(0, 1, -1).Day()
}
