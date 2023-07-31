package logic

var Web = lWeb{}

type WebService interface {
	ReloadService() error
	CloseService()
}

type lWeb struct {
	WebService
}

func (l *lWeb) RegisterService(s WebService) {
	l.WebService = s
}
