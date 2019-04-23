package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func grabInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func joinAll(tFields []string) string {
	var sb strings.Builder

	for _, field := range tFields {
		// As per the doc, WriteString always
		// returns a nil error. But
		// go insists we "check" it.
		_, err := sb.WriteString(field)
		_ = err
	}

	return sb.String()
}

func beenPwned(hashedPassword, haveIBeenPwnedResp string) (bool, error) {
	// Create a scanner from the response.
	scanner := bufio.NewScanner(strings.NewReader(haveIBeenPwnedResp))

	// Scan the response.
	for scanner.Scan() {
		currentHash := scanner.Text()
		// the api returns the hash - (first five chars)
		// since a standard hash is 40 chars this means
		// that the hash is the fist 35 characters.
		currentHash = currentHash[:35]

		// Check to see if the first 10 chars are the same
		if hashedPassword[len(hashedPassword)-10:] == currentHash[len(currentHash)-10:] {
			return true, nil
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("could not scan text: %v", err)
	}

	return false, nil
}

func insertMiddle(s, set string) string {
	// Chose a random rune from the set.
	char := set[rand.Intn(len(s))]
	// Insert and return the random char into the middle of the string.
	return s[:len(s)/2] + string(char) + s[len(s)/2:]
}

func sha1HashAsString(data []byte) string {
	hash := sha1.Sum(data)
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

func queryHaveIBeenPwned(query string) (string, error) {
	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", query)
	resp, err := http.Get(url)
	// Check for 2 sources of error:
	// 0) Error from err var.
	// 1) Bad status code
	if err != nil {
		return "", fmt.Errorf("could not not hit api with query %q: %v", query, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response from api with query %q: %d (%v)",
			query,
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	// Otherwise read the body and return it
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read body response: %v", err)
	}
	return string(bodyBytes), nil
}

func main() {
	// Init the seed.
	rand.Seed(time.Now().Unix())

	// Ask for input
	fmt.Println("Type in three uncommon words to use in your password.",
		"These words can include your favorite band, snack, etc!",
		"Please enter words with a space in between")

	// Parse the fields from the user and join them.
	fields := strings.Fields(grabInput())
	input := joinAll(fields)

	// Insert a random punctuation mark.
	punctuation := "!@#$%^&*(){}<>?."
	input = insertMiddle(input, punctuation)

	// Hash input.
	hashedInput := sha1HashAsString([]byte(input))

	// And queryHaveIBeenPwned
	resp, err := queryHaveIBeenPwned(hashedInput[:5])
	if err != nil {
		log.Fatalln(err)
	}

	pwned, err := beenPwned(hashedInput, resp)

	if err != nil {
		log.Fatalln(err)
	}

	if pwned {
		fmt.Println("Your password has been cracked!")
		return
	}

	fmt.Println("Your password has not been cracked! Good choice!")
}
