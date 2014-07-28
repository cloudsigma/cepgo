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
// /dev/ttyS1
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
// returns it
func fetchViaSerialPort(key string) (interface{}, error) {
	var result interface{}

	config := &serial.Config{Name: SerialPort, Baud: Baud}
	connection, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(requestPattern, key)
	if _, err := connection.Write([]byte(query)); err != nil {
		return result, err
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
