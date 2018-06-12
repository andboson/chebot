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
)

func GetCalendarEventsList() []string {
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
	list, err := srv.CalendarList.List().Do()
	log.Printf("list: %+v - %s", *list.Items[0], list.Items)

	events, err := srv.Events.List(calendarId).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		TimeMax(to).
		OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		log.Println("No upcoming events found.")
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
				event := fmt.Sprintf("%s (%s - %s )  %s", item.Summary, start, end, item.Creator)
				result = append(result, event)
			}
		}
	}
	
	return result
}
