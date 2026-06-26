package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type BrevoRecipient struct {
	Email string `json:"to"`
}

type BrevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BrevoPayload struct {
	Sender      BrevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

// SendBrevoEmail drops mail packets out via standard HTTP requests safely.
func SendBrevoEmail(toEmail string, ticketID int, status, message string) {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		fmt.Println("LOG: Brevo API key missing from host environment variables.")
		return
	}

	subjectLine := fmt.Sprintf("Update on your IT Support Ticket #%d [%s]", ticketID, status)
	if ticketID == 0 {
		subjectLine = "Lordis Account System Security Alert"
	}

	payload := BrevoPayload{
		Sender:  BrevoSender{Name: "Lordis help", Email: "oluwadamilareoshodi@gmail.com"},
		To:      []BrevoRecipient{{Email: toEmail}},
		Subject: subjectLine,
		HTMLContent: fmt.Sprintf(`
			<html>
			<body style="font-family: sans-serif; color: #333;">
				<h2 style="color: #4f46e5;">Lordis Notification Dispatcher</h2>
				<p>Hello,</p>
				<p>A system process updated your dashboard records status to: <strong style="color: #2563eb;">%s</strong></p>
				<div style="background: #f3f4f6; padding: 15px; border-left: 4px solid #4f46e5; margin: 15px 0;">
					%s
				</div>
				<p style="font-size: 12px; color: #666;">This is an automated operational transmission from the CitiData Centre infrastructure tier.</p>
			</body>
			</html>`, status, message),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR: Failed connecting to Brevo cloud endpoint: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Printf("SUCCESS: Firewalled outbound API mail routed flawlessly to %s\n", toEmail)
	} else {
		fmt.Printf("ERROR: Brevo endpoint rejected transmission payload with code: %d\n", resp.StatusCode)
	}
}