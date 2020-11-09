package main

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func main() {
	for {
		userInfo := getRandomUser()
		fmt.Println(userInfo)

		sendRequest(userInfo)
		time.Sleep(2 * time.Second)
		validEmail(userInfo)
	}
}

func getRandomUser() UserInfo {
	resultRand := rand()
	gender := 0
	if resultRand {
		gender = 1
	}

	fmt.Println()
	email := randomdata.Email()
	email = strings.Split(email, "@")[0] + "@yopmail.com"
	userInfo := UserInfo{
		FirstName: randomdata.FirstName(gender),
		LastName:  randomdata.LastName(),
		Email:     email,
		ZipCode:   "42000",
	}
	return userInfo
}

func rand() bool {
	c := make(chan struct{})
	close(c)
	select {
	case <-c:
		return true
	case <-c:
		return false
	}
}

func validEmail(info UserInfo) {
	fmt.Println("valid email...")

	email := strings.Split(info.Email, "@")[0]

	nbEmail, err := getNumEmail(email)
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		validEmail(info)
		return
	}

	link, err := getLink(email, nbEmail)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("link", link)
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.Status == "200 OK" {
		fmt.Println("Signature OK")
	} else {
		fmt.Println(resp.Status)
	}
}

func getNumEmail(email string) (string, error) {
	fmt.Println("get num email...")

	cmd := exec.Command("/home/lambert/yogo_linux_amd64", "inbox", "list", email, "10")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.HasSuffix(line, "Nous sommes pour") {
			return strings.Split(line, " ")[0], nil
		}
	}
	return "", errors.New("email not found")
}

func getLink(email, nbEmail string) (string, error) {
	fmt.Println("get link...", email, nbEmail)

	cmd := exec.Command("/home/lambert/yogo_linux_amd64", "inbox", "show", email, nbEmail)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Bot probably detected :/")
		return "", err
	}
	words := strings.Split(string(output), " ")
	for _, word := range words {
		if strings.HasPrefix(word, "https://a.noussommespour.fr/inscription") {
			return word, nil
		}
	}
	return "", errors.New("link not found")
}

type UserInfo struct {
	FirstName string
	LastName  string
	Email     string
	ZipCode   string
}

func sendRequest(userInfo UserInfo) {
	fmt.Println("send Request...")
	url := "https://noussommespour.fr/wp-admin/admin-ajax.php"

	str := fmt.Sprintf("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"post_id\"\r\n\r\n2\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_id\"\r\n\r\n9352cee\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[first_name]\"\r\n\r\n%s\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[last_name]\"\r\n\r\n%s\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[email]\"\r\n\r\n%s\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[contact_phone]\"\r\n\r\n\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[location_zip]\"\r\n\r\n%s\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"form_fields[agir_referer]\"\r\n\r\n\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"action\"\r\n\r\nelementor_pro_forms_send_form\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"referrer\"\r\n\r\nhttps://noussommespour.fr/\r\n-----011000010111000001101001--\r\n", userInfo.FirstName, userInfo.LastName, userInfo.Email, userInfo.ZipCode)

	payload := strings.NewReader(str)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cookie", "__cfduid=db7720dcb09c960da705547abcd63deae1604920044; agir_id=9ae4a99e-90db-4cf5-9f6f-197774ad3536; agir_email=marinelepen@yopmail.com; agir_location_zip=75003; agir_first_name=Marine; agir_last_name=Lepen")
	req.Header.Add("authority", "noussommespour.fr")
	req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")
	req.Header.Add("content-type", "multipart/form-data; boundary=---011000010111000001101001")
	req.Header.Add("origin", "https://noussommespour.fr")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://noussommespour.fr/")
	req.Header.Add("accept-language", "en-US,en;q=0.9,fr;q=0.8")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
}