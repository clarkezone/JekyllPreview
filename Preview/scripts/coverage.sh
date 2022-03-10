export TEST_JEKPREV_REPO_NOAUTH=https://github.com/clarkezone/JekyllPreview.git
export TEST_JEKPREV_LOCALDIR=/tmp/jekpreview_test
export TEST_JEKPREV_BRANCHSWITCH=BugFix
export TEST_JEKPREV_SECURE_REPO_NOAUTH=true
export TEST_JEKPREV_SECURECLONEPW=unused
go test -covermode=count -coverprofile=count.out fmt
