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

func main() {
	// Init Seed
	rand.Seed(time.Now().Unix())
	// Ask for input
	fmt.Println("Type in three uncommon words to use in your password. These words can", "include your favorite band, snack, etc! Please enter words with a space inbetween")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fields := strings.Fields(text)
	// punctuationStr := "!@#$%^&*(){}"
	// punct := punctuationStr[rand.Intn(len(punctuationStr))]

	// fields[0] = fields[0][:len(fields)/2] + string(punct) + fields[1][len(fields)/2:]

	passwd := ""

	for _, str := range fields {
		passwd += str
	}

	sum := sha1.Sum([]byte(passwd))
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
