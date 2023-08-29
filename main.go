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
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type mailHeader struct {
	From    string `json:"f"`
	Subject string `json:"s"`
}

type msgPreview struct {
	UID string `json:"uid"`
}

type msg struct {
	UID     string `json:"uid"`
	From    string `json:"f"`
	Subject string `json:"s"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

type inbox struct {
	MsgPreviews []msgPreview `json:"msgs"`
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 20 * time.Second,
	}
}

func getMailMsg(latestMsgPreview msgPreview) (msg, error) {
	msgURL := "https://inboxes.com/api/v2/message/" + latestMsgPreview.UID
	fmt.Println("Retrieving message URL:", msgURL)
	var m msg

	c := newHTTPClient()
	resp, err := c.Get(msgURL)
	if err != nil {
		return m, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return m, err
	}

	if latestMsgPreview.UID != m.UID {
		return m, fmt.Errorf("UID of preview %s and mail %s do not match", latestMsgPreview.UID, m.UID)
	}

	return m, nil
}

func getInbox(address string) (inbox, error) {
	webInboxURL := "https://inboxes.com/api/v2/inbox/" + address
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

	numberMsgs := len(addressInbox.MsgPreviews)
	if numberMsgs == 0 {
		fmt.Println("no messages in inbox")
		os.Exit(0)
	}
	fmt.Println("Found", numberMsgs, "messages")

	latestMsgID := addressInbox.MsgPreviews[0]
	msg, err := getMailMsg(latestMsgID)
	if err != nil {
		fmt.Println("failed to get message", err)
		os.Exit(1)
	}

	fmt.Println("From   :", msg.From)
	fmt.Println("Subject:", msg.Subject)
	fmt.Println("Text:")
	fmt.Println(msg.Text)
	fmt.Println("HTML:")
	fmt.Println(msg.HTML)
}
