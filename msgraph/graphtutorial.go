// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

// <ProgramSnippet>
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/prembhaskal/onedrive-to-others-sync/msgraph/graphhelper"

	"github.com/joho/godotenv"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func main() {
	fmt.Println("Go Graph Tutorial")
	fmt.Println()

	// Load .env files
	// .env.local takes precedence (if present)
	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	graphHelper := graphhelper.NewGraphHelper()

	initializeGraph(graphHelper)

	greetUser(graphHelper)

	var choice int64 = -1

	for {
		fmt.Println("Please choose one of the following options:")
		fmt.Println("0. Exit")
		fmt.Println("1. Display access token")
		fmt.Println("2. List my inbox")
		fmt.Println("3. Send mail")
		fmt.Println("4. Make a Graph call")
		fmt.Println("5. Access Drive")

		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			choice = -1
		}

		switch choice {
		case 0:
			// Exit the program
			fmt.Println("Goodbye...")
		case 1:
			// Display access token
			displayAccessToken(graphHelper)
		case 2:
			// List emails from user's inbox
			listInbox(graphHelper)
		case 3:
			// Send an email message
			sendMail(graphHelper)
		case 4:
			// Run any Graph code
			makeGraphCall(graphHelper)
		case 5:
			accessDrive(graphHelper)
		default:
			fmt.Println("Invalid choice! Please try again.")
		}

		if choice == 0 {
			break
		}
	}
}

// </ProgramSnippet>

// <InitializeGraphSnippet>
func initializeGraph(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.InitializeGraphForUserAuth()
	if err != nil {
		log.Panicf("Error initializing Graph for user auth: %v\n", err)
	}
}

// </InitializeGraphSnippet>

// <GreetUserSnippet>
func greetUser(graphHelper *graphhelper.GraphHelper) {
	user, err := graphHelper.GetUser()
	if err != nil {
		log.Panicf("Error getting user: %v\n", err)
	}

	fmt.Printf("Hello, %s!\n", *user.GetDisplayName())

	// For Work/school accounts, email is in Mail property
	// Personal accounts, email is in UserPrincipalName
	email := user.GetMail()
	if email == nil {
		email = user.GetUserPrincipalName()
	}

	fmt.Printf("Email: %s\n", *email)
	fmt.Println()
}

// </GreetUserSnippet>

// <DisplayAccessTokenSnippet>
func displayAccessToken(graphHelper *graphhelper.GraphHelper) {
	token, err := graphHelper.GetUserToken()
	if err != nil {
		log.Panicf("Error getting user token: %v\n", err)
	}

	fmt.Printf("User token: %s", *token)
	fmt.Println()
}

// </DisplayAccessTokenSnippet>

// <ListInboxSnippet>
func listInbox(graphHelper *graphhelper.GraphHelper) {
	messages, err := graphHelper.GetInbox()
	if err != nil {
		log.Panicf("Error getting user's inbox: %v", err)
	}

	// Load local time zone
	// Dates returned by Graph are in UTC, use this
	// to convert to local
	location, err := time.LoadLocation("Local")
	if err != nil {
		log.Panicf("Error getting local timezone: %v", err)
	}

	// Output each message's details
	for _, message := range messages.GetValue() {
		fmt.Printf("Message: %s\n", *message.GetSubject())
		fmt.Printf("  From: %s\n", *message.GetFrom().GetEmailAddress().GetName())

		status := "Unknown"
		if *message.GetIsRead() {
			status = "Read"
		} else {
			status = "Unread"
		}
		fmt.Printf("  Status: %s\n", status)
		fmt.Printf("  Received: %s\n", (*message.GetReceivedDateTime()).In(location))
	}

	// If GetOdataNextLink does not return nil,
	// there are more messages available on the server
	nextLink := messages.GetOdataNextLink()

	fmt.Println()
	fmt.Printf("More messages available? %t\n", nextLink != nil)
	fmt.Println()
}

// </ListInboxSnippet>

// <SendMailSnippet>
func sendMail(graphHelper *graphhelper.GraphHelper) {
	// Send mail to the signed-in user
	// Get the user for their email address
	user, err := graphHelper.GetUser()
	if err != nil {
		log.Panicf("Error getting user: %v", err)
	}

	// For Work/school accounts, email is in Mail property
	// Personal accounts, email is in UserPrincipalName
	email := user.GetMail()
	if email == nil {
		email = user.GetUserPrincipalName()
	}

	subject := "Testing Microsoft Graph"
	body := "Hello world!"
	err = graphHelper.SendMail(&subject, &body, email)
	if err != nil {
		log.Panicf("Error sending mail: %v", err)
	}

	fmt.Println("Mail sent.")
	fmt.Println()
}

// </SendMailSnippet>

// <MakeGraphCallSnippet>
func makeGraphCall(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.MakeGraphCall()
	if err != nil {
		log.Panicf("Error making Graph call: %v", err)
	}
}

func accessDrive(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.ListFiles()
	if err != nil {
		log.Fatalf("error in list files: %v\n", err)
	}

	err = graphHelper.UploadFile()
	if err != nil {
		printOdataError(err)
		log.Fatalf("error in upload files: %#v\n", err)
	}
}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error: %v", typed.Error())
		if terr := typed.GetError(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
	fmt.Println("")
}

// </MakeGraphCallSnippet>
