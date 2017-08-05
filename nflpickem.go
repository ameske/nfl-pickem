package nflpickem

// Service is the interface implemented by types that can provide
// all the various services needed by the NFL Pickem Pool
type Service interface {
	Weeker
	GamesRetriever
	PasswordUpdater
	Picker
	PickRetriever
	ResultFetcher
	WeekTotalFetcher
	CredentialChecker
	DataSummarizer
	UserAdder
	GameAdder
	DateAdder
	PickCreater
}

type Notifier interface {
	Notify(to string, week int, picks []Pick) error
}

type DataSummarizer interface {
	Years() ([]int, error)
}
