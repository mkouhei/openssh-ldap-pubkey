package main

import (
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
	sshPublicKeyName = "sshPublicKey"
)

type ldapEnv struct {
	host   string
	port   int
	base   string
	filter string
	tls    bool
	skip   bool
	uid    string
}

func (l *ldapEnv) argparse(args []string) {
	if len(args) == 0 {
		args = os.Args
	}
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	h := flags.String("host", l.host, "LDAP server host")
	p := flags.Int("port", l.port, "LDAP server port")
	b := flags.String("base", l.base, "search base")
	f := flags.String("filter", l.filter, "search filter")
	t := flags.Bool("tls", l.tls, "LDAP connect over TLS")
	s := flags.Bool("skip", false, "Insecure skip verify")
	flags.Parse(args[1:])

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

	if len(flags.Args()) != 1 {
		log.Fatal("Specify username")
	}
	l.uid = flags.Args()[0]
}

func isAddr(host string) bool {
	return !(net.ParseIP(host).To4() == nil && net.ParseIP(host).To16() == nil)
}

func (l *ldapEnv) connect() *ldap.Conn {
	c, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (l *ldapEnv) connectTLS() *ldap.Conn {
	certs := *x509.NewCertPool()
	tlsConfig := &tls.Config{
		RootCAs: &certs,
	}

	if isAddr(l.host) || l.skip {
		tlsConfig.InsecureSkipVerify = true
	}

	c, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.host, l.port), tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func simpleBind(c *ldap.Conn) {
	bindRequest := ldap.NewSimpleBindRequest("", "", nil)
	_, err := c.SimpleBind(bindRequest)
	if err != nil {
		log.Fatal(err)
	}
}

func (l *ldapEnv) search(c *ldap.Conn) []*ldap.Entry {
	searchRequest := ldap.NewSearchRequest(
		l.base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf(l.filter, l.uid), []string{sshPublicKeyName}, nil)
	sr, err := c.Search(searchRequest)
	if err != nil {
		log.Fatal(err, sshPublicKeyName)
	}
	return sr.Entries
}

func printPubkey(entries []*ldap.Entry) {
	if len(entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}

	if len(entries[0].GetAttributeValues("sshPublicKey")) == 0 {
		log.Fatal("User does not use ldapPublicKey.")
	}
	for _, pubkey := range entries[0].GetAttributeValues("sshPublicKey") {
		fmt.Println(pubkey)
	}
}

func main() {
	l := &ldapEnv{}
	l.loadNslcdConf()
	l.argparse([]string{})

	c := &ldap.Conn{}
	if l.tls {
		c = l.connectTLS()
	} else {
		c = l.connect()
	}
	defer c.Close()

	simpleBind(c)
	printPubkey(l.search(c))
}
