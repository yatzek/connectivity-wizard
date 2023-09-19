package main

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

var fqdnToCheck = "potato.com"

func main() {
	log.Printf("dns check for %s\n", fqdnToCheck)

	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		 log.Fatal("error loading /etc/resolv.conf", err)
	}

	dnsHostPort := net.JoinHostPort(config.Servers[0], config.Port)

	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(fqdnToCheck), dns.TypeA)

	client := new(dns.Client)
	resp, rtt, err := client.Exchange(message, dnsHostPort)
	if err != nil {
		log.Fatal("error sending dns query", err)
 	}

	// print response
	log.Printf("round trip time: %v\n", rtt)
	log.Printf("%v\n", resp)
}