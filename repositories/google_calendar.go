package repositories

import (
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"time"
	"fmt"
	"log"
	"golang.org/x/net/context"
	"github.com/andboson/chebot/models"
	"strings"
)

func GetCalendarEventsList(text string) []string {
	var calendarId =  models.Conf.CalendarId
	var result []string
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return result
	}

	ctx := context.Background()
	config2, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Printf("Unable to create config: %v %+v", err)
		return result
	}

	c := config2.Client(ctx)
	srv, err := calendar.New(c)
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v", err)
		return result
	}

	t := time.Now().Format(time.RFC3339)
	to := time.Now().Add(time.Hour * time.Duration((18 - time.Now().Hour()))).Format(time.RFC3339)
	if strings.Contains(text, "tomorrow") {
		t = time.Now().Add(time.Hour * time.Duration((25 - time.Now().Hour()))).Add(6 * time.Hour).Format(time.RFC3339)
		to = time.Now().Add(time.Hour * time.Duration((25 - time.Now().Hour()))).Add( 18 * time.Hour).Format(time.RFC3339)
	}

	events, err := srv.Events.List(calendarId).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		TimeMax(to).
		OrderBy("startTime").Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err, calendarId)
		return result
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		log.Printf("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			start := item.Start.DateTime
			if start == "" {
				start = item.Start.Date
			}
			end := item.End.DateTime
			if end == "" {
				end = item.End.Date
			}
			parsedStart, _ := time.Parse(time.RFC3339, start)
			parsedEnd, _ := time.Parse(time.RFC3339, end)
			end = parsedEnd.Format("15:04")
			start = parsedStart.Format("15:04")
			
			if item.Summary !="" {
				event := fmt.Sprintf("%s (%s - %s )  %s", item.Summary, start, end, item.Creator.DisplayName)
				result = append(result, event)
			}
		}
	}
	
	return result
}
