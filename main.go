package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/ldap.v2"
)

const (
	version          = "0.1.3"
	sshPublicKeyName = "sshPublicKey"
)

var (
	ver        string
	errVersion = errors.New("show version")

	license = `openssh-ldap-pubkey %s

Copyright (C) 2015-2018 Kouhei Maeda
License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
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
