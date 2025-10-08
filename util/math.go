package util

import "math"

func Round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// StepsToDistance 步数转换为距离（平均身高167＝企业里男170cm，比例73.1%，女158，比例26.9%）
func StepsToDistance(steps int64) float64 {
	return float64((float64(steps) * 167.0 * (45.0 / 100.0) / 100.0))
	// return float64((steps * 50 * (165 - 132) / 2700))
	// calories = 70 * distance  * 1036 / 1000 / 1000;
	// return float64(steps) * 0.6
}
