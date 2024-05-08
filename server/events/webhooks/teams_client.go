package webhooks

import (
	"fmt"
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"strings"
)

const (
	teamsSuccessColour string = "2DC72D"
	teamsFailureColour string = "8C1A1A"
)

type TeamsClient struct {
	client *goteamsnotify.TeamsClient
	Token  string
}

func NewTeamsClient(token string) TeamsClient {
	return TeamsClient{
		client: goteamsnotify.NewTeamsClient(),
		Token:  token,
	}
}

func (t *TeamsClient) AuthTest() error {
	err := t.client.ValidateWebhook(t.Token)

	return err
}

func (t *TeamsClient) TokenIsSet() bool {
	return t.Token != ""
}

func (t *TeamsClient) PostMessage(_ string, applyResult ApplyResult) error {
	var (
		err     error
		msgCard *messagecard.MessageCard
	)

	msgCard, err = CreateMessageCard(applyResult)
	if err != nil {
		return err
	}

	err = t.client.Send(t.Token, msgCard)
	if err != nil {
		return err
	}

	fmt.Println("Send message successful trough teams client.")
	return nil
}

func CreateMessageCard(applyResult ApplyResult) (*messagecard.MessageCard, error) {
	var colour, successWord string

	// Apply status judgment
	if applyResult.Success {
		colour = teamsSuccessColour
		successWord = "Succeeded"
	} else {
		colour = teamsFailureColour
		successWord = "Failed"
	}

	// Message text generator
	msgText := fmt.Sprintf("Apply %s for [%s](%s)", strings.ToLower(successWord), applyResult.Repo.FullName, applyResult.Pull.URL)

	// Message card section
	msgCardSection := messagecard.NewSection()
	msgCardSection.ActivityTitle = msgText
	msgCardSection.Markdown = true
	msgCardSection.Facts = []messagecard.SectionFact{
		{
			Name:  "Workspace",
			Value: applyResult.Workspace,
		},
		{
			Name:  "User",
			Value: applyResult.User.Username,
		},
		{
			Name:  "Directory",
			Value: applyResult.Directory,
		},
	}

	// Message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = fmt.Sprintf("Atlantis Apply Notice (%s)", successWord)
	msgCard.Text = ""
	msgCard.Summary = "Atlantis"
	msgCard.ThemeColor = colour
	msgCard.Sections = []*messagecard.Section{msgCardSection}

	return msgCard, nil
}
