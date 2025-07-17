package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

type UserTimezone struct {
	Name     string
	Username string
	Timezone string
	Offset   int
}

type SlackClient struct {
	client *slack.Client
}

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		client: slack.New(token),
	}
}

func (sc *SlackClient) GetUserTimezones() ([]UserTimezone, error) {
	users, err := sc.client.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	var userTimezones []UserTimezone
	for _, user := range users {
		if user.IsBot || user.Deleted {
			continue
		}

		tz := user.TZ
		if tz == "" {
			tz = "UTC"
		}

		userTimezones = append(userTimezones, UserTimezone{
			Name:     user.RealName,
			Username: user.Name,
			Timezone: tz,
			Offset:   user.TZOffset,
		})
	}

	return userTimezones, nil
}

func (ut UserTimezone) GetCurrentTime() time.Time {
	loc, err := time.LoadLocation(ut.Timezone)
	if err != nil {
		log.Printf("Error loading timezone %s: %v", ut.Timezone, err)
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}

func (ut UserTimezone) GetTimeIn(t time.Time) time.Time {
	loc, err := time.LoadLocation(ut.Timezone)
	if err != nil {
		log.Printf("Error loading timezone %s: %v", ut.Timezone, err)
		return t.UTC()
	}
	return t.In(loc)
}

func (ut UserTimezone) GetTimezoneOffset() string {
	now := time.Now()

	// Get the user's timezone location
	userLoc, err := time.LoadLocation(ut.Timezone)
	if err != nil {
		log.Printf("Error loading timezone %s: %v", ut.Timezone, err)
		return "?"
	}

	// Get the local timezone location
	localLoc := now.Location()

	// Get the same moment in both timezones
	nowInUserTZ := now.In(userLoc)
	nowInLocalTZ := now.In(localLoc)

	// Calculate the offset difference in seconds, then convert to hours
	_, userOffset := nowInUserTZ.Zone()
	_, localOffset := nowInLocalTZ.Zone()

	diffSeconds := userOffset - localOffset
	diffHours := diffSeconds / 3600

	if diffHours == 0 {
		return "Same"
	} else if diffHours > 0 {
		return fmt.Sprintf("+%dh", diffHours)
	} else {
		return fmt.Sprintf("%dh", diffHours)
	}
}

func getSlackToken() string {
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		token = os.Getenv("SLACK_USER_TOKEN")
	}
	return token
}
