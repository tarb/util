package www

import (
	"net/http"
	"sync"
)

var cliMut sync.RWMutex
var client *http.Client = http.DefaultClient

// SetClient sets which http.Client will be making the custom
// http requests. Uses http.DefaultClient by default
func SetClient(cli *http.Client) {
	cliMut.Lock()
	client = cli
	cliMut.Unlock()
}

// getClient returns the current http.Client to make requests
func getClient() *http.Client {
	defer cliMut.RUnlock()
	cliMut.RLock()

	return client
}
