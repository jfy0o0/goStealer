package gsconv

type (
	// errorStack is the interface for Stack feature.
	errorStack interface {
		Error() string
		Stack() string
	}
)

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}

	// StructTagPriority defines the default priority tags for Map*/Struct* functions.
	// Note, the "gsconv", "param", "params" tags are used by old version of package.
	// It is strongly recommended using short tag "c" or "p" instead in the future.
	StructTagPriority = []string{"gsconv", "param", "params", "c", "p", "json"}
)

type doConvertInput struct {
	FromValue  interface{}   // Value that is converted from.
	ToTypeName string        // Target value type name in string.
	ReferValue interface{}   // Referred value, a value in type `ToTypeName`.
	Extra      []interface{} // Extra values for implementing the converting.
}
