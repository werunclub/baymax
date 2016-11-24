package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMaxTime(t *testing.T) {
	a := time.Now()
	b := a.AddDate(0, 0, -1)
	c := a.AddDate(0, 0, -100)
	d := a.AddDate(0, 0, 10)
	e := a.AddDate(0, 0, 99)

	maxTime := MaxTime(a, b, c, d, e)
	assert.Equal(t, maxTime, e)

	minTime := MinTime(a, b, c, d, e)
	assert.Equal(t, minTime, c)
}
