module temp.com/JekyllBlogPreview

go 1.12

require (
	github.com/clarkezone/hookserve v0.0.0-20200325144548-21b11caacc02
	github.com/go-git/go-git/v5 v5.0.0
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
	k8s.io/client-go v0.23.3
)

//replace github.com/clarkezone/go-execobservable => ../go-execobservable
