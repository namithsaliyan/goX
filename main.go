package main

import (
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "net/smtp"
    "os"
    "strings"
)

// Configuration structure to hold SMTP settings
type Config struct {
    SMTPUser   string
    SMTPPass   string
    AdminEmail string
    CCEmails   string
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

// Function to initialize the configuration from environment variables
func initConfig() {
    config = Config{
        SMTPUser:   os.Getenv("SMTP_USER"),
        SMTPPass:   os.Getenv("SMTP_PASS"),
        AdminEmail: os.Getenv("ADMIN_EMAIL"),
        CCEmails:   os.Getenv("CC_EMAILS"),
    }
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

func handleContactForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var msg Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        log.Printf("Failed to decode JSON: %v", err)
        return
    }

    if err := sendEmail(msg); err != nil {
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        log.Printf("Failed to send email: %v", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"status": "success", "message": "Thank you for your message! We will get back to you soon."}
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to send response", http.StatusInternalServerError)
        log.Printf("Failed to encode JSON response: %v", err)
    }
}

func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    initConfig()

    mux := http.NewServeMux()
    mux.HandleFunc("/submit", handleContactForm)

    handler := withCORS(mux)

    log.Println("Server started at :8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}
