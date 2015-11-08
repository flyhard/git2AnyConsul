package main

import "github.com/VictorLowther/go-git/git"

func updateRepo(repo *git.Repo, branch string) {
	Info.Print("Updating GIT repo")
	repo.Checkout(branch)
	repo.Fetch([]string{"origin"})
	res, err := repo.Ref("remotes/origin/" + branch)
	if err != nil {
		Error.Fatal("Failed to get ref to origin/master: ", err)
	}
	currentRes, err := repo.CurrentRef()
	if err != nil {
		Error.Fatal("Failed to get ref to master: ", err)
	}
	err = currentRes.MergeWith(res)
	if err != nil {
		Error.Fatal("Failed to merge to master: ", err)
	}
}

func aquireGitRepo(repo string, dataDir string) (r *git.Repo) {
	Info.Print("Starting to clone repo '", repo, "'")
	r, err := git.Clone(repo, dataDir)
	if err != nil {
		Error.Fatal(err)
	}
	clean, _ := r.IsClean()
	if clean {
		Info.Print("Done cloning repo '", repo, "'. repo is clean now.")
	}
	return
}
