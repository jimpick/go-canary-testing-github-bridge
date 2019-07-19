package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v26/github"
)

var webhookSecretKey = []byte("ipfs_secret")

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println("Push",
				event.GetRepo().GetFullName(),
				event.GetRef(),
				event.GetAfter())
		case *github.PullRequestEvent:
			fmt.Println("PullRequest",
				event.GetRepo().GetFullName(),
				event.GetNumber(),
				event.GetAction(), // opened, synchronize
				event.GetPullRequest().GetHead().GetRepo().GetFullName(),
				event.GetPullRequest().GetHead().GetRef(),
				event.GetPullRequest().GetHead().GetSHA())
		case *github.IssueCommentEvent:
			fmt.Println("IssueComment",
				event.GetRepo().GetFullName(),
				event.GetIssue().GetNumber(),
				event.GetAction()) // created
		default:
			fmt.Printf("GitHub event type %T\n", event)
		}
		fmt.Fprintf(w, "OK")
		fmt.Println()
	})
	log.Fatal(http.ListenAndServe(":14001", nil))
}
