package helpers

import (
	"os"
	"reflect"
)

func ReverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func SaveStringToFile(filePath string, str string) error {
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(str)

	if err != nil {
		return err
	}

	return nil
}

func SliceContainsString(slice []string, str string) bool {
	for _, a := range slice {
		if a == str {
			return true
		}
	}
	return false
}
