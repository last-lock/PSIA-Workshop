package plai_ldap_plugin

import (
	"context"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

// Initializes persistent connection to AD Domain Server using ldap
//
// Parameters:
//   - a (int): The first integer.
//   - b (int): The second integer.
//
// Returns:
//   - int: The sum of a and b.
//
// Example:
//   result := Sum(5, 3)
//   fmt.Println(result) // Output: 8
func Persistentconnectioninit(adServer string, baseDN string, bindDN string, bindPassword string) {
    // // Active Directory server details
    // adServer := "ldap://localhost:389"
    // baseDN := "dc=example,dc=com"
    // bindDN := "cn=admin,dc=example,dc=com"
    // bindPassword := "password"

    // Connect to the Active Directory server
    l, err := ldap.DialURL(adServer)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()

    // Bind to the server
    err = l.Bind(bindDN, bindPassword)
    if err != nil {
        log.Fatal(err)
    }

    // Syncrepl configuration search with notify, not persistent connection
    syncConfig := &ldap.SearchRequest{
        BaseDN:       baseDN,
        Scope:        ldap.ScopeWholeSubtree,
        Filter:       "(objectClass=*)",
        Attributes:   []string{"dn", "cn"}, //FIXME Update to other fields 
    }

    // Context for ldap search
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    mode := ldap.SyncRequestModeRefreshAndPersist
    var cookie []byte = nil

    // Start the Syncrepl session
    syncStream:= l.Syncrepl(ctx, syncConfig, 64, mode, cookie, false)

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







