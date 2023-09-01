package pagination

import "fmt"

func ConvertToString(a any) string {
	switch v := a.(type) {
	case int64:
		return fmt.Sprintf("%d", int(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ToPointer[T any](p T) *T {
	return &p
}
