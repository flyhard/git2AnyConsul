# git2AnyConsul
Clones and polls a git repo and pushes its contents into Consul.

This project is implementing the same basic idea as git2Consul, but does not need a local consul agent.

## Commandline
    Usage of bin/gitProperties2Consul:
      -dataDir string
            The location of the data store (default "./data")
      -host string
            Address of consul server (default "127.0.0.1")
      -port int
            consul port (default 8500)
      -repo string
            The location of the git repo (default "")
      -branch
            The branch to use (default "master")
      -interval
            The interval between polls of the repo (default 10)


## completed features
* reading a GIT repo
* writing a key to Consul if different from current value
* reading JSON files as part of the tree
* Polling the Git repo
* Reading files that are not in JSON format into Consul Values
* CLI parameter for selecting the branch
* Configurable poll interval

## incomplete features
