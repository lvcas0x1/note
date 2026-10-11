---
date: 2026-09-22
---
##### Change commit user for different workpath

```shell
# .gitconfig
[user]
	email = email-address@company.com
	name = username

[includeif "gitdir:~/Documents/Code/workpath/"]
	path = ~/.gitconfig-work
```
