install: 
	go install ./...


lint:
	golangci-lint run

tag:
	git tag $(TAG)
	git push origin $(TAG)

tag-delete:
	git tag --delete $(TAG)
	git push --delete origin $(TAG)
	gh release delete $(TAG) -y
