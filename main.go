package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type pageData struct {
	Title          string
	Year           int
	Notice         string
	ContactSuccess bool
	ScheduleSent   bool
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/Images/", http.StripPrefix("/Images/", http.FileServer(http.Dir("Images"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("Images"))))
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/contact", contactHandler)
	mux.HandleFunc("/schedule", scheduleHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           loggingMiddleware(mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("House of Bounce portal running on http://localhost:%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func renderHome(w http.ResponseWriter, notice string, contactSuccess bool, scheduleSent bool) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		log.Printf("template parse error: %v", err)
		return
	}

	data := pageData{
		Title:          "House of Bounce | Bounce House Rentals in Maine",
		Year:           time.Now().Year(),
		Notice:         notice,
		ContactSuccess: contactSuccess,
		ScheduleSent:   scheduleSent,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Printf("template execution error: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderHome(w, "", false, false)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/#contact", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderHome(w, "We could not read your message. Please try again.", false, false)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || message == "" {
		renderHome(w, "Please fill in name, email, and message for contact requests.", false, false)
		return
	}

	log.Printf("CONTACT REQUEST | Name: %s | Email: %s | Phone: %s | Message: %s", name, email, phone, message)
	go sendEmail(
		"House of Bounce - New Contact Message from "+name,
		fmt.Sprintf("Name: %s\nEmail: %s\nPhone: %s\n\nMessage:\n%s", name, email, phone, message),
	)
	renderHome(w, "Thanks for reaching out. We will contact you shortly.", true, false)
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/#schedule", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderHome(w, "We could not read your scheduling request. Please try again.", false, false)
		return
	}

	clientName := strings.TrimSpace(r.FormValue("client_name"))
	eventDate := strings.TrimSpace(r.FormValue("event_date"))
	eventCity := strings.TrimSpace(r.FormValue("event_city"))
	equipment := strings.TrimSpace(r.FormValue("equipment"))
	notes := strings.TrimSpace(r.FormValue("notes"))

	if clientName == "" || eventDate == "" || eventCity == "" {
		renderHome(w, "Please complete name, event date, and city to request a booking.", false, false)
		return
	}

	log.Printf("SCHEDULING REQUEST | Client: %s | Date: %s | City: %s | Equipment: %s | Notes: %s", clientName, eventDate, eventCity, equipment, notes)
	go sendEmail(
		"House of Bounce - Scheduling Request from "+clientName,
		fmt.Sprintf("Client: %s\nEvent Date: %s\nCity: %s\nEquipment: %s\n\nNotes:\n%s", clientName, eventDate, eventCity, equipment, notes),
	)
	renderHome(w, "Your scheduling request was submitted. We will confirm availability soon.", false, true)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s from %s (%s)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func sendEmail(subject, body string) {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	to := os.Getenv("SMTP_TO")
	if from == "" || password == "" || to == "" {
		log.Println("EMAIL: SMTP env vars not configured, skipping")
		return
	}
	host := "smtp.gmail.com"
	auth := smtp.PlainAuth("", from, password, host)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: House of Bounce <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		to, from, subject, body,
	))
	if err := smtp.SendMail(host+":587", auth, from, []string{to}, msg); err != nil {
		log.Printf("EMAIL ERROR: %v", err)
	} else {
		log.Printf("EMAIL SENT: %s", subject)
	}
}
