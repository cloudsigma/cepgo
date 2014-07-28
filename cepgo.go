// Cepko implements easy-to-use communication with CloudSigma's VMs through a
// virtual serial port without bothering with formatting the messages properly
// nor parsing the output with the specific and sometimes confusing shell tools
// for that purpose.
//
// Having the server definition accessible by the VM can ve useful in various
// ways. For example it is possible to easily determine from within the VM,
// which network interfaces are connected to public and which to private
// network. Another use is to pass some data to initial VM setup scripts, like
// setting the hostname to the VM name or passing ssh public keys through
// server meta.
//
// Example usage:
//
//   package main
//
//   import (
//           "fmt"
//
//           "github.com/cloudsigma/cepgo"
//   )
//
//   func main() {
//           c := cepgo.NewCepgo()
//           result, err := c.Meta()
//           if err != nil {
//                   panic(err)
//           }
//           fmt.Printf("%#v", result)
//   }
//
// Output:
//
//   map[string]interface {}{
//   	"optimize_for":"custom",
//   	"ssh_public_key":"ssh-rsa AAA...",
//   	"description":"[...]",
//   }
//
// For more information take a look at the Server Context section API Docs:
// http://cloudsigma-docs.readthedocs.org/en/latest/server_context.html
package cepgo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/tarm/goserial"
)

const (
	requestPattern = "<\n%s\n>"
	EOT            = '\x04'
)

var (
	SerialPort string = "/dev/ttyS1"
	Baud       int    = 115200
)

// Sets the serial port. If the operating system is windows CloudSigma's server
// context is at COM2 port, otherwise (linux, freebsd, darwin) the port is
// /dev/ttyS1.
func init() {
	if runtime.GOOS == "windows" {
		SerialPort = "COM2"
	}
}

// The default fetcher makes the connection to the serial port,
// writes given query and reads until the EOT symbol. Then tries
// to unmarshal the result to json and returns it. If the
// unmarshalling fails it's safe to assume the result it's just
// a string (e.g. the Key() method has been executed) and
// returns it.
func fetchViaSerialPort(key string) (interface{}, error) {
	var result interface{}

	config := &serial.Config{Name: SerialPort, Baud: Baud}
	connection, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(requestPattern, key)
	if _, err := connection.Write([]byte(query)); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(connection)
	answer, err := reader.ReadBytes(EOT)
	if err != nil {
		return nil, err
	}

	usefulAnswer := answer[0 : len(answer)-1]
	err = json.Unmarshal(usefulAnswer, &result)
	if err != nil {
		return string(usefulAnswer), nil
	}
	return result, nil
}

// Queries to the serial port can be executed only from instance of this type.
// The result from each of them can be either a map[string]interface or a
// single in case of single value is returned.
type Cepgo struct {
	fetcher func(string) (interface{}, error)
}

// Creates a Cepgo instance with the default serial port fetcher
func NewCepgo() *Cepgo {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchViaSerialPort
	return cepgo
}

// Creates a Cepgo instance with custom fetcher fetcher
func NewCepgoFetcher(fetcher func(string) (interface{}, error)) *Cepgo {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetcher
	return cepgo
}

// Fetches a single key
func (c *Cepgo) Key(key string) (interface{}, error) {
	return c.fetcher(key)
}

// Fetches all the server context. Equivalent of c.Key("")
func (c *Cepgo) All() (interface{}, error) {
	return c.fetcher("")
}

// Fetches only the object meta field. Equivalent of c.Key("/meta/")
func (c *Cepgo) Meta() (interface{}, error) {
	return c.fetcher("/meta/")
}

// Fetches only the global context. Equivalent of c.Key("/global_context/")
func (c *Cepgo) GlobalContext() (interface{}, error) {
	return c.fetcher("/global_context/")
}
