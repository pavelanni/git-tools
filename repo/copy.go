package repo

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Copy(repo, branch, dst string) error {
	fs, err := cloneToMemFs(repo, branch)
	if err != nil {
		return err
	}
	//printDir(fs, "./")
	err = copyDirFromMemFs(fs, ".", dst)
	if err != nil {
		return err
	}
	return nil
}

func cloneToMemFs(repo, branch string) (billy.Filesystem, error) {
	fs := memfs.New()
	// Git objects storer based on memory
	storer := memory.NewStorage()

	var refName plumbing.ReferenceName
	if branch == "" {
		refName = plumbing.HEAD
	} else {
		refName = plumbing.NewBranchReferenceName(branch)
	}
	// Clones the repository into the worktree (fs) and stores all the .git
	// content into the storer
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           repo,
		ReferenceName: refName,
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func printDir(fs billy.Filesystem, path string) {
	files, err := fs.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func copyFileFromMemFs(fs billy.Filesystem, src, dst string) error {
	srcFile, err := fs.Open(src)
	if err != nil {
		log.Println("error opening file", src)
		return err
	}
	if err = os.MkdirAll(filepath.Dir(os.ExpandEnv(dst)), 0755); err != nil {
		log.Println("error creating dir", filepath.Dir(os.ExpandEnv(dst)))
		return err
	}
	dstFile, err := os.Create(os.ExpandEnv(dst))
	if err != nil {
		return err
	}
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func copyDirFromMemFs(fs billy.Filesystem, src, dst string) error {
	srcStat, err := fs.Stat(src)
	if !srcStat.IsDir() {
		log.Println("not a dir", src)
		return err
	}
	files, err := fs.ReadDir(src)
	if err != nil {
		log.Println("error reading dir", src)
		return err
	}
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())
		if file.IsDir() {
			if err = copyDirFromMemFs(fs, srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err = copyFileFromMemFs(fs, srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
