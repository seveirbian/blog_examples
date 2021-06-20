package common

import (
	"strconv"
	"time"
)

var (
	Topic = "bian_test_" + strconv.Itoa(time.Now().YearDay())
	URL = "tcp://localhost:1883"
)
