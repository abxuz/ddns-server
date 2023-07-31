package logic

var Dns = lDns{}

type DnsService interface {
	ReloadRecord() error
	ReloadService() error
	CloseService()
}

type lDns struct {
	DnsService
}

func (l *lDns) RegisterService(s DnsService) {
	l.DnsService = s
}
