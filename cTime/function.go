package cTime

import (
	"time"
)

func MskLocation() (*time.Location, error) {
	return time.LoadLocation(Moscow)
}

// возвращает московское время. Если была получена ошибка
// в получении московской временной зоны, то добавляет к UTC +3 часа
func Msk() time.Time {

	location, err := MskLocation()
	if err != nil {
		return time.Now().UTC().Add(time.Hour * 3)
	}

	return time.Now().In(location)
}

// аналог time.Since(time.Now()) для Msk()
func Since(start time.Time) time.Duration {
	return Msk().Sub(start)
}
