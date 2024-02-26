package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func AddString(set map[string]bool, s string) {
	set[s] = true
}

// Function to get unique values from the set
func GetUniqueValues(set map[string]bool) []string {
	uniqueValues := make([]string, 0, len(set))
	for key := range set {
		uniqueValues = append(uniqueValues, key)
	}
	return uniqueValues
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func FilterArray(a, b []string) []string {
	// Create a map to store elements of array b
	bMap := make(map[string]bool)
	for _, val := range b {
		bMap[val] = true
	}

	// Initialize an empty result slice
	result := []string{}

	// Iterate through array a
	for _, val := range a {
		// If the element is not present in b, add it to the result
		if !bMap[val] {
			result = append(result, val)
		}
	}

	return result
}
