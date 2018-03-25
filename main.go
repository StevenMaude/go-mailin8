package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type mail struct {
	From    string `json:"f"`
	Subject string `json:"s"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

type msg struct {
	UID string `json:"uid"`
}

type inbox struct {
	Msgs []msg `json:"msgs"`
}

func getMail(latestMsg msg) error {
	msgURL := "https://getnada.com/api/v1/messages/" + latestMsg.UID
	fmt.Println("Retrieving URL:", msgURL)

	resp, err := http.Get(msgURL)
	if err != nil {
		return err
	}

	// TODO: move out display of mail from getting mail.
	defer resp.Body.Close()

	mailMessage := mail{}
	err = json.NewDecoder(resp.Body).Decode(&mailMessage)
	if err != nil {
		return err
	}

	fmt.Println("\nFrom   :", mailMessage.From)
	fmt.Println("Subject:", mailMessage.Subject)
	fmt.Println("Plain text:")
	fmt.Println(mailMessage.Text)

	fmt.Println("HTML:")
	fmt.Println(mailMessage.HTML)

	return nil
}

func getInbox(address string) (inbox, error) {
	webInboxURL := "https://getnada.com/api/v1/inboxes/" + address
	fmt.Println("Retrieving URL:", webInboxURL)

	addressInbox := inbox{}
	resp, err := http.Get(webInboxURL)
	if err != nil {
		return addressInbox, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&addressInbox)
	// No need for error check here as we return mbxDetails and err whether
	// we have an error or not.
	return addressInbox, err
}

func main() {
	// TODO: consider allow to retrieve more than one message.
	if len(os.Args) != 2 {
		fmt.Println("Usage: mailin8 <address>")
		os.Exit(1)
	}

	address := os.Args[1]
	addressInbox, err := getInbox(address)
	if err != nil {
		fmt.Println("failed to get message ID:", err)
		os.Exit(1)
	}

	numberMsgs := len(addressInbox.Msgs)
	if numberMsgs == 0 {
		fmt.Println("no messages in inbox")
		os.Exit(0)
	}
	fmt.Println("Found", numberMsgs, "messages")

	latestMsg := addressInbox.Msgs[0]
	err = getMail(latestMsg)
	if err != nil {
		fmt.Println("failed to get mail:", err)
		os.Exit(1)
	}
}
