package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

func waitForConsul(host string, port int) (kv *api.KV) {

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
	leader = ""
	for retry < 5 && leader == "" {
		leader, err := status.Leader()
		if err != nil {
			Error.Fatal("failed to get the Consul leader")
		} else {
			if leader != "" {
				Info.Print("Connected to Consul with leader:", leader)
				break
			}
		}
		time.Sleep(1 * time.Second)
		retry++
	}

	// Get a handle to the KV API
	kv = client.KV()
	return
}

func storeData(kv *api.KV, key string, value []byte) {
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		Error.Fatal("Failed reading from Consul for key:", key, "Error:", err)
	}
	if pair == nil || string(pair.Value) != string(value) {
		var p *api.KVPair
		if pair == nil {
			p = &api.KVPair{Key: key, Value: value}
		} else {
			p = &api.KVPair{Key: key, Value: value, ModifyIndex: pair.ModifyIndex}
		}
		success, meta, err := kv.CAS(p, nil)
		if err != nil {
			Error.Fatal(err)
		}
		if !success {
			storeData(kv, key, value)
		} else {
			Trace.Print("Request took: ", meta.RequestTime.Seconds(), " seconds")
			Info.Print("Updated key '", key, "' to '", string(value), "'")
		}
	} else {
		Trace.Print("not updating key", key)
	}
}
