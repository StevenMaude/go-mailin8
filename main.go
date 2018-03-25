package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type mail struct {
	Data struct {
		Subject string `json:"subject"`
		Parts   []struct {
			Body string `json:"body"`
		} `json:"parts"`
		Headers struct {
			From string `json:"from"`
		} `json:"headers"`
	} `json:"data"`
}

type publicMsg struct {
	ID string `json:"id"`
	To string `json:"to"`
}

type mailboxDetails struct {
	PublicMsgs []publicMsg `json:"messages"`
}

type msg struct {
	UID string `json:"uid"`
}

type inbox struct {
	Msgs []msg `json:"msgs"`
}

func getMailboxDetails(localPart string) (mailboxDetails, error) {
	webInboxURL := "https://www.mailinator.com/fetch_inbox?zone=public&to=" + localPart
	fmt.Println("Retrieving URL:", webInboxURL)

	mbxDetails := mailboxDetails{}
	resp, err := http.Get(webInboxURL)
	if err != nil {
		return mbxDetails, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&mbxDetails)
	// No need for error check here as we return mbxDetails and err whether
	// we have an error or not.
	return mbxDetails, err
}

func getMail(latestMsg publicMsg, cookies []*http.Cookie) error {
	msgURL := "https://www.mailinator.com/fetch_email?zone=public&msgid=" + latestMsg.ID
	fmt.Println("Retrieving URL:", msgURL)
	req, err := http.NewRequest("GET", msgURL, nil)
	if err != nil {
		return err
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := &http.Client{}

	mailResp, err := client.Do(req)
	defer mailResp.Body.Close()
	if err != nil {
		return err
	}

	mailMessage := mail{}
	err = json.NewDecoder(mailResp.Body).Decode(&mailMessage)
	if err != nil {
		return err
	}

	fmt.Println("\nFrom   :", mailMessage.Data.Headers.From)
	fmt.Println("Subject:", mailMessage.Data.Subject)
	fmt.Println("Plain text:")
	fmt.Println(mailMessage.Data.Parts[0].Body)

	if len(mailMessage.Data.Parts) == 2 {
		fmt.Println("HTML:")
		fmt.Println(mailMessage.Data.Parts[1].Body)
	}

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
		fmt.Println("Usage: mailin8 <local-part>")
		os.Exit(1)
	}

	address := os.Args[1]
	addressInbox, err := getInbox(address)
	if err != nil {
		fmt.Println("failed to get message ID:", err)
		os.Exit(1)
	}

	fmt.Println(addressInbox)
	numberMsgs := len(addressInbox.Msgs)
	if numberMsgs == 0 {
		fmt.Println("no messages in inbox")
		os.Exit(0)
	}
	fmt.Println(numberMsgs)

	// latestMsg := mbxDetails.PublicMsgs[numberMsgs-1]

	//	err = getMail(latestMsg, cookies)
	//	if err != nil {
	//		fmt.Println("failed to get mail:", err)
	//		os.Exit(1)
	//	}
}
