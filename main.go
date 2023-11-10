package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var k = koanf.New(".")

func main() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.StringSlice("conf", []string{"config/config_development.yml"}, "path to one or more .yml config files")

	f.StringP("file", "f", "test_data/sample.mp3", "path to one or more files")

	f.Parse(os.Args[1:])

	configFiles, _ := f.GetStringSlice("conf")
	for i, c := range configFiles {
		if err := k.Load(file.Provider(c), yaml.Parser()); err != nil {
			log.Fatalf("error loading [%d] config: %v", i, err)
		}
	}

	filePath, _ := f.GetString("file")

	config := getConfig(k)

	connectToVosk(config, filePath)

}

var m Message

func connectToVosk(c *Config, filePath string) {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.Vosk.Host, c.Vosk.Port),
		Path:   "",
	}

	log.Printf("connecting to %s/n", u.String())

	con, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer con.Close()

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buf := make([]byte, 8000)
		dat, err := f.Read(buf)

		if dat == 0 && err == io.EOF {
			err = con.WriteMessage(websocket.TextMessage, []byte("end"))
			if err != nil {
				log.Fatal("write:", err)
			}

			break
		}

		err = con.WriteMessage(websocket.BinaryMessage, buf)
		if err != nil {
			log.Fatal("write:", err)
		}

		_, _, err = con.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
	}

	_, msg, err := con.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}

	con.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	err = json.Unmarshal(msg, &m)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("message: %+v\n", m)
}
