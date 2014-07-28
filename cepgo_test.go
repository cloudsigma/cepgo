package cepgo

import (
	"errors"
	"strings"
	"testing"
)

func fetchMock(key string) (interface{}, error) {
	context := map[string]interface{}{
		"cpu": 1000,
		"cpus_instead_of_cores": false,
		"global_context": map[string]interface{}{
			"some_global_key": "some_global_val",
		},
		"mem": 1073741824,
		"meta": map[string]interface{}{
			"ssh_public_key": "ssh-rsa AAAAB3NzaC1yc2E.../hQ5D5 john@doe",
		},
		"name":         "test_server",
		"smp":          1,
		"tags":         []string{"much server", "very performance"},
		"uuid":         "65b2fb23-8c03-4187-a3ba-8b7c919e8890",
		"vnc_password": "9e84d6cb49e46379",
	}
	if key == "" {
		return context, nil
	} else {
		result, ok := context[strings.Trim(key, "/")]
		if ok {
			return result, nil
		} else {
			return nil, errors.New("No such key in the context")
		}
	}
}

func TestAll(t *testing.T) {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchMock

	result, err := cepgo.All()
	if err != nil {
		t.Error(err)
	}

	for _, key := range []string{"meta", "name", "uuid", "global_context"} {
		if _, ok := result.(map[string]interface{})[key]; !ok {
			t.Errorf("%s not in all keys", key)
		}
	}
}

func TestKey(t *testing.T) {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchMock

	result, err := cepgo.Key("uuid")
	if err != nil {
		t.Error(err)
	}

	if _, ok := result.(string); !ok {
		t.Error("Fetching the uuid did not return a string")
	}
}

func TestMeta(t *testing.T) {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchMock

	result, err := cepgo.Meta()
	if err != nil {
		t.Error(err)
	}

	if meta, ok := result.(map[string]interface{}); !ok {
		t.Error("Fetching the meta did not return a map[string]interface{}")
	} else if _, ok := meta["ssh_public_key"]; !ok {
		t.Error("ssh_public_key is not in the meta")
	}
}

func TestGlobalContext(t *testing.T) {
	cepgo := new(Cepgo)
	cepgo.fetcher = fetchMock

	result, err := cepgo.GlobalContext()
	if err != nil {
		t.Error(err)
	}

	if meta, ok := result.(map[string]interface{}); !ok {
		t.Error("Fetching the global context did not return a map[string]interface{}")
	} else if _, ok := meta["some_global_key"]; !ok {
		t.Error("some_global_key is not in the global context")
	}
}
