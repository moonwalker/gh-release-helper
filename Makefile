VERSION=v0.4.0

reltest:
	@goreleaser --snapshot --skip-publish --rm-dist

release:
	@git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
	@goreleaser --rm-dist
