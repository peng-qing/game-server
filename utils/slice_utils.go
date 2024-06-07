package utils

import "GameServer/types"

func ConvertIntSliceToInt64s[T types.Int](arr []T) []int64 {
	result := make([]int64, len(arr))

	for i, v := range arr {
		result[i] = int64(v)
	}

	return result
}

func ConvertUintSliceToUint64s[T types.Uint](arr []T) []uint64 {
	result := make([]uint64, len(arr))

	for i, v := range arr {
		result[i] = uint64(v)
	}

	return result
}

func ConvertFloatSliceToFloat64s[T types.Float](arr []T) []float64 {
	result := make([]float64, len(arr))
	for i, v := range arr {
		result[i] = float64(v)
	}
	return result
}
