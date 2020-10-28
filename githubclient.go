package main

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

type release struct {
	repName   string
	releaseAt time.Time
	version   string
	url       string
}

func getReleasesAfterDate(t time.Time, token string) []release {
	var result []release
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	ops := github.NotificationListOptions{All: true, Since: t}
	notifications, _, err := client.Activity.ListNotifications(ctx, &ops)
	if err != nil {
		log.Printf("Error on get notifications: %v\n", err)
	}
	for _, notification := range notifications {
		if notification.Subject.GetType() == "Release" {
			log.Printf("Repository: %v, subject: %v. Title: %v", notification.GetRepository().GetName(), notification.Subject.GetType(), notification.Subject.GetTitle())
			release := release{releaseAt: notification.GetUpdatedAt(), repName: notification.GetRepository().GetName(), url: notification.GetRepository().GetHTMLURL(), version: notification.Subject.GetTitle()}
			result = append(result, release)
		}
	}

	return result
}
