package main

import "flag"

func parseCli() (host string, port int, dataDir string, repo string, branch string) {
	flag.StringVar(&host, "host", "127.0.0.1", "Address of consul server")

	flag.IntVar(&port, "port", 8500, "consul port")

	flag.StringVar(&dataDir, "dataDir", "./data", "The location of the data store")

	flag.StringVar(&repo, "repo", "", "The location of the git repo")

	flag.StringVar(&branch, "branch", "master", "The branch to use")

	flag.Parse()
	return
}
