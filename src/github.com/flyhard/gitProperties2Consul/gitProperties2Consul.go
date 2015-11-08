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
	"encoding/json"
	"github.com/VictorLowther/go-git/git"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogging(
traceHandle io.Writer,
infoHandle io.Writer,
warningHandle io.Writer,
errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate | log.Ltime | log.Lshortfile)
}

func storeData(kv *api.KV, key string, value []byte) {
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		Error.Fatal("Failed reading from Consul for key:", key, "Error:", err)
	}
	if pair == nil || string(pair.Value) != string(value) {
		p := &api.KVPair{Key:key, Value:value}
		_, err = kv.Put(p, nil)
		if err != nil {
			Error.Fatal(err)
		}
		Info.Print("Updated key ", key, "to", value)
	} else {
		Trace.Print("not updating key", key)
	}
}

func processJson(path string, m map[string]interface{}, kv *api.KV) {
	for key, value := range m {
		localPath := path + "/" + key
		switch vv := value.(type) {
		case string:
			Trace.Print(localPath, "is string", vv)
			storeData(kv, localPath, []byte(value.(string)))
		case int:
			Trace.Print(localPath, "is int", vv)
		case []interface{}:
			Trace.Print(localPath, "is an array:")
			for i, u := range vv {
				Trace.Print(i, u)
			}
		case map[string]interface{}:
			Trace.Print(localPath, "is a map")
			processJson(localPath, value.(map[string]interface{}), kv)
		default:
			Error.Print(key, "is of a type I don't know how to handle: ", vv)
		}
	}
}

func loadJson(basename string, filename string, kv *api.KV) bool {
	b, err := ioutil.ReadFile(basename + "/" + filename)
	if err != nil {
		Error.Fatal("File reading failed:", err)
	}
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		Error.Print("Unmarshallingfailed:", err)
		return false
	}
	if f == nil {
		Warning.Print("no readable data in file", filename)
		return false
	}
	m := f.(map[string]interface{})
	dotIndex := strings.LastIndex(filename, ".")
	slashIndex := strings.LastIndex(filename, "/")
	if dotIndex > slashIndex {
		filename = filename[0:dotIndex]
	}

	processJson(filename, m, kv)
	return true
}

func processDir(basename string, dirname string, kv *api.KV) {
	files, err := ioutil.ReadDir(basename + "/" + dirname)
	if err != nil {
		Error.Fatal("Directory reading failed:", err)
	}
	for _, element := range files {
		name := element.Name()
		if dirname != "" {
			name = dirname + "/" + name
		}
		if strings.HasPrefix(element.Name(), ".") {
			continue
		}
		if element.IsDir() {
			Trace.Print("Processing dir:", name + "/")
			processDir(basename, name, kv)
		} else {
			Trace.Print("Processing file", name)
			isJson := loadJson(basename, name, kv)
			if !isJson {
				p := &api.KVPair{Key:name, Value:[]byte("TODO")}
				_, err = kv.Put(p, nil)
				if err != nil {
					Error.Fatal(err)
				}
			}
		}
	}
}

func main() {
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

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

	status := client.Status()
	retry := 0
	var leader string
	for retry < 5 && leader == ""{
		leader, err = status.Leader()
		if err != nil {
			Error.Fatal("failed to get the Consul leader")
		}else {
			Info.Print("Connected to Consul with leader:", leader)
		}
		time.Sleep(1 * time.Second)
		retry ++
	}

	// Get a handle to the KV API
	kv := client.KV()
	r, err := git.Clone("https://github.com/flyhard/testData", dataDir)
	if err != nil {
		Error.Fatal(err)
	}
	clean, _ := r.IsClean()
	if clean {
		Info.Print("Done cloning. repo is clean now.")
	}
	processDir(dataDir, "", kv)

}
