package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem"
	"github.com/ameske/nfl-pickem/parser/results"
)

// scheduleUpdates sets up goroutines that will import the results of games and update the
// picks after every wave of games completes.
func scheduleUpdates(db nflpickem.Updater) {
	// Friday at 8:00
	go func() {
		nextFriday := adjustIfPast(nextDay(time.Friday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextFriday)
		time.Sleep(nextFriday.Sub(time.Now()))
		for {
			go update(db, false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Sunday at 18:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 18))
		logNextScheduleUpdate(nextSunday)
		time.Sleep(nextSunday.Sub(time.Now()))
		for {
			go update(db, false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Sunday at 21:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 21))
		logNextScheduleUpdate(nextSunday)
		time.Sleep(nextSunday.Sub(time.Now()))
		for {
			go update(db, false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Monday at 8:00
	go func() {
		nextMonday := adjustIfPast(nextDay(time.Monday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextMonday)
		time.Sleep(nextMonday.Sub(time.Now()))
		for {
			go update(db, false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Tuesday at 8:00. Here we need to update the current week - 1
	go func() {
		nextTuesday := adjustIfPast(nextDay(time.Tuesday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextTuesday)
		time.Sleep(nextTuesday.Sub(time.Now()))
		for {
			go update(db, true)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()
}

func update(db nflpickem.Updater, updatePreviousWeek bool) {
	nflWeek, err := db.CurrentWeek(time.Now())
	if err != nil {
		log.Println(err)
		return
	}

	if updatePreviousWeek {
		nflWeek.Week -= 1
	}

	results, err := getGameResults(nflWeek.Year, nflWeek.Week)
	if err != nil {
		log.Println(err)
		return
	}

	for _, result := range results {
		log.Printf("Updating Game: %v", result)
		err := db.UpdateGame(nflWeek.Week, nflWeek.Year, result.Home, result.HomeScore, result.AwayScore)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func getGameResults(year, week int) ([]results.Result, error) {
	url := fmt.Sprintf("http://www.nfl.com/schedules/%d/REG%d", year, week)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	p := results.NewParser(resp.Body)

	return p.Parse()
}

func adjustIfPast(next time.Time) time.Time {
	now := time.Now()
	ty, tm, td := now.Date()
	ny, nm, nd := next.Date()

	// if the next day is today, but the hour we want has past then we must advance next a week
	if (ty == ny && tm == nm && td == nd) && now.Hour() > next.Hour() {
		next = next.AddDate(0, 0, 7)
	}

	return next
}

func nextDay(day time.Weekday) time.Time {
	now := time.Now()

	// We only want to go forwards, so use modular arith to force going ahead
	diff := int(day-now.Weekday()+7) % 7

	next := now.AddDate(0, 0, diff)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	return next
}

func logNextScheduleUpdate(t time.Time) {
	slog.Info("Scheduling update for " + t.Format(time.RFC1123))
}
