package castor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// TODO: use repo directory instad of WD
// git rev-parse --git-dir
// outputs .git if root of repo
// outputs absolute path otherwise

var castorfile string

func init() {
	cur, err := user.Current()
	if err != nil {
		log.Fatal(err) // no intention to handle this case
	}

	castorfile = filepath.Join(cur.HomeDir, ".castor")
}

func wipBranch() (string, error) {
	dir, err := repoDir()
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(castorfile)
	if err != nil {
		if os.IsNotExist(err) {
			return "master", nil
		}
		return "", err
	}

	branches := map[string]string{}

	err = json.Unmarshal(b, &branches)
	if err != nil {
		return "", err
	}

	branch, ok := branches[dir]
	if !ok {
		return "master", nil
	}

	return branch, nil
}

// TODO: handle missing files properly
func setWipBranch() error {
	branch, err := currentBranch()
	if err != nil {
		return err
	}

	dir, err := repoDir()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(castorfile)
	if err != nil {
		if os.IsNotExist(err) {
			b = []byte("{}")
		} else {
			fmt.Println("does exist but can't open")
			return err
		}
	}

	fmt.Println("did not end")

	branches := map[string]string{}

	err = json.Unmarshal(b, &branches)
	if err != nil {
		fmt.Println("unmarshall: ", err.Error())
		return err
	}

	branches[dir] = branch

	b, err = json.Marshal(branches)
	if err != nil {
		fmt.Println("marshall: ", err.Error())
		return err
	}

	err = ioutil.WriteFile("~/.castor", b, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(castorfile)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.Write(b)
			return err
		}
		return err
	}

	return nil
}
