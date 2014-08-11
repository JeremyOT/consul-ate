package consul

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type CheckStatus string

const (
	CheckPass CheckStatus = "pass"
	CheckWarn CheckStatus = "warn"
	CheckFail CheckStatus = "fail"
	APIRoot               = "/v1"
)

type Client struct {
	address string
}

// Create a new client, ensuring that address is properly formatted
func NewClient(address string) *Client {
	if !strings.HasPrefix(address, "http") {
		address = "http://" + address
	}
	if address[len(address)-1] == '/' {
		address = address[:len(address)-1]
	}
	return &Client{address: address}
}

// Register a new service and return the ID of the new service and the ID of
// the generated check
func (c *Client) RegisterService(name, id string, tags []string, port int, check map[string]string) (serviceId, checkId string, err error) {
	if id == "" {
		id = name
	}
	body := map[string]interface{}{"Name": name, "ID": id}
	if tags != nil {
		body["Tags"] = tags
	}
	if port > 0 {
		body["Port"] = port
	}
	if check != nil {
		body["Check"] = check
	}
	data, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err := http.Post(c.address+path.Join(APIRoot, "agent/service/register"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", "", err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		return "", "", errors.New(resp.Status)
	}
	return id, fmt.Sprintf("service:%s", id), nil
}

// Deregister the specified service
func (c *Client) DeregisterService(id string) error {
	resp, err := http.Get(c.address + path.Join(APIRoot, "agent/service/deregister", id))
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

// Update the given check with the given status. Optionally, specify a note to be
// sent with the status
func (c *Client) UpdateCheck(id, note string, status CheckStatus) error {
	query := ""
	if note != "" {
		v := url.Values{}
		v.Set("note", note)
		query = "?" + v.Encode()
	}
	resp, err := http.Get(c.address + path.Join(APIRoot, "agent/check", string(status), id) + query)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

// Update the check with a "pass" and the given note at the given interval until quit is closed
func (c *Client) RegisterCheckHeartbeat(id, note string, interval time.Duration, quit chan int) {
	if err := c.UpdateCheck(id, note, CheckPass); err != nil {
		log.Println("Error updating check:", err)
	}
	clock := time.Tick(interval)
	for {
		select {
		case <-clock:
			if err := c.UpdateCheck(id, note, CheckPass); err != nil {
				log.Println("Error updating check:", err)
			}
		case <-quit:
			return
		}
	}
}

func (c *Client) Address() string {
	return c.address
}
