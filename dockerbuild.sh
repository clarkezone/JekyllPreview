docker build . \
 --build-arg ARG_TEST_JEKPREV_REPO_NOAUTH=https://gitlab.com/clarkezone/jekyllpreviewtestnoauth.git \
 --build-arg ARG_TEST_JEKPREV_LOCALDIR=test_repo \
 --build-arg ARG_TEST_JEKPREV_BRANCHSWITCH=test \
 --build-arg ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH=https://gitlab.com/clarkezone/jekyllpreviewtestauth.git \
 --build-arg ARG_TEST_JEKPREV_SECURECLONEPW=H72fqkQh3HVekAtnfhbX \
 --add-host=master.localhost:127.0.0.1 \
 --add-host=pepper.localhost:127.0.0.1 \
 --tag clarkezone/jekpreview:sd
