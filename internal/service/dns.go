package service

import (
	"ddns-server/internal/config"
	"ddns-server/internal/dao"
	"ddns-server/internal/logic"
	"net"
	"strings"
	"sync"

	"github.com/abxuz/b-tools/bmap"
	"github.com/miekg/dns"
)

var Dns = sDns{
	lock: &sync.RWMutex{},
}

type sDns struct {
	lock    *sync.RWMutex
	server  *dns.Server
	config  *config.Dns
	records map[string]*config.Record
}

func init() {
	logic.Dns.RegisterService(&Dns)
}

func (s *sDns) ReloadRecord() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	rs, err := dao.Record.List()
	if err != nil {
		return err
	}
	s.records = bmap.NewMapFromSlice(rs, func(r *config.Record) string { return r.Domain })
	return nil
}

func (s *sDns) ReloadService() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	config, err := dao.Dns.Get()
	if err != nil {
		return err
	}

	if s.config == nil || s.config.Listen == "" {
		if config == nil || config.Listen == "" {
			return nil
		}

		s.server = &dns.Server{
			Addr:    config.Listen,
			Net:     "udp",
			Handler: s,
		}
		go s.server.ListenAndServe()
		s.config = config
		return nil
	}

	if config == nil || config.Listen == "" {
		s.server.Shutdown()
		s.server = nil
		s.config = config
		return nil
	}

	if config.Listen == s.config.Listen {
		return nil
	}

	s.server.Shutdown()
	s.server = &dns.Server{
		Addr:    config.Listen,
		Net:     "udp",
		Handler: s,
	}
	go s.server.ListenAndServe()
	s.config = config
	return nil
}

func (s *sDns) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if r.Opcode != dns.OpcodeQuery {
		return
	}

	if len(r.Question) != 1 {
		return
	}

	question := r.Question[0]
	name := strings.TrimSuffix(question.Name, ".")
	record, ok := s.records[strings.ToLower(name)] // Google DNS会随机大小写
	if !ok {
		w.WriteMsg(new(dns.Msg).SetRcode(r, dns.RcodeNameError))
		return
	}

	replyMsg := new(dns.Msg).SetReply(r)
	switch question.Qtype {
	case dns.TypeA:
		if record.Ipv4 == "" {
			break
		}
		addr := net.ParseIP(record.Ipv4)
		if addr == nil {
			break
		}
		if addr = addr.To4(); addr == nil {
			break
		}
		replyMsg.Answer = append(replyMsg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: question.Qtype,
				Class:  question.Qclass,
				Ttl:    1,
			},
			A: addr,
		})
	case dns.TypeAAAA:
		if record.Ipv6 == "" {
			break
		}
		addr := net.ParseIP(record.Ipv6)
		if addr == nil {
			break
		}
		if addr.To4() != nil {
			break
		}
		replyMsg.Answer = append(replyMsg.Answer, &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: question.Qtype,
				Class:  question.Qclass,
				Ttl:    1,
			},
			AAAA: addr,
		})
	}

	w.WriteMsg(replyMsg)
}

func (s *sDns) CloseService() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.server != nil {
		s.server.Shutdown()
		s.server = nil
		s.config = nil
	}
}
