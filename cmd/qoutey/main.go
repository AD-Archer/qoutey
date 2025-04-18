package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Config struct {
	SMTP struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	Email struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
	} `json:"email"`
	Schedule struct {
		Times []string `json:"times"`
	} `json:"schedule"`
	Quotes        []string `json:"quotes"`
	UsedQuotes    []string `json:"usedQuotes"`
	MaxRepetition int      `json:"maxRepetition"`
}

func main() {
	// Set up logging to file
	setupLogging()

	log.Println("Starting Quote Emailer application")

	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Check if test mode is enabled
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("Test mode: Sending a quote immediately...")
		log.Println("Test mode: Sending a quote immediately...")
		sendQuote(config)
		return
	}

	// Send startup notification email
	log.Println("Sending startup notification email...")
	sendStartupNotification(config)

	// Set up cron scheduler
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "CRON: ", log.LstdFlags))))

	// Use the schedule times from config
	for _, timeStr := range config.Schedule.Times {
		// Convert time string (e.g. "7:00") to cron expression (e.g. "0 7 * * *")
		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid time format: %s. Expected format: HH:MM", timeStr)
			continue
		}

		hour := parts[0]
		minute := parts[1]
		cronExpression := fmt.Sprintf("%s %s * * *", minute, hour)

		log.Printf("Setting up schedule: %s (cron: %s)", timeStr, cronExpression)

		// Create a closure to capture the current timeStr value
		scheduleTime := timeStr
		c.AddFunc(cronExpression, func() {
			log.Printf("Scheduled task: Sending quote at %s", scheduleTime)
			sendQuote(config)
		})
	}

	// Start the scheduler
	c.Start()

	// Create a readable representation of the schedule times
	scheduleDisplay := strings.Join(config.Schedule.Times, ", ")
	fmt.Printf("Quote emailer started. Running schedule at: %s\n", scheduleDisplay)
	fmt.Println("Application is running in the foreground. Keep this terminal window open.")
	fmt.Println("Check logs at ./qoutey.log")
	log.Printf("Quote emailer started. Running schedule at: %s", scheduleDisplay)

	// Keep the application running
	select {}
}

func setupLogging() {
	// Create log file
	logFile, err := os.OpenFile("qoutey.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to both file and console
	log.SetOutput(os.Stdout)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// Set log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func loadConfig(filename string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default config
		config := &Config{
			SMTP: struct {
				Server   string `json:"server"`
				Port     int    `json:"port"`
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Server:   "smtp.example.com",
				Port:     587,
				Username: "your-email@example.com",
				Password: "your-password",
			},
			Email: struct {
				From    string   `json:"from"`
				To      []string `json:"to"`
				Subject string   `json:"subject"`
			}{
				From:    "quotes@example.com",
				To:      []string{"your-email@example.com"},
				Subject: "Your Daily Quote",
			},
			Schedule: struct {
				Times []string `json:"times"`
			}{
				Times: []string{"0 7 * * *", "0 12 * * *", "0 19 * * *"},
			},
			Quotes: []string{
				"The only way to do great work is to love what you do. - Steve Jobs",
				"Life is what happens when you're busy making other plans. - John Lennon",
				"The future belongs to those who believe in the beauty of their dreams. - Eleanor Roosevelt",
			},
			UsedQuotes:    []string{},
			MaxRepetition: 5, // Avoid repeating quotes until at least 5 others have been sent
		}

		// Save default config
		saveConfig(filename, config)
		return config, nil
	}

	// Read config file
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}

func saveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func sendQuote(config *Config) {
	// Select quote
	quote := selectQuote(config)

	// Format message
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"\r\n"+
		"%s\r\n",
		config.Email.To[0],
		config.Email.Subject,
		quote)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", config.SMTP.Username, config.SMTP.Password, config.SMTP.Server)
	smtpAddr := fmt.Sprintf("%s:%d", config.SMTP.Server, config.SMTP.Port)

	// Send email
	err := smtp.SendMail(
		smtpAddr,
		auth,
		config.Email.From,
		config.Email.To,
		[]byte(message),
	)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return
	}

	log.Printf("Successfully sent quote: %s", quote)

	// Save updated used quotes
	saveConfig("config.json", config)
}

func selectQuote(config *Config) string {
	// If all quotes have been used, reset
	if len(config.Quotes) <= len(config.UsedQuotes) {
		// Keep only the most recent quotes to avoid repetition
		if len(config.UsedQuotes) > config.MaxRepetition {
			config.UsedQuotes = config.UsedQuotes[len(config.UsedQuotes)-config.MaxRepetition:]
		}
	}

	// Find a quote that hasn't been used recently
	availableQuotes := []string{}
	for _, quote := range config.Quotes {
		isUsed := false
		for _, usedQuote := range config.UsedQuotes {
			if quote == usedQuote {
				isUsed = true
				break
			}
		}
		if !isUsed {
			availableQuotes = append(availableQuotes, quote)
		}
	}

	// If all quotes have been used recently, use the least recently used one
	if len(availableQuotes) == 0 {
		// Get the oldest used quote (first in the list)
		selectedQuote := config.UsedQuotes[0]
		// Remove it from used quotes and add it to the end
		config.UsedQuotes = append(config.UsedQuotes[1:], selectedQuote)
		return selectedQuote
	}

	// Select a random quote from available ones
	selectedQuote := availableQuotes[rand.Intn(len(availableQuotes))]
	// Add to used quotes
	config.UsedQuotes = append(config.UsedQuotes, selectedQuote)

	return selectedQuote
}

func sendStartupNotification(config *Config) {
	// Format message
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s - Startup Notification\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"\r\n"+
		"Qoutey service has started successfully at %s.\r\n"+
		"The quote service is running and scheduled for 7am, 12pm, and 7pm.\r\n",
		config.Email.To[0],
		config.Email.Subject,
		time.Now().Format("2006-01-02 15:04:05"))

	// Connect to SMTP server
	auth := smtp.PlainAuth("", config.SMTP.Username, config.SMTP.Password, config.SMTP.Server)
	smtpAddr := fmt.Sprintf("%s:%d", config.SMTP.Server, config.SMTP.Port)

	// Send email
	err := smtp.SendMail(
		smtpAddr,
		auth,
		config.Email.From,
		config.Email.To,
		[]byte(message),
	)

	if err != nil {
		log.Printf("Failed to send startup notification email: %v", err)
		return
	}

	log.Printf("Successfully sent startup notification email")
}
