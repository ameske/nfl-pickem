package results

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"golang.org/x/net/html"
)

type ConditionFunc func(html.Token) bool

type Result struct {
	Away      string `json:"away"`
	AwayScore int    `json:"awayScore"`
	Home      string `json:"home"`
	HomeScore int    `json:"homeScore"`
}

func (r Result) String() string {
	return fmt.Sprintf("%s (%d) at %s (%d)", r.Away, r.AwayScore, r.Home, r.HomeScore)
}

type Parser struct {
	*html.Tokenizer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		Tokenizer: html.NewTokenizer(r),
	}
}

// Parse finds all the results for games that have finished in the given week the parser is set to.
func (p *Parser) Parse() ([]Result, error) {
	results := make([]Result, 0)

	for p.nextMatchup() == nil {
		away, err := p.extractNodeText(AwayTeam)
		if err != nil {
			log.Println("Couldn't find away team")
			return nil, err
		}

		awayScoreStr, err := p.extractNodeText(AwayTeamScore)
		if err != nil {
			continue
		}
		awayScore, err := strconv.ParseInt(awayScoreStr, 10, 64)
		if err != nil {
			return nil, err
		}

		homeScoreStr, err := p.extractNodeText(HomeTeamScore)
		if err != nil {
			continue
		}
		homeScore, err := strconv.ParseInt(homeScoreStr, 10, 64)
		if err != nil {
			return nil, err
		}

		home, err := p.extractNodeText(HomeTeam)
		if err != nil {
			log.Println("Couldn't find home team")
			return nil, err
		}

		results = append(results, Result{Away: away, AwayScore: int(awayScore), Home: home, HomeScore: int(homeScore)})
	}

	return results, nil
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

func MatchupStart(t html.Token) bool {
	return classEquals(t, "list-matchup-row-team")
}

func HomeTeam(t html.Token) bool {
	return classEquals(t, "team-name home ") || classEquals(t, "team-name home lost")
}

func HomeTeamScore(t html.Token) bool {
	return classEquals(t, "team-score home lost") || classEquals(t, "team-score home ")
}

func AwayTeam(t html.Token) bool {
	return classEquals(t, "team-name away ") || classEquals(t, "team-name away lost")
}

func AwayTeamScore(t html.Token) bool {
	return classEquals(t, "team-score away lost") || classEquals(t, "team-score away ")
}

func classEquals(t html.Token, class string) bool {
	for _, a := range t.Attr {
		if a.Key == "class" && a.Val == class {
			return true
		}
	}

	return false
}
