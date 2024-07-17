package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

// Initializes persistent connection to AD Domain Server using ldaps
func PersistentConnectionInit(adServer string, baseDN string, bindDN string, bindPassword string) {
	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Note: This should be set to false in production and proper certificates should be used
	}

	// Connect to the LDAP server using DialURL
	l, err := ldap.DialURL(adServer, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer l.Close()

	// Bind to the server
	err = l.Bind(bindDN, bindPassword)
	if err != nil {
		log.Fatalf("Failed to bind: %v", err)
	}

	// Syncrepl configuration search with notify, not persistent connection
	syncConfig := &ldap.SearchRequest{
		BaseDN:     baseDN,
		Scope:      ldap.ScopeWholeSubtree,
		Filter:     "(objectClass=*)",
		Attributes: []string{"dn", "cn"}, //FIXME Update to other fields
	}

	// Context for ldap search
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mode := ldap.SyncRequestModeRefreshAndPersist
	var cookie []byte = nil

	// Start the Syncrepl session
	syncStream := l.Syncrepl(ctx, syncConfig, 64, mode, cookie, false)

	// Handle Syncrepl updates
	for syncStream.Next() {
		entry := syncStream.Entry()
		if entry != nil {
			fmt.Printf("%s has DN %s\n", entry.GetAttributeValue("cn"), entry.DN) //FIXME Update data structure
		}
		controls := syncStream.Controls()
		if len(controls) != 0 {
			fmt.Printf("%s", controls)
		}
	}
	if err := syncStream.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Use the LDAPS URL scheme to ensure a secure connection
	PersistentConnectionInit("ldaps://localhost:636", "\"dc=example,dc=com\", testUser", "cn=admin,dc=example,dc=com", "password1!")
}
