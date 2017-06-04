package schedule

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ConditionFunc func(html.Token) bool

type Matchup struct {
	Date time.Time `json:"date"`
	Away string    `json:"away"`
	Home string    `json:"home"`
}

func (m Matchup) String() string {
	return fmt.Sprintf("%s\t%s at %s", m.Date, m.Away, m.Home)
}

type Parser struct {
	*html.Tokenizer
	Month int
	Date  int
	Year  int
}

func NewParser(year int, r io.Reader) *Parser {
	return &Parser{
		Tokenizer: html.NewTokenizer(r),
		Month:     0,
		Date:      0,
		Year:      year,
	}
}

func (p *Parser) Parse() ([]Matchup, error) {
	matchups := make([]Matchup, 0)

	for p.nextMatchup() == nil {
		date, err := p.extractNodeText(Date)
		if err != nil {
			return nil, err
		}

		meridian, err := p.extractNodeText(Meridian)
		if err != nil {
			return nil, err
		}

		away, err := p.extractNodeText(AwayTeam)
		if err != nil {
			return nil, err
		}

		home, err := p.extractNodeText(HomeTeam)
		if err != nil {
			return nil, err
		}

		hourMin := strings.Split(date, ":")

		hour, err := strconv.ParseInt(hourMin[0], 10, 64)
		if err != nil {
			return nil, err
		}
		if meridian == "pm" {
			hour += 12
		}

		min, err := strconv.ParseInt(hourMin[1], 10, 64)
		if err != nil {
			return nil, err
		}

		t := time.Date(p.Year, time.Month(p.Month), p.Date, int(hour), int(min), 0, 0, time.Now().Location())

		matchups = append(matchups, Matchup{Date: t, Away: away, Home: home})
	}

	return matchups, nil
}

func (p *Parser) nextMatchup() error {
	return p.advanceUntil(MatchupStart)
}

func (p *Parser) advanceUntil(f ConditionFunc) error {
	for {
		tt := p.Next()
		if tt == html.ErrorToken {
			return p.Err()

		}

		t := p.Token()

		if ScheduleDate(t) {
			err := p.updateDate()
			if err != nil {
				return err
			}

			continue
		}

		if f(t) {
			return nil
		}
	}
}

func (p *Parser) extractNodeText(f ConditionFunc) (string, error) {
	for {
		tt := p.Next()
		if tt == html.ErrorToken {
			return "", p.Err()
		}

		if tt == html.StartTagToken {
			t := p.Token()
			if f(t) {
				p.Next()
				t := p.Token()
				return t.Data, nil
			}
		}
	}
}

// This one is pretty hard coded, advance 3 tokens and we should have it
func (p *Parser) updateDate() error {
	p.Next()
	p.Next()
	p.Next()

	p.Next()
	t := p.Token()
	month, day, err := parseDate(t.Data)
	if err != nil {
		return err
	}

	p.Month = month
	p.Date = day

	return nil
}

func ScheduleDate(t html.Token) bool {
	return classEquals(t, "schedules-list-date")
}

func MatchupStart(t html.Token) bool {
	return classEquals(t, "list-matchup-row-time")
}

func HomeTeam(t html.Token) bool {
	return classEquals(t, "team-name home ")
}

func AwayTeam(t html.Token) bool {
	return classEquals(t, "team-name away ")
}

func Date(t html.Token) bool {
	return classEquals(t, "time")
}

func Meridian(t html.Token) bool {
	return classEquals(t, "am") || classEquals(t, "pm")
}

func classEquals(t html.Token, class string) bool {
	for _, a := range t.Attr {
		if a.Key == "class" && a.Val == class {
			return true
		}
	}

	return false
}

func parseDate(date string) (int, int, error) {
	parts := strings.Split(date, " ")
	if len(parts) != 3 {
		log.Println("Invalid string: ", date)
		return -1, -1, errors.New("unable to parse")
	}

	var month int
	switch parts[1] {
	case "September":
		month = 9
	case "October":
		month = 10
	case "November":
		month = 11
	case "December":
		month = 12
	case "January":
		month = 1
	}

	day, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	return month, int(day), nil
}
