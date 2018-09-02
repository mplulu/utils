package utils

import (
	"math/rand"
)

type ByInt64 []int64

func (a ByInt64) Len() int      { return len(a) }
func (a ByInt64) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByInt64) Less(i, j int) bool {
	return a[i] < a[j]
}

type ByInt []int

func (a ByInt) Len() int      { return len(a) }
func (a ByInt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByInt) Less(i, j int) bool {
	return a[i] < a[j]
}

func ShuffleStringSlice(slice []string) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func CloneStringSlice(slice []string) []string {
	newSlice := make([]string, len(slice), len(slice))
	for index, element := range slice {
		newSlice[index] = element
	}
	return newSlice
}

func ShuffleBoolSlice(slice []bool) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func RemoveStringElement(slice []string, element string) []string {
	for i, inSlice := range slice {
		if inSlice == element {
			copy(slice[i:], slice[i+1:])
			slice[len(slice)-1] = ""
			return slice[:len(slice)-1]
		}
	}
	return slice
}
