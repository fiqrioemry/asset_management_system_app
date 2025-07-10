package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fiqrioemry/asset_management_system_app/server/config"
	gomail "gopkg.in/gomail.v2"
)

type EmailTemplate struct {
	Subject  string
	Template string
}

type EmailData struct {
	UserName    string
	Email       string
	ResetLink   string
	OTPCode     string
	ExpiryTime  string
	AppName     string
	SupportURL  string
	CompanyName string
}

// Email templates
var emailTemplates = map[string]EmailTemplate{
	"reset_password": {
		Subject: "Reset Your Password - {{.AppName}}",
		Template: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f8f9fa; padding: 20px; text-align: center; border-radius: 8px; margin-bottom: 30px; }
        .content { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .button { display: inline-block; background: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; font-weight: bold; margin: 20px 0; }
        .button:hover { background: #0056b3; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 14px; color: #666; text-align: center; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.AppName}}</h1>
        <p>Password Reset Request</p>
    </div>
    
    <div class="content">
        <h2>Hello {{.UserName}},</h2>
        
        <p>We received a request to reset your password for your {{.AppName}} account associated with <strong>{{.Email}}</strong>.</p>
        
        <p>Click the button below to reset your password:</p>
        
        <a href="{{.ResetLink}}" class="button">Reset Password</a>
        
        <div class="warning">
            <strong>Important:</strong> This link will expire in {{.ExpiryTime}}. If you didn't request this password reset, please ignore this email.
        </div>
        
        <p>If the button doesn't work, copy and paste this link into your browser:</p>
        <p style="word-break: break-all; background: #f8f9fa; padding: 10px; border-radius: 5px;">{{.ResetLink}}</p>
        
        <p>If you're having trouble, contact our support team at <a href="{{.SupportURL}}">{{.SupportURL}}</a></p>
        
        <p>Best regards,<br>The {{.CompanyName}} Team</p>
    </div>
    
    <div class="footer">
        <p>This email was sent to {{.Email}}. If you didn't request this, please ignore this email.</p>
        <p>&copy; {{.CompanyName}}. All rights reserved.</p>
    </div>
</body>
</html>`,
	},

	"otp_verification": {
		Subject: "Your OTP Code - {{.AppName}}",
		Template: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f8f9fa; padding: 20px; text-align: center; border-radius: 8px; margin-bottom: 30px; }
        .content { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .otp-box { background: #e3f2fd; border: 2px solid #2196f3; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0; }
        .otp-code { font-size: 32px; font-weight: bold; color: #1976d2; letter-spacing: 5px; margin: 10px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 14px; color: #666; text-align: center; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.AppName}}</h1>
        <p>OTP Verification</p>
    </div>
    
    <div class="content">
        <h2>Hello {{.UserName}},</h2>
        
        <p>Your One-Time Password (OTP) for {{.AppName}} is:</p>
        
        <div class="otp-box">
            <div class="otp-code">{{.OTPCode}}</div>
            <p>Enter this code to verify your identity</p>
        </div>
        
        <div class="warning">
            <strong>Important:</strong> This OTP will expire in {{.ExpiryTime}}. Don't share this code with anyone.
        </div>
        
        <p>If you didn't request this OTP, please ignore this email and contact our support team immediately.</p>
        
        <p>Best regards,<br>The {{.CompanyName}} Team</p>
    </div>
    
    <div class="footer">
        <p>This email was sent to {{.Email}}. If you didn't request this, please contact support.</p>
        <p>&copy; {{.CompanyName}}. All rights reserved.</p>
    </div>
</body>
</html>`,
	},

	"welcome": {
		Subject: "Welcome to {{.AppName}}!",
		Template: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; border-radius: 8px; margin-bottom: 30px; }
        .content { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .button { display: inline-block; background: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; font-weight: bold; margin: 20px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 14px; color: #666; text-align: center; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Welcome to {{.AppName}}!</h1>
        <p>Your account has been created successfully</p>
    </div>
    
    <div class="content">
        <h2>Hello {{.UserName}},</h2>
        
        <p>Welcome to {{.AppName}}! We're excited to have you on board.</p>
        
        <p>Your account (<strong>{{.Email}}</strong>) has been successfully created and is ready to use.</p>
        
        <a href="{{.SupportURL}}" class="button">Get Started</a>
        
        <p>If you have any questions or need help getting started, don't hesitate to reach out to our support team.</p>
        
        <p>Best regards,<br>The {{.CompanyName}} Team</p>
    </div>
    
    <div class="footer">
        <p>&copy; {{.CompanyName}}. All rights reserved.</p>
    </div>
</body>
</html>`,
	},
}

// SendTemplateEmail sends email using predefined templates
func SendTemplateEmail(templateName, toEmail string, data EmailData) error {
	tmpl, exists := emailTemplates[templateName]
	if !exists {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Set default values if not provided
	if data.AppName == "" {
		data.AppName = getEnvOrDefault("APP_NAME", "Asset Management System")
	}
	if data.CompanyName == "" {
		data.CompanyName = getEnvOrDefault("COMPANY_NAME", "Your Company")
	}
	if data.SupportURL == "" {
		data.SupportURL = getEnvOrDefault("SUPPORT_URL", "mailto:support@yourcompany.com")
	}

	// Parse and execute subject template
	subjectTmpl, err := template.New("subject").Parse(tmpl.Subject)
	if err != nil {
		return fmt.Errorf("failed to parse subject template: %w", err)
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	// Parse and execute body template
	bodyTmpl, err := template.New("body").Parse(tmpl.Template)
	if err != nil {
		return fmt.Errorf("failed to parse body template: %w", err)
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	subject := subjectBuf.String()
	htmlBody := bodyBuf.String()

	// Create plain text version (strip HTML tags)
	plainTextBody := stripHTML(htmlBody)

	return SendEmail(subject, toEmail, plainTextBody, htmlBody)
}

// SendResetPasswordEmail sends password reset email
func SendResetPasswordEmail(toEmail, userName, resetLink string, expiryDuration time.Duration) error {
	data := EmailData{
		UserName:   userName,
		Email:      toEmail,
		ResetLink:  resetLink,
		ExpiryTime: formatDuration(expiryDuration),
	}

	return SendTemplateEmail("reset_password", toEmail, data)
}

// SendOTPEmail sends OTP verification email
func SendOTPEmail(toEmail, userName, otpCode string, expiryDuration time.Duration) error {
	data := EmailData{
		UserName:   userName,
		Email:      toEmail,
		OTPCode:    otpCode,
		ExpiryTime: formatDuration(expiryDuration),
	}

	return SendTemplateEmail("otp_verification", toEmail, data)
}

// SendWelcomeEmail sends welcome email
func SendWelcomeEmail(toEmail, userName string) error {
	data := EmailData{
		UserName: userName,
		Email:    toEmail,
	}

	return SendTemplateEmail("welcome", toEmail, data)
}

// LoadTemplatesFromFile loads email templates from external files
func LoadTemplatesFromFile(templatesDir string) error {
	if templatesDir == "" {
		return nil // Use default templates
	}

	files, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return fmt.Errorf("failed to read template files: %w", err)
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".html")

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", file, err)
		}

		// You can extend this to parse subject from file headers
		emailTemplates[name] = EmailTemplate{
			Subject:  "{{.AppName}} - Notification",
			Template: string(content),
		}
	}

	return nil
}

// SendNotificationEmail - backward compatibility
func SendNotificationEmail(to string, from string, title string, message string) error {
	plainText := message
	html := fmt.Sprintf("<p>%s</p>", message)

	return SendEmail(title, to, plainText, html)
}

// SendEmail - original function with improvements
func SendEmail(subject, toEmail, plainTextBody, htmlBody string) error {
	m := gomail.NewMessage()
	from := getEnvOrDefault("USER_EMAIL", "noreply@yourcompany.com")

	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plainTextBody)
	m.AddAlternative("text/html", htmlBody)

	if err := config.MailDialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", toEmail, err)
	}

	return nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func stripHTML(html string) string {
	// Simple HTML tag removal for plain text
	html = strings.ReplaceAll(html, "<br>", "\n")
	html = strings.ReplaceAll(html, "<br/>", "\n")
	html = strings.ReplaceAll(html, "<br />", "\n")
	html = strings.ReplaceAll(html, "</p>", "\n\n")
	html = strings.ReplaceAll(html, "</div>", "\n")
	html = strings.ReplaceAll(html, "</h1>", "\n")
	html = strings.ReplaceAll(html, "</h2>", "\n")
	html = strings.ReplaceAll(html, "</h3>", "\n")

	// Remove all HTML tags
	for strings.Contains(html, "<") && strings.Contains(html, ">") {
		start := strings.Index(html, "<")
		end := strings.Index(html[start:], ">")
		if end == -1 {
			break
		}
		html = html[:start] + html[start+end+1:]
	}

	// Clean up extra whitespace
	html = strings.ReplaceAll(html, "\n\n\n", "\n\n")
	html = strings.TrimSpace(html)

	return html
}

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}
