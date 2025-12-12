.PHONY: CHANGELOG.md
CHANGELOG.md:
	bin/git-cliff -o CHANGELOG.md
