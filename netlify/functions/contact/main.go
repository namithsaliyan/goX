package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Configuration structure to hold SMTP settings
type Config struct {
	SMTPUser   string `json:"smtp_user"`
	SMTPPass   string `json:"smtp_pass"`
	AdminEmail string `json:"admin_email"`
	CCEmails   string `json:"cc_emails"`
}

type Message struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

const (
	smtpServer  = "smtp.gmail.com"
	smtpPort    = "587"
	emailSubject = "Contact Form Submission"
)

var config Config

// Function to load the configuration from environment variables
func loadConfig() {
	config.SMTPUser = os.Getenv("SMTP_USER")
	config.SMTPPass = os.Getenv("SMTP_PASS")
	config.AdminEmail = os.Getenv("ADMIN_EMAIL")
	config.CCEmails = os.Getenv("CC_EMAILS")
}

// Function to send an email using SMTP
func sendEmail(msg Message) error {
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, smtpServer)

	emailTemplate := `
        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #e9ecef; }
                .container { width: 100%; max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; }
                .header { background-color: #007bff; color: #ffffff; padding: 20px; text-align: center; }
                .content { padding: 20px; }
                .footer { background-color: #f1f1f1; padding: 10px; text-align: center; color: #6c757d; }
                h1 { margin: 0; }
                p { margin: 0 0 10px; }
                .highlight { color: #007bff; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Contact Form Submission</h1>
                </div>
                <div class="content">
                    <p><strong class="highlight">Name:</strong> {{.Name}}</p>
                    <p><strong class="highlight">Email:</strong> {{.Email}}</p>
                    <p><strong class="highlight">Message:</strong><br>{{.Message}}</p>
                </div>
                <div class="footer">
                    <p>&copy; 2024 Your Company. All rights reserved.</p>
                </div>
            </div>
        </body>
        </html>
    `

	t, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return err
	}

	var body strings.Builder
	err = t.Execute(&body, msg)
	if err != nil {
		return err
	}

	headers := []string{
		"From: " + config.SMTPUser,
		"To: " + config.AdminEmail,
		"Cc: " + config.CCEmails,
		"Subject: " + emailSubject,
		"Content-Type: text/html; charset=\"utf-8\"",
	}

	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body.String()

	toAddresses := strings.Split(config.AdminEmail, ",")
	ccAddresses := strings.Split(config.CCEmails, ",")
	allRecipients := append(toAddresses, ccAddresses...)

	log.Printf("Sending email with subject: %s", emailSubject)
	log.Printf("Message:\n%s", message)

	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, config.SMTPUser, allRecipients, []byte(message))
	return err
}

// AWS Lambda handler function
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Invalid request method",
		}, nil
	}

	var msg Message
	if err := json.Unmarshal([]byte(request.Body), &msg); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Failed to parse request body",
		}, nil
	}

	if err := sendEmail(msg); err != nil {
		log.Printf("Failed to send email: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to send email",
		}, nil
	}

	response := map[string]string{"status": "success", "message": "Thank you for your message! We will get back to you soon."}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func main() {
	loadConfig()
	lambda.Start(handleRequest)
}
