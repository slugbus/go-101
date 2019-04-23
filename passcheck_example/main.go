package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
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

func main() {
	// Init Seed
	rand.Seed(time.Now().Unix())
	// Ask for input
	fmt.Println("Type in three uncommon words to use in your password.",
		"These words can include your favorite band, snack, etc!",
		"Please enter words with a space in between")

	// Parse the fields from the user and join them.
	fields := strings.Fields(grabInput())
	input := joinAll(fields)

	sum := sha1.Sum([]byte(input))
	values := strings.ToUpper(fmt.Sprintf("%x", sum))

	// resp, err := http.Get(fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", "C8FED"))
	resp, err := http.Get(fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", values[:5]))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	reader2 := bufio.NewScanner(strings.NewReader(bodyStr))

	for reader2.Scan() {
		crackedPasswd := strings.Split(reader2.Text(), ":")[0]

		if crackedPasswd[len(crackedPasswd)-10:] == values[len(values)-10:] {
			fmt.Println("Your password has been cracked!")
			os.Exit(1)
		}

	}

	fmt.Println("Your password has not been cracked! Good choice!")
}
