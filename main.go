package main

import (
	"fmt"
	"log"
	"net"
	"net/smtp" // https://www.geeksforgeeks.org/sending-email-using-smtp-in-golang/
	"time"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	// main page FRONTEND
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter).
		AddInputField("From:", "", 25, nil, nil).
		AddPasswordField("Password:", "", 25, '*', nil). // passwd is google application password not the actual
		// https://support.google.com/accounts/answer/185833?hl=en <<<<------ pls see
		AddInputField("To:", "", 25, nil, nil).
		AddInputField("Subject:", "", 120, nil, nil).
		AddInputField("Mail-Body:", "", 120, nil, nil)

		// button configurations and what goes behind them kinda like their backend
	form.AddButton("SEND!", func() {
		emailItem := form.GetFormItem(0).(*tview.InputField)
		passwordItem := form.GetFormItem(1).(*tview.InputField)
		toEMAIL := form.GetFormItem(2).(*tview.InputField)
		subjectMAIL := form.GetFormItem(3).(*tview.InputField)
		bodyMAIL := form.GetFormItem(4).(*tview.InputField)

		// how to put gmail picture side by side?
		email := emailItem.GetText()
		password := passwordItem.GetText()
		toEmailId := toEMAIL.GetText()
		subject := subjectMAIL.GetText()
		mailbody := bodyMAIL.GetText()

		err := sendLoginNotification(email, password, toEmailId, subject, mailbody) // needs to be fixed // FIX WHAT?!?!!?!?
		if err != nil {
			showDialog(app, fmt.Sprintf("failed to send email: %v", err), form)
		} else {
			showDialog(app, "Gmail sent successfully - Please check!", form)
		}
	}).
		AddButton("QUIT", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Gmail-CLI by Dev_Vaayen").SetTitleAlign(tview.AlignCenter)

	if err := app.SetRoot(form, true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}

// accepts stuff from main func and uses it to send mail
func sendLoginNotification(email, password, toEmailId, subject, mailbody string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", email, password, smtpServer)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP address: %w", err)
	}

	to := []string{toEmailId}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + mailbody + "\r\n" +
		fmt.Sprintf("===================================================================\r\nThis mail was sent on %s from IP address %s.\r\n", timestamp, ip))

	// actual execution of sending
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, email, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// for security purpose also sending where the mail was sent from
// https://stackoverflow.com/a/37382208  <---- major influence
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// func showInbox(app *tview.Application, email, password string) {
// 	inboxView := tview.NewTextView().
// 		SetText("not yet implemented!!!!!!!!!!!!!!!!!!!!!!!!!").
// 		SetDynamicColors(true).
// 		SetWrap(true)

// 	// app.SetRoot(inboxView, true).Draw()
// 	app.SetRoot(inboxView, true).SetFocus(inboxView)
// }

// after sending the mail user should learn whether mail was sent or not
func showDialog(app *tview.Application, message string, returnTo tview.Primitive) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Return To Main Page"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if returnTo != nil {
				app.SetRoot(returnTo, true).SetFocus(returnTo)
			}
		})

	app.SetRoot(modal, true).SetFocus(modal)
}
