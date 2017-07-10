package main

import "testing"
import "time"

func TestCalculateExpireDate(t *testing.T) {
	// get some sample data
	var expireDate = time.Now()
	var futureDate = expireDate.Add(time.Hour * time.Duration(10))

	newDate := calculateExpireDate(expireDate, "10H")
	if !newDate.Equal(futureDate) {
		t.Errorf("Got %s expected %s", newDate.String(), futureDate.String())
	}

	expireDate = time.Now()
	futureDate = expireDate.AddDate(0, 0, 2*7)

	newDate = calculateExpireDate(expireDate, "2w")
	if !newDate.Equal(futureDate) {
		t.Errorf("Got %s expected %s", newDate.String(), futureDate.String())
	}
}
