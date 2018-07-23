package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"crypto/tls"
	"crypto/x509"

	"gopkg.in/ldap.v2"
)

const (
	version          = "0.1.2"
	sshPublicKeyName = "sshPublicKey"
)

type ldapEnv struct {
	host   string
	port   int
	base   string
	filter string
	tls    bool
	skip   bool
	debug  bool
	uid    string
	binddn string
	bindpw string
}

var (
	ver        string
	errVersion = errors.New("show version")

	license = `openssh-ldap-pubkey %s

Copyright (C) 2015 Kouhei Maeda
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
This is free software, and you are welcome to redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`
)

func (l *ldapEnv) argparse(args []string, ver string) error {
	if len(args) == 0 {
		args = os.Args
	}
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	h := flags.String("host", l.host, "LDAP server host")
	p := flags.Int("port", l.port, "LDAP server port")
	b := flags.String("base", l.base, "search base")
	f := flags.String("filter", l.filter, "search filter")
	t := flags.Bool("tls", l.tls, "LDAP connect over TLS")
	s := flags.Bool("skip", l.skip, "Insecure skip verify")
	v := flags.Bool("version", false, "show version")
	d := flags.Bool("debug", false, "debug mode")
	flags.Parse(args[1:])

	if *v {
		fmt.Printf(license, ver)
		return errVersion
	}
	if l.host != *h {
		l.host = *h
	}
	if l.port != *p {
		l.port = *p
	}
	if l.base != *b {
		l.base = *b
	}
	if l.filter != *f {
		l.filter = *f
	}
	if l.tls != *t {
		l.tls = *t
	}
	if l.skip != *s {
		l.skip = *s
	}
	if l.debug != *d {
		l.debug = *d
	}

	if len(flags.Args()) != 1 {
		return errors.New("Specify username")
	}
	l.uid = flags.Args()[0]
	return nil
}

func isAddr(host string) bool {
	return !(net.ParseIP(host).To4() == nil && net.ParseIP(host).To16() == nil)
}

func (l *ldapEnv) getAddr() string {
	addrs, err := net.LookupHost(l.host)
	if err != nil {
		return l.host
	}
	if net.ParseIP(addrs[0]).To16() != nil {
		return fmt.Sprintf("[%s]", addrs[0])
	}
	return addrs[0]
}

func (l *ldapEnv) connect() (*ldap.Conn, error) {
	host := l.host
	if !isAddr(host) {
		host = l.getAddr()
	}
	return ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, l.port))
}

func (l *ldapEnv) connectTLS() (*ldap.Conn, error) {
	certs := *x509.NewCertPool()
	tlsConfig := &tls.Config{
		RootCAs: &certs,
	}

	if isAddr(l.host) || l.skip {
		tlsConfig.InsecureSkipVerify = true
	}
	host := l.host
	if !isAddr(l.host) {
		tlsConfig.ServerName = l.host
		host = l.getAddr()
	}
	return ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", host, l.port), tlsConfig)
}

func logging(err error) {
	if err == errVersion {
		os.Exit(0)
	} else if err != nil {
		log.Fatal(err)
	}
}

func simpleBind(c *ldap.Conn, l *ldapEnv) error {
	bindRequest := ldap.NewSimpleBindRequest(l.binddn, l.bindpw, nil)
	_, err := c.SimpleBind(bindRequest)
	return err
}

func (l *ldapEnv) search(c *ldap.Conn) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		l.base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf(l.filter, l.uid), []string{sshPublicKeyName}, nil)
	sr, err := c.Search(searchRequest)
	return sr.Entries, err
}

func printPubkey(entries []*ldap.Entry) error {
	if len(entries) != 1 {
		return errors.New("User does not exist or too many entries returned")
	}

	if len(entries[0].GetAttributeValues("sshPublicKey")) == 0 {
		return errors.New("User does not use ldapPublicKey")
	}
	for _, pubkey := range entries[0].GetAttributeValues("sshPublicKey") {
		fmt.Println(pubkey)
	}
	return nil
}

func main() {
	l := &ldapEnv{}
	l.loadNslcdConf()
	var err error
	var entries []*ldap.Entry
	if ver == "" {
		ver = version
	}
	logging(l.argparse([]string{}, ver))
	c := &ldap.Conn{}
	if l.debug {
		var bindpw = ""
		if l.bindpw != "" {
			bindpw = "<bindpw can found in nslcd.conf>"
		}
		log.Printf("[debug] host  : %s\n", l.host)
		log.Printf("[debug] port  : %d\n", l.port)
		log.Printf("[debug] tls	  : %v\n", l.tls)
		log.Printf("[debug] base  : %s\n", l.base)
		log.Printf("[debug] skip  : %v\n", l.skip)
		log.Printf("[debug] filter: %s\n", l.filter)
		log.Printf("[debug] uid	  : %s\n", l.uid)
		log.Printf("[debug] binddn: %s\n", l.binddn)
		log.Printf("[debug] bindpw: %s\n", bindpw)
	}
	if l.tls {
		c, err = l.connectTLS()
		logging(err)
	} else {
		c, err = l.connect()
		logging(err)
	}
	defer c.Close()

	logging(simpleBind(c, l))
	entries, err = l.search(c)
	logging(err)
	logging(printPubkey(entries))
}
