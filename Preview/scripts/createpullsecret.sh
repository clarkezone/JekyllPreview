export CR_PAT_BASE64=$(echo -n 'clarkezone:'$CR_PAT | base64)
echo '{"auths":{"ghcr.io":{"auth":"'$CR_PAT_BASE64'"}}}' \
| kubectl create secret generic ghcr-pat --type=kubernetes.io/dockerconfigjson --from-file=.dockerconfigjson=/dev/stdin
