## Flow for working with `gipr`

**Listing active my PRs**

```
$ gipr prs --my
- #3840 [HA-1119] Add Transaction History ... frontend to-review 2 change requests
- #3577 [HA-995] Add "Payout Settings" pages  frontend to-accept accepted
```

**Listing PRs assigned to me**

```
$ gipr prs --assigned
- #3887 Add check for tester notes                               changes requested - accepted by 5 people
- #3840 [HA-1119] Add Transaction History page t... frontend wip
- #3789 [HA-1039] Setup Cypress and create login... frontend     changes requested - waiting review
```

**Open a PR in the browser/editor**

```
$ gipr open 3840
```

**Checkingout a PR's branch**

NOTE: _git commands here are only used to show status of the repo_.

```
$ git status                 
On branch feature/ha-1119
Your branch is up-to-date with 'origin/feature/ha-1119'.
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   package.json

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	static/js/myredux/handleActions.ts

no changes added to commit (use "git add" and/or "git commit -a")
```

```
$ git branch                                                    
  develop
* feature/ha-1119
```

```
$ gipr view 3887
Your changes in branch `feature/ha-1119` are safe.
Switched to `docs/check-test-stage-needed`

Here's a list of what changed

static/js/app.js +10 -20
```

```
$ git branch                                                    
  develop
  feature/ha-1119
* docs/check-test-stage-needed 
```

```
$ git status
On branch feature/ha-1119
Your branch is ahead of 'origin/feature/ha-1119' by 1 commit.
  (use "git push" to publish your local commits)
nothing to commit, working directory clean
```

**Checkout the WIP branch**

```
$ gipr goback
Going back to `feature/ha-1119`
Applying your work in progress...
```

```
$ git status                 
On branch feature/ha-1119
Your branch is up-to-date with 'origin/feature/ha-1119'.
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   package.json

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	static/js/myredux/handleActions.ts
```
