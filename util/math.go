package util

import "math"

func Round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// StepsToDistance 步数转换为距离
func StepsToDistance(steps int) float64 {
	return float64((steps * 50 * (165 - 132) / 2700))
	// calories = 70 * distance  * 1036 / 1000 / 1000;
	// return float64(steps) * 0.6
}
