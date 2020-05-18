VERSION=v0.2.0

reltest:
	@goreleaser --snapshot --skip-publish --rm-dist

release:
	@git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
	@goreleaser --rm-dist
