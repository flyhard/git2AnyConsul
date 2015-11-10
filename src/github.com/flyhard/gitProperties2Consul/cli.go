package main

import "flag"

type CliParameters struct {
	host     string
	port     int
	dataDir  string
	repo     string
	branch   string
	interval int
}

func parseCli() (parameters CliParameters) {
	flag.StringVar(&parameters.host, "host", "127.0.0.1", "Address of consul server")

	flag.IntVar(&parameters.port, "port", 8500, "consul port")

	flag.IntVar(&parameters.interval, "interval", 10, "The interval between polls of the repo")

	flag.StringVar(&parameters.dataDir, "dataDir", "./data", "The location of the data store")

	flag.StringVar(&parameters.repo, "repo", "", "The location of the git repo")

	flag.StringVar(&parameters.branch, "branch", "master", "The branch to use")

	flag.Parse()
	return
}
