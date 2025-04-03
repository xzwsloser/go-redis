package config

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ServerProperties defines the global config properties of the Redis Server
type ServerProperties struct {
	Bind           string   `cfg:"bind"`
	Port           int      `cfg:"port"`
	AppendOnly     bool     `cfg:"appendOnly"`
	AppendFilename string   `cfg:"appendFilename"`
	MaxClients     int      `cfg:"maxClients"`
	RequirePass    string   `cfg:"requirePass"`
	Databases      string   `cfg:"databases"`
	Peers          []string `cfg:"peers"`
	Self           string   `cfg:"self"`
}

var Properties *ServerProperties

// the default properties of server
func init() {
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       7379,
		AppendOnly: false,
	}
}

func parse(src io.Reader) *ServerProperties {
	config := &ServerProperties{}

	// read the config file
	rawMap := make(map[string]string)
	// use scanner to read line by line
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		// the comments , skip
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 {
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}

		if err := scanner.Err(); err != nil {
			// TODO: using the logger to record the error status
			return config
		}

		// using reflect to reject the properties
		t := reflect.TypeOf(config)
		v := reflect.ValueOf(config)
		n := t.Elem().NumField()
		for i := 0; i < n; i++ {
			field := t.Elem().Field(i)
			fieldValue := v.Elem().Field(i)
			key, ok := field.Tag.Lookup("cfg")
			if !ok {
				key = field.Name
			}
			value, ok := rawMap[strings.ToLower(key)]
			if ok {
				switch field.Type.Kind() {
				case reflect.String:
					fieldValue.SetString(value)
				case reflect.Int:
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err == nil {
						fieldValue.SetInt(intValue)
					}
				case reflect.Bool:
					boolValue := "yes" == value
					fieldValue.SetBool(boolValue)
				case reflect.Slice:
					if field.Type.Elem().Kind() == reflect.String {
						slices := strings.Split(value, ",")
						fieldValue.Set(reflect.ValueOf(slices))
					}
				}
			}
		}
	}
	return config
}

func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	Properties = parse(file)
}
