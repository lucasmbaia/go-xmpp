package utils

func ContainsString(value string, slice []string) (int, bool) {
  for index, str := range slice {
    if value == str {
      return index, true
    }
  }

  return 0, false
}
