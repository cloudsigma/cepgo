package cepgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/tarm/goserial"
)

const requestPattern = "<\n%s\n>"

var (
	SerialPort string = "/dev/ttyS1"
	Baud       int    = 115200
)

// Sets the serial port. If the operating system is windows CloudSigma's server
// context is at COM2 port, otherwise (linux, freebsd, darwin) the port is
// /dev/ttyS1
func init() {
	if runtime.GOOS == "windows" {
		SerialPort = "COM2"
	}
}

func fetchViaSerialPort(key string) (interface{}, error) {
	var result interface{}

	config := &serial.Config{Name: SerialPort, Baud: Baud}
	connection, err := serial.OpenPort(config)
	if err != nil {
		return result, err
	}

	query := fmt.Sprintf(requestPattern, key)
	if _, err := connection.Write([]byte(query)); err != nil {
		return result, err
	}

	answer, err := ioutil.ReadAll(connection)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(answer, result)
	return result, err
}

type Cepgo struct {
	fetcher func(string) (interface{}, error)
}

func NewCepgo() *Cepgo {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchViaSerialPort
	return cepgo
}

func (c *Cepgo) Key(key string) (interface{}, error) {
	return c.fetcher(key)
}

func (c *Cepgo) All() (interface{}, error) {
	return c.fetcher("")
}

func (c *Cepgo) Meta() (interface{}, error) {
	return c.fetcher("/meta/")
}

func (c *Cepgo) GlobalContext() (interface{}, error) {
	return c.fetcher("/global_context/")
}
