package main

import (
	"fmt"
	"log"
	"net"
	"net/smtp" // https://www.geeksforgeeks.org/sending-email-using-smtp-in-golang/ 
	"time"

	"github.com/rivo/tview"
)

var app *tview.Application
var globalEmail, globalPassword string

func main() {
	// main page FRONTEND
	app = tview.NewApplication()
	loginPage()
}

func loginPage() {
	globalEmail = ""
	globalPassword = ""
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter).
//		AddInputField("From:", "", 25, nil, nil).
		AddPasswordField("From:", "", 25, '*', nil).
		AddPasswordField("Password:", "", 25, '*', nil) // passwd is google application password not the actual
	// https://support.google.com/accounts/answer/185833?hl=en <<<<------ pls see

	form.AddButton("LOGIN!", func() {
		emailItem := form.GetFormItem(0).(*tview.InputField)
		passwordItem := form.GetFormItem(1).(*tview.InputField)

		// how to put gmail picture side by side?
		email := emailItem.GetText()
		password := passwordItem.GetText()

		err := loginFromThisDeviceAlert(email, password) // needs to be fixed // FIX WHAT?!?!!?!?
		if err != nil {
			showDialog(app, fmt.Sprintf("Login failed: %v", err), form)
		} else {
			// once user is logged into Gmail, they will be sent to next window to compose a mail
			globalEmail = email
			globalPassword = password
			landingPageMaybeLikeWebUI() // go to composeMail directly on successful login
		}
	}).AddButton("Exit", func() {
		app.Stop()
	})
	// form.AddImage // learnhow this works - can set up logo on login page
	form.SetBorder(true).SetTitle(" Gmail-TUI ").SetTitleAlign(tview.AlignCenter)
	form.SetButtonsAlign(tview.AlignCenter)
	app.SetRoot(form, true).EnableMouse(true).EnablePaste(true) // Set login form as the initial root
	if err := app.Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}

func landingPageMaybeLikeWebUI() {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)

	form.AddButton("Compose Mail", func() {
		composeMail()
	}).AddButton("Inbox (WIP)", func() {

	}).AddButton("Drafts (WIP)", func() {

	}).AddButton("Starred (WIP)", func() {

	}).AddButton("Exit", func() {
		app.Stop()
		app.Stop()
	})
	form.SetBorder(true).SetTitle(" Gmail-TUI ").SetTitleAlign(tview.AlignCenter)
	form.SetButtonsAlign(tview.AlignCenter).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true).EnableMouse(true).EnablePaste(true) // Set login form as the initial root
	if err := app.Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}

// In login page, login is being validated based on whether or not mail could be sent to self using entered email and password
func loginFromThisDeviceAlert(email, password string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", email, password, smtpServer)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP address: %w", err)
	}

	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + email + "\r\n" +
		fmt.Sprintf("\r\nGmail-TUI Login Alert on a device with %s IP address at %s.\r\n", ip, timestamp)) /// <--- this is success

	// actual execution of sending
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, email, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func composeMail() {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter).
		AddInputField("To:", "", 25, nil, nil).
		AddInputField("Subject:", "", 66, nil, nil).
		AddTextArea("Mail-Body:", "", 66, 0, 0, nil)

		// button configurations and what goes behind them kinda like their backend
	form.AddButton("SEND!", func() {
		toEMAIL := form.GetFormItem(0).(*tview.InputField)
		subjectMAIL := form.GetFormItem(1).(*tview.InputField)
		bodyMAIL := form.GetFormItem(2).(*tview.TextArea)

		// how to put gmail picture side by side?
		toEmailId := toEMAIL.GetText()
		subject := subjectMAIL.GetText()
		mailbody := bodyMAIL.GetText()

		err := sendSmtpMailLogic(globalEmail, globalPassword, toEmailId, subject, mailbody) // needs to be fixed // FIX WHAT?!?!!?!?
		if err != nil {
			showDialog(app, fmt.Sprintf("failed to send email: %v", err), form)
		} else {
			showDialog(app, "Gmail sent successfully - Please check!", form)
		}
	}).
		AddButton("Back", func() {
			landingPageMaybeLikeWebUI()
		})
	form.SetBorder(true).SetTitle(" Gmail-TUI ").SetTitleAlign(tview.AlignCenter)

	if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

}

// for accepting stuff from main func and uses it to send mail - after test will remove timestamping in mails
func sendSmtpMailLogic(globalEmail, globalPassword, toEmailId, subject, mailbody string) error {
	// https://www.geeksforgeeks.org/sending-email-using-smtp-in-golang/ <---- implementation of understanding from here
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", globalEmail, globalPassword, smtpServer)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP address: %w", err)
	}

	to := []string{toEmailId}
	msg := []byte("To: " + toEmailId + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + mailbody + "\r\n" +
		fmt.Sprintf("===================================================================\r\nThis mail was sent on %s from IP address %s.\r\n", timestamp, ip))

	// actual execution of sending
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, globalEmail, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// for security purpose also sending where the mail was sent from --->
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

// after sending the mail user should learn whether or not mail was sent - shown uses dialog box
func showDialog(app *tview.Application, message string, returnTo tview.Primitive) {
	// ----------> https://github.com/rivo/tview/wiki/Modal  <----------- very important 
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
