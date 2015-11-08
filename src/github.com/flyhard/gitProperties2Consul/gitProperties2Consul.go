package main

import (
	"encoding/json"
	"github.com/VictorLowther/go-git/git"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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
			Trace.Print("Processing dir:", name+"/")
			processDir(basename, name, kv)
		} else {
			Trace.Print("Processing file", name)
			isJson := loadJson(basename, name, kv)
			if !isJson {
				p := &api.KVPair{Key: name, Value: []byte("TODO")}
				_, err = kv.Put(p, nil)
				if err != nil {
					Error.Fatal(err)
				}
			}
		}
	}
}

func loop(dataDir string, kv *api.KV, repo *git.Repo) {
	updateRepo(repo)
	processDir(dataDir, "", kv)
	time.Sleep(10 * time.Second)
	loop(dataDir, kv, repo)
}

func main() {
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	host, port, dataDir, repo := parseCli()

	kv := waitForConsul(host, port)

	repository := aquireGitRepo(repo, dataDir)
	loop(dataDir, kv, repository)
}
