package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Attendee struct {
	EmailAddress struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	Type string `json:"type"`
}

type Event struct {
	Subject           string     `json:"subject"`
	StartTime         string     `json:"start_time"`
	EndTime           string     `json:"end_time"`
	Attendees         []Attendee `json:"attendees"`
	Organizer         string     `json:"organizer"`
	MeetingTemplateId string     `json:"meeting_template_id"`
}

var apiAuth = map[string]any{
	"secret":   "[Appreg Secret]",
	"tenantID": "[Tenant ID]",
	"appID":    "[App ID]",
}

func createEventHandler(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	token := getBearerToken()

	joinWebUrl, audioConferencing := createEvent(event, token)
	if joinWebUrl == "" || joinWebUrl == "response not ok joinWebUrl" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": joinWebUrl + " was returned"})
		return
	}
	if audioConferencing == nil || audioConferencing == "response not ok audioConferencing" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Audio conferencing content not found in response"})
		return
	}

	body := sendInvite(event, token, audioConferencing, joinWebUrl)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"body":    body,
	})
}

func createEvent(event Event, token string) (string, any) {
	//URI := "https://graph.microsoft.com/v1.0/me/events"
	URI := "https://graph.microsoft.com/v1.0/users/" + event.Organizer + "/onlineMeetings"

	eventPayload := map[string]interface{}{
		"subject":           event.Subject,
		"startDateTime":     event.StartTime,
		"endDateTime":       event.EndTime,
		"meetingTemplateId": "[Meeting Template ID]",
	}

	body, err := json.Marshal(eventPayload)
	if err != nil {
		return "", ""
	}

	req, err := http.NewRequest("POST", URI, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return "", ""
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("%d\n", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Graph API error: %s\n", string(bodyBytes))
		return "", ""
	}

	var respBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Print(err)
		return "", ""
	}

	audioConferencing, ok := respBody["audioConferencing"]
	if !ok {
		return "response not ok audioConferencing", "response not ok audioConferencing"
	}

	joinWebUrl, ok := respBody["joinUrl"].(string)
	if !ok {
		return "response not ok joinUrl", "response not ok joinUrl"
	}

	return joinWebUrl, audioConferencing
}

func sendInvite(event Event, token string, audioConferencing interface{}, joinWebUrl string) string {
	//URI := "https://graph.microsoft.com/v1.0/me/events"
	URI := "https://graph.microsoft.com/v1.0/users/" + event.Organizer + "/events"

	// Format audioConferencing map into a readable string
	tollFreeNumber := audioConferencing.(map[string]interface{})["tollFreeNumber"]
	if tollFreeNumber == nil {
		tollFreeNumber = ""
	}
	tollNumber := audioConferencing.(map[string]interface{})["tollNumber"]
	if tollNumber == nil {
		tollNumber = ""
	}
	dialinUrl := audioConferencing.(map[string]interface{})["dialinUrl"]
	if dialinUrl == nil {
		dialinUrl = ""
	}

	bodyContent := `<br><b>Meeting Details:</b><br><br>
	<b>Join URL:</b> ` + joinWebUrl + `<br><br>
	<b>Audio Conferencing:</b><br>
	&nbsp&nbsp&nbsp <b>Toll Free Number:</b> ` + tollFreeNumber.(string) + `<br>
	&nbsp&nbsp&nbsp <b>Toll Number:</b> ` + tollNumber.(string) + `<br>
	&nbsp&nbsp&nbsp <b>Dial-in URL:</b> ` + dialinUrl.(string) + `<br>
	`

	eventPayload := map[string]interface{}{
		"subject": event.Subject,
		"start": map[string]string{
			"dateTime": event.StartTime,
			"timeZone": "UTC",
		},
		"end": map[string]string{
			"dateTime": event.EndTime,
			"timeZone": "UTC",
		},
		"isOnlineMeeting": false,
		"attendees":       event.Attendees,
		"body": map[string]string{
			"contentType": "HTML",
			"content":     bodyContent,
		},
	}

	reqBody, err := json.Marshal(eventPayload)
	if err != nil {
		return "failed to marshal event payload: " + err.Error()
	}

	req, err := http.NewRequest("POST", URI, io.NopCloser(bytes.NewReader(reqBody)))
	if err != nil {
		return "failed to create HTTP request: " + err.Error()
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "failed to send request to Graph API: " + err.Error()
	}
	defer resp.Body.Close()

	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Print(err)
		return ""
	}

	respBodyArr, ok := respBody.(map[string]interface{})["body"].(map[string]interface{})["content"]
	if !ok {
		return "response body is not a valid map"
	}

	return respBodyArr.(string)
}

func getBearerToken() string {
	body := "client_id=" + apiAuth["appID"].(string) +
		"&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default" +
		"&client_secret=" + apiAuth["secret"].(string) +
		"&grant_type=client_credentials"

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/"+apiAuth["tenantID"].(string)+"/oauth2/v2.0/token", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println("Error creating HTTP Request, error is: ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to token API, error is: ", err)

	}
	defer resp.Body.Close()

	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		fmt.Print(err)
		return ""
	}

	return respBody.(map[string]interface{})["access_token"].(string)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/create-event", createEventHandler)
	r.Run(":8080")
}
