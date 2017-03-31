package main

// StringInSlice Go has no generics.
func StringInSlice(s string, l []string) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == s {
			return true
		}
	}
	return false
}
