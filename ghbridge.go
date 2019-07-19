package ghbridge

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/go-github/v26/github"
)

// PushEvent is a GitHub push event
type PushEvent struct {
	repo string
	ref  string
	sha  string
}

// PullRequestEvent is a GitHub pull request event
type PullRequestEvent struct {
	repo     string
	number   int
	action   string
	fromRepo string
	fromRef  string
	sha      string
}

// IssueCommentEvent is a GitHub issue comment event
type IssueCommentEvent struct {
	repo   string
	number int
	action string
}

// UnmappedEvent is a GitHub event that isn't mapped
type UnmappedEvent struct {
	eventType string
}

// GetHandler builds a GitHub webhook handler for net/http
func GetHandler(c chan interface{}, webhookSecretKey []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := github.ValidatePayload(r, webhookSecretKey)
		if err != nil {
			panic(err)
		}
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			panic(err)
		}
		// fmt.Println("GitHub event", event)
		switch event := event.(type) {
		case *github.PushEvent:
			c <- &PushEvent{
				repo: event.GetRepo().GetFullName(),
				ref:  event.GetRef(),
				sha:  event.GetAfter(),
			}
		case *github.PullRequestEvent:
			c <- &PullRequestEvent{
				repo:     event.GetRepo().GetFullName(),
				number:   event.GetNumber(),
				action:   event.GetAction(),
				fromRepo: event.GetPullRequest().GetHead().GetRepo().GetFullName(),
				fromRef:  event.GetPullRequest().GetHead().GetRef(),
				sha:      event.GetPullRequest().GetHead().GetSHA(),
			}
		case *github.IssueCommentEvent:
			c <- &IssueCommentEvent{
				repo:   event.GetRepo().GetFullName(),
				number: event.GetIssue().GetNumber(),
				action: event.GetAction(),
			}
		default:
			c <- &UnmappedEvent{
				eventType: reflect.TypeOf(event).String(),
			}
		}
		fmt.Fprintf(w, "OK")
	}
}
