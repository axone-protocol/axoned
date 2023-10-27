package util

// Map applies the given function to each element of the given slice and returns a new slice with the results.
func Map[T, M any](s []T, f func(T) M) []M {
	m := make([]M, len(s))
	for i, v := range s {
		m[i] = f(v)
	}
	return m
}
