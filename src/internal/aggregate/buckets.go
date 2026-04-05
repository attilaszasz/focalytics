package aggregate

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func yearBucketKey(value time.Time) string {
	return fmt.Sprintf("%04d", value.Year())
}

func dayBucketKey(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}

func focalLengthBucket(value float64) (string, string) {
	rounded := int(math.Round(value * 10))
	return fmt.Sprintf("%06d", rounded), fmt.Sprintf("%smm", formatScaledDecimal(rounded, 10))
}

func apertureBucket(value float64) (string, string) {
	rounded := int(math.Round(value * 10))
	return fmt.Sprintf("%06d", rounded), fmt.Sprintf("f/%s", formatScaledDecimal(rounded, 10))
}

func shutterSpeedBucket(value float64) (string, string) {
	micros := int(math.Round(value * 1_000_000))
	label := fmt.Sprintf("%ss", formatDurationSeconds(value))
	return fmt.Sprintf("%012d", micros), label
}

func isoBucket(value int) (string, string) {
	return fmt.Sprintf("%06d", value), fmt.Sprintf("ISO %d", value)
}

func formatScaledDecimal(value, scale int) string {
	integer := value / scale
	fraction := value % scale
	if fraction == 0 {
		return strconv.Itoa(integer)
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%d.%d", integer, fraction), "0"), ".")
}

func formatDurationSeconds(value float64) string {
	if value > 0 && value < 1 {
		reciprocal := math.Round(1 / value)
		if reciprocal > 0 {
			return fmt.Sprintf("1/%d", int(reciprocal))
		}
	}
	if math.Mod(value, 1) == 0 {
		return strconv.Itoa(int(value))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", value), "0"), ".")
}
