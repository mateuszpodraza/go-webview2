# Changelog

### [1.0.1](https://www.github.com/mattpodraza/webview2/compare/v1.0.0...v1.0.1) (2021-05-03)


### CI

* fix the release name ([de625fe](https://www.github.com/mattpodraza/webview2/commit/de625fe9a2e653f8977adca691d24846a31f4962))


### Chores

* start showing CI commits in the changelog ([b41dd23](https://www.github.com/mattpodraza/webview2/commit/b41dd2375099b702c87f341c99915b2ff66daf5d))

## 1.0.0 (2021-05-03)


### ⚠ BREAKING CHANGES

* introduce better error handling, use the user32 helper
* embed DLLs using the `go:embed` directive
* change the import path, bump min. Go to 1.16

### Features

* add the `user32` package which wraps many of the syscall calls ([17e76a9](https://www.github.com/mattpodraza/webview2/commit/17e76a9678310a602f70b85ca28e65ab3ed9c883))


### Documentation

* add a minimal example ([4e9cb5d](https://www.github.com/mattpodraza/webview2/commit/4e9cb5d45ca7cbaf229ecc06a5a064da2979520a))
* update LICENSE ([5634e48](https://www.github.com/mattpodraza/webview2/commit/5634e48a4f8c3c55b907db07a7bf9b26c7999554))
* update README.md ([f8b3120](https://www.github.com/mattpodraza/webview2/commit/f8b3120cce0d497540289bfa754210c86f5c05a6))


### Chores

* ignore .vscode and built binaries ([104afbe](https://www.github.com/mattpodraza/webview2/commit/104afbe1b72716d4f7018b8edaea1a03bfcd3f0a))


### Refactoring

* change the import path, bump min. Go to 1.16 ([cbf6f57](https://www.github.com/mattpodraza/webview2/commit/cbf6f57c5e76d19147804fe6dd288e1cb2b79275))
* embed DLLs using the `go:embed` directive ([b287646](https://www.github.com/mattpodraza/webview2/commit/b287646acdcd485ef1c3d4068a48d4603e213868))
* introduce better error handling, use the user32 helper ([09698be](https://www.github.com/mattpodraza/webview2/commit/09698be696bc23cbb4389ffb58c787a502e41560))
