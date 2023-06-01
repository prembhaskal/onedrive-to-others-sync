// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

// <GraphHelperSnippet>
package graphhelper

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	nethttp "net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

type GraphHelper struct {
	deviceCodeCredential *azidentity.DeviceCodeCredential
	userClient           *msgraphsdk.GraphServiceClient
	graphUserScopes      []string
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

// </GraphHelperSnippet>

// <UserAuthConfigSnippet>
func (g *GraphHelper) InitializeGraphForUserAuth() error {
	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	scopes := os.Getenv("GRAPH_USER_SCOPES")
	g.graphUserScopes = strings.Split(scopes, ",")

	// Create the device code credential
	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		TenantID: tenantId,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		return err
	}

	g.deviceCodeCredential = credential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, g.graphUserScopes)
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.userClient = client

	return nil
}

// </UserAuthConfigSnippet>

// <GetUserTokenSnippet>
func (g *GraphHelper) GetUserToken() (*string, error) {
	token, err := g.deviceCodeCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.graphUserScopes,
	})
	if err != nil {
		return nil, err
	}

	return &token.Token, nil
}

// </GetUserTokenSnippet>

// <GetUserSnippet>
func (g *GraphHelper) GetUser() (models.Userable, error) {
	query := users.UserItemRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "mail", "userPrincipalName"},
	}

	return g.userClient.Me().Get(context.Background(),
		&users.UserItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &query,
		})
}

// </GetUserSnippet>

// <GetInboxSnippet>
func (g *GraphHelper) GetInbox() (models.MessageCollectionResponseable, error) {
	var topValue int32 = 25
	query := users.ItemMailFoldersItemMessagesRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"from", "isRead", "receivedDateTime", "subject"},
		// Get at most 25 results
		Top: &topValue,
		// Sort by received time, newest first
		Orderby: []string{"receivedDateTime DESC"},
	}

	return g.userClient.Me().MailFolders().
		ByMailFolderId("inbox").
		Messages().
		Get(context.Background(),
			&users.ItemMailFoldersItemMessagesRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

// </GetInboxSnippet>

// <SendMailSnippet>
func (g *GraphHelper) SendMail(subject *string, body *string, recipient *string) error {
	// Create a new message
	message := models.NewMessage()
	message.SetSubject(subject)

	messageBody := models.NewItemBody()
	messageBody.SetContent(body)
	contentType := models.TEXT_BODYTYPE
	messageBody.SetContentType(&contentType)
	message.SetBody(messageBody)

	toRecipient := models.NewRecipient()
	emailAddress := models.NewEmailAddress()
	emailAddress.SetAddress(recipient)
	toRecipient.SetEmailAddress(emailAddress)
	message.SetToRecipients([]models.Recipientable{
		toRecipient,
	})

	sendMailBody := users.NewItemSendMailPostRequestBody()
	sendMailBody.SetMessage(message)

	// Send the message
	return g.userClient.Me().SendMail().Post(context.Background(), sendMailBody, nil)
}

// </SendMailSnippet>

// <MakeGraphCallSnippet>
func (g *GraphHelper) MakeGraphCall() error {
	// INSERT YOUR CODE HERE
	return nil
}

// </MakeGraphCallSnippet>

func (g *GraphHelper) ListFiles() error {
	ctx := context.Background()

	// result, err := g.userClient.Me().Drives().Get(context.Background(), nil)
	drvs, err := g.userClient.Me().Drive().Get(context.Background(), nil)
	if err != nil {
		return err
	}

	fmt.Printf("drive id: %s, description: %s\n", *drvs.GetId(), *drvs.GetDescription())

	rootItem, err := g.userClient.Drives().ByDriveId(*drvs.GetId()).Root().Get(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Printf("root id: %s, name: %s\n", *rootItem.GetId(), *rootItem.GetName())

	result, err := g.userClient.Drives().ByDriveId(*drvs.GetId()).Items().ByDriveItemId(*rootItem.GetId()).Children().Get(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Printf("result: %#v\n", result)

	fmt.Println("list items")
	for _, item := range result.GetValue() {
		fmt.Printf("item -> id: %s, name: %s\n", *item.GetId(), *item.GetName())
	}

	// fmt.Printf("ListFiles results: %#v\n", result)
	// fmt.Printf("ListFiles result.GetValue: %#v\n", result.GetValue())
	// for _, ds := range result.GetValue() {
	// 	fmt.Printf("drive id: %s, description: %s\n", *ds.GetId(), *ds.GetDescription())
	// }

	return err
}

func (g *GraphHelper) UploadFile() error {
	ctx := context.Background()
	filename := "test2.out"
	// item -> id: 01VUOUUTCMY6YN5SQELZEIFBXBN6XZJXBM, name: home
	driveID := "b!4sPTmBDuQkCY6AyFVL1V6YCpvxwLFl9InPT023kInlMzOIMLIcjLSbmtYklBw0un"
	homeItemId := "01VUOUUTCMY6YN5SQELZEIFBXBN6XZJXBM"
	// reqBody := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	// res, err := g.userClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(homeItemId).CreateUploadSession().Post(ctx, reqBody, nil)
	uploadSession, err := CreateLocalUploadSession(g.userClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(homeItemId), filename).Post(ctx, nil, nil)
	if err != nil {
		return err
	}

	fmt.Printf("uploadSession: %#v\n", uploadSession)
	fmt.Printf("uploadSession: url: %s, expire time: %s\n", *uploadSession.GetUploadUrl(), uploadSession.GetExpirationDateTime())

	return g.uploadFileChunks(*uploadSession.GetUploadUrl())
}

func (g *GraphHelper) uploadFileChunks(upSessionUrl string) error {
	ctx := context.Background()
	filetoUploadName := "/tmp/data.out"
	fileToUpload, err := os.Open(filetoUploadName)
	if err != nil {
		return err
	}
	defer fileToUpload.Close()

	filestat, err := fileToUpload.Stat()
	if err != nil {
		return err
	}
	filesize := filestat.Size()

	filereader := bufio.NewReader(fileToUpload)

	chunksize := 5 * 1024 * 1024 // 5 MB

	var totalread int64
	// create the http client

	httpclient := g.createHttpClient(30, 30)

	isEOF := false
	var chunk int
	for totalread < filesize || !isEOF {
		// var buf bytes.Reader
		// read until chunksize
		var off int
		buf := make([]byte, chunksize)
		for off < chunksize {
			n, err := filereader.Read(buf[off:])
			off += n
			if errors.Is(err, io.EOF) {
				isEOF = true
				break
			} else if err != nil {
				return err
			}
		}

		putRequest, err := nethttp.NewRequestWithContext(ctx, "PUT", upSessionUrl, bytes.NewReader(buf[:off]))
		if err != nil {
			return fmt.Errorf("new request context error: %w", err)
		}
		beginIndex := totalread
		endIndex := totalread + int64(off - 1)

		putRequest.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", beginIndex, endIndex, filesize))

		fmt.Printf("chunk: %d, request: %#v\n", chunk, putRequest)

		resp, err := httpclient.Do(putRequest)
		if err != nil {
			return fmt.Errorf("upload resp err: %w", err)
		}

		fmt.Printf("chunk uploaded response, chunk: %d, begin: %d, end: %d, resp: %#v\n", chunk, beginIndex, endIndex, resp)

		totalread += int64(off)
		chunk++
	}
	return nil
}

func (g *GraphHelper) createHttpClient(requestTimeout, connectTimeout int) *nethttp.Client {
	return &nethttp.Client{
		Timeout: time.Duration(requestTimeout) * time.Second,
		Transport: &nethttp.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Duration(connectTimeout) * time.Second,
			}).Dial,
		},
	}
}

func CreateLocalUploadSession(m *drives.ItemItemsDriveItemItemRequestBuilder, filename string) *drives.ItemItemsItemCreateUploadSessionRequestBuilder {
	urlTemplateWithPath := fmt.Sprintf("{+baseurl}/drives/{drive%%2Did}/items/{driveItem%%2Did}:/%s:/createUploadSession", filename)
	up := &drives.ItemItemsItemCreateUploadSessionRequestBuilder{
		BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(m.BaseRequestBuilder.RequestAdapter, urlTemplateWithPath, m.BaseRequestBuilder.PathParameters),
	}
	return up
}
