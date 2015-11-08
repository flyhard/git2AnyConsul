# git2AnyConsul
Clones and polls a git repo and pushes its contents into Consul.

This project is implementing the same basic idea as git2Consul, but does not need a local consul agent.

## completed features
* reading a GIT repo
* writing a key to Consul if different from current value
* reading JSON files as part of the tree

## incomplete features
* Polling the Git repo
* Reading files that are not in JSON format into Consul Values
