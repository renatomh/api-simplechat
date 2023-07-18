package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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

// Here we can communicate with Random User's API
// https://randomuser.me/

const url string = "https://randomuser.me/api/"

// Defining the random user API response struct (created with https://mholt.github.io/json-to-go/)
type RandomUserAPIResponse struct {
	Results []RandomUser `json:"results"`
	Info    struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}
type RandomUser struct {
	Gender string `json:"gender"`
	Name   struct {
		Title string `json:"title"`
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Location struct {
		Street struct {
			Number int    `json:"number"`
			Name   string `json:"name"`
		} `json:"street"`
		City        string      `json:"city"`
		State       string      `json:"state"`
		Country     string      `json:"country"`
		Postcode    interface{} `json:"postcode"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
		Timezone struct {
			Offset      string `json:"offset"`
			Description string `json:"description"`
		} `json:"timezone"`
	} `json:"location"`
	Email string `json:"email"`
	Login struct {
		UUID     string `json:"uuid"`
		Username string `json:"username"`
		Password string `json:"password"`
		Salt     string `json:"salt"`
		Md5      string `json:"md5"`
		Sha1     string `json:"sha1"`
		Sha256   string `json:"sha256"`
	} `json:"login"`
	Dob struct {
		Date time.Time `json:"date"`
		Age  int       `json:"age"`
	} `json:"dob"`
	Registered struct {
		Date time.Time `json:"date"`
		Age  int       `json:"age"`
	} `json:"registered"`
	Phone string `json:"phone"`
	Cell  string `json:"cell"`
	ID    struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"id"`
	Picture struct {
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Thumbnail string `json:"thumbnail"`
	} `json:"picture"`
	Nat string `json:"nat"`
}

// QueryRandomUsersAPI queries the Random User API with the set number of results
func QueryRandomUsersAPI(results int) ([]RandomUser, error) {
	// Making the API request to get random users
	response, err := http.Get(fmt.Sprintf("%s?results=%d", url, results))
	if err != nil {
		fmt.Println(err)
		return []RandomUser{}, err
	}

	// Setting response body to be closed in the end
	defer response.Body.Close()
	// Reading the response body in a bytes slice
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return []RandomUser{}, err
	}

	// Here we can see the body response as string (JSON)
	//fmt.Println(string(body))

	// Now, we can parse the response with the defined struct
	data := RandomUserAPIResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return []RandomUser{}, err
	}

	// Returning the resulting data for the function call
	return data.Results, nil
}
