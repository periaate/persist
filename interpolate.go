package persist

import "math"

const (
	// DefaultSize is the default size of the set
	DefaultSize = 100

	// DefaultLoadFactor is the default load factor of the set
	DefaultLoadFactor = 0.75

	// GrowthFactorSmall is the growth factor for small sets
	GrowthFactorSmall = 20

	// GrowthFactorLarge is the growth factor for large sets
	GrowthFactorLarge = 1.2

	// SmallSetSize is the size at which the set is considered small
	SmallSetSize = 10000

	// LargeSetSize is the size at which the set is considered large
	LargeSetSize = 100000000
)

func NewInterpolateFunc(sizeS, sizeL uint64, resizeS, resizeL float64) func(uint64) uint64 {
	return func(value uint64) uint64 {
		if value >= sizeL {
			return uint64(resizeL * float64(value))
		}

		logL := math.Log10(float64(sizeL))
		logFactor := 0 - logL

		logValue := (math.Log10(float64(value)) - logL) / logFactor

		linearFactor := (resizeL - resizeS) / (1.0 - 0.0)
		res := resizeS + (linearFactor * (1.0 - logValue))
		return uint64(res * float64(value))
	}
}

func DefaultInterpolate() func(uint64) uint64 {
	return NewInterpolateFunc(SmallSetSize, LargeSetSize, GrowthFactorSmall, GrowthFactorLarge)
}
