package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	nslcdConf     = "/etc/nslcd.conf"
	defaultFilter = "(&(objectClass=posixAccount)(uid=%s))"
)

func getNslcdConfPath() string {
	conf := os.Getenv("NSLCD_CONF")
	if conf == "" {
		conf = nslcdConf
	}
	return conf
}

func (l *ldapEnv) loadNslcdConf() {
	conf := getNslcdConfPath()
	b, err := ioutil.ReadFile(conf)
	if err != nil {
		l.host = "localhost"
		l.port = 389
		l.base = "dc=example,dc=org"
		l.filter = defaultFilter
		l.tls = false
	}
	for _, s := range strings.Split(string(b), "\n") {
		v := strings.Split(s, " ")
		switch {
		case v[0] == "uri":
			u, err := url.Parse(v[1])
			if err != nil {
				log.Fatal(err)
			}
			if u.Scheme == "ldaps" {
				l.tls = true
			} else {
				l.tls = false
			}
			if strings.Contains(u.Host, ":") {
				h := strings.Split(u.Host, ":")
				l.host = h[0]
				p, err := strconv.Atoi(h[1])
				if err != nil {
					log.Fatal(err)
				}
				l.port = p
			} else {
				l.host = u.Host
				if l.tls {
					l.port = 686
				} else {
					l.port = 389
				}
			}
		case v[0] == "base":
			l.base = v[1]
		case v[0] == "pam_authz_search":
			if strings.Contains(v[1], "$username") {
				l.filter = strings.Replace(v[1], "$username", "%s", 1)
			} else {
				l.filter = v[1]
			}
		default:
			if l.filter == "" {
				l.filter = defaultFilter
			}
		}
	}
}
