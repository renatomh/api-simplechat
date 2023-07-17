package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Seeding the random generator
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	// Initializing the string builder
	var sb strings.Builder
	k := len(alphabet)

	// Adding random chars
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		// Writing to the string builder
		sb.WriteByte(c)
	}

	// Finally, we return the string built
	return sb.String()
}

// RandomUsername generates a random username
func RandomUsername() string {
	// Username format will be like "jack.doe"
	return (RandomString(int(RandomInt(3, 6))) + "." + RandomString(int(RandomInt(3, 6))))
}

// RandomStatus generates a random contact status
func RandomStatus() string {
	statuses := []string{"Pending", "Accepted", "Rejected"}
	return statuses[len(statuses)]
}
