<div align=center>

# Gmail-TUI by Dev_Vaayen

![img](https://i.imgur.com/hRZ4VRz.gif)

</div>

A simple TUI application that aims to replicate the Gmail Web-UI in a TUI-Environment - Is this even possible? I don't know yet but let's find out! Special thanks to Rivo for their [TUI Library](https://github.com/rivo/tview/tree/master).

DevLogs for this project can be found below:     
- [[DevLog #01] Gmail-TUI: Replicating The Gmail-Web Experience In Terminal](https://dev.to/dev_vaayen/devlog-01-gmail-tui-replicating-the-gmail-web-experience-in-terminal-1lk1)

## Future Plans
As of now, users can only compose mails and send them to Email-IDs from their Gmail-ID using the [Go SMTP-Library](https://www.geeksforgeeks.org/sending-email-using-smtp-in-golang/). I plan on implementing the following to this TUI-Application:
- A login page for entering email-ID and password
- Composing and sending mails - Implemented today!
- Listing received emails with email-IDs in the Inbox
- Opening the content of the received mail after clicking it
- Viewing sent email in Sent-Box
- A small panel on the left side to choose from the Compose, Inbox, Drafts, Sent buttons.

## Instructions
1. Create application password for your Gmail-ID | [Link](https://support.google.com/accounts/answer/185833?hl=en)
2. Clone this repository and inside it, run `./Gmail-TUI`
3. Enter the required details in the placeholders - Please only enter the application password created in step 1
4. Use `tab` and `shift+tab` to navigate between the buttons and hit `enter` on the **SEND!** button
5. Check your inbox after the success dialogue box has been displayed 
