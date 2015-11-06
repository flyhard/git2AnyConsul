package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func processDir(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		Error.Fatal("File reading failed:", err)
	}
	for _, element := range files {
		name := dirname + "/" + element.Name()
		if strings.HasPrefix(element.Name(), ".") {
			continue
		}
		if element.IsDir() {
			Info.Print(name + "/")
			processDir(name)
		} else {
			Info.Print(name)
		}
	}
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	var host string
	flag.StringVar(&host, "host", "127.0.0.1", "Address of consul server")

	var port int
	flag.IntVar(&port, "port", 8500, "consul port")

	var dataDir string
	flag.StringVar(&dataDir, "dataDir", "./data", "The location of the data store")

	flag.Parse()

	// Get a new client
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", host, port)
	Info.Print("Config: ", config)
	client, err := api.NewClient(config)
	if err != nil {
		Error.Fatal(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	processDir(dataDir)

	// PUT a new KV pair
	p := &api.KVPair{Key: "foo", Value: []byte("test")}
	_, err = kv.Put(p, nil)
	if err != nil {
		Error.Fatal(err)
	}

	// Lookup the pair
	pair, _, err := kv.Get("foo", nil)
	if err != nil {
		Error.Fatal(err)
	}
	Info.Print("KV: ", pair)

}
