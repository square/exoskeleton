.PHONY: CHANGELOG.md
CHANGELOG.md: .hermit/rust/bin/git-cliff
	.hermit/rust/bin/git-cliff -o CHANGELOG.md

.hermit/rust/bin/git-cliff:
	bin/cargo install git-cliff
