// Copyright 2018 Steven Maude
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type mailHeader struct {
	From    string `json:"f"`
	Subject string `json:"s"`
}

type msg struct {
	UID string `json:"uid"`
}

type inbox struct {
	Msgs []msg `json:"msgs"`
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 20 * time.Second,
	}
}

func getMailHeader(latestMsg msg) (mailHeader, error) {
	msgURL := "https://getnada.com/api/v1/messages/" + latestMsg.UID
	fmt.Println("Retrieving message URL:", msgURL)
	var mh mailHeader

	c := newHTTPClient()
	resp, err := c.Get(msgURL)
	if err != nil {
		return mh, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&mh)
	if err != nil {
		return mh, err
	}

	return mh, nil
}

func getMailBody(latestMsg msg) (string, error) {
	htmlURL := "https://getnada.com/api/v1/messages/html/" + latestMsg.UID
	fmt.Println("Retrieving HTML", htmlURL)

	c := newHTTPClient()
	htmlResp, err := c.Get(htmlURL)
	if err != nil {
		return "", err
	}
	defer htmlResp.Body.Close()

	html, err := ioutil.ReadAll(htmlResp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}

func getInbox(address string) (inbox, error) {
	webInboxURL := "https://getnada.com/api/v1/inboxes/" + address
	fmt.Println("Retrieving URL:", webInboxURL)

	var addressInbox inbox
	c := newHTTPClient()
	resp, err := c.Get(webInboxURL)
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
		fmt.Printf("Usage: %v <address>\n", filepath.Base(os.Args[0]))
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

	latestMsgID := addressInbox.Msgs[0]
	mh, err := getMailHeader(latestMsgID)
	if err != nil {
		fmt.Println("failed to get message metadate", err)
		os.Exit(1)
	}

	mb, err := getMailBody(latestMsgID)
	if err != nil {
		fmt.Println("failed to get mail:", err)
		os.Exit(1)
	}

	fmt.Println("From   :", mh.From)
	fmt.Println("Subject:", mh.Subject)
	fmt.Println("HTML:")
	fmt.Println(mb)
}
