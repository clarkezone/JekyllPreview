curl -X POST http://0.0.0.0:8090/postreceive \
   -H 'Content-Type: application/json' -H 'X-GitLab-Event: Push Hook' \
   -d @webhook.json
