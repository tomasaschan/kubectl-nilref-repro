# kubectl-nilref-repro

This repository exists as a reproduction of a nilref exception in k8s.io/kubectl@v0.29.0.

Some context can be found here: https://github.com/kubernetes/kubernetes/commit/eb32969ab8f3d1455c33a595bba5b3ce4d8a03b8#r135086674

In order to run the repro locally, clone this repo and then run `make generate manifests test`.

The failing test tries to ensure that a resources is successfully created by the `kubebuilder-declarative-pattern`-based controller, but the resource is never created because the controller runs into the nilref panic. The panic, including a full stack trace, is available in the test logs.

Example output:

```shell
λ make test
test -s /home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/controller-gen && /home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/controller-gen --version | grep -q v0.13.0 || \
GOBIN=/home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
test -s /home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/setup-envtest || GOBIN=/home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
KUBEBUILDER_ASSETS="/home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/bin/k8s/1.28.0-linux-amd64" go test ./... -coverprofile cover.out
?       github.com/tomasaschan/kubectl-nilref-repro/api/v1alpha1        [no test files]
?       github.com/tomasaschan/kubectl-nilref-repro/cmd [no test files]
Running Suite: Controller Suite - /home/tomasl/src/github.com/tomasaschan/kubectl-nilref-repro/internal/controller
==================================================================================================================
Random Seed: 1702980420

Will run 1 of 1 specs
configmap/repro created
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
        panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x265b976]

goroutine 190 [running]:
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile.func1()
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:116 +0x1e5
panic({0x2c32f60?, 0x45eb500?})
        /usr/lib/go/src/runtime/panic.go:914 +0x21f
k8s.io/kubectl/pkg/cmd/apply.(*Patcher).patchSimple(0xc000544fb0, {0x346c630?, 0xc00068a5c8}, {0xc000232a80, 0x37a, 0x380}, {0xc00063b9a0, 0x7}, {0xc00063b980, 0x5}, ...)
        /home/tomasl/src/pkg/mod/k8s.io/kubectl@v0.29.0/pkg/cmd/apply/patcher.go:160 +0x6d6
k8s.io/kubectl/pkg/cmd/apply.(*Patcher).Patch(0xc000544fb0, {0x346c630, 0xc00068a5c8}, {0xc000232a80, 0x37a, 0x380}, {0x2f792a1, 0xe}, {0xc00063b9a0, 0x7}, ...)
        /home/tomasl/src/pkg/mod/k8s.io/kubectl@v0.29.0/pkg/cmd/apply/patcher.go:360 +0xa5
k8s.io/kubectl/pkg/cmd/apply.(*ApplyOptions).applyOneObject(0xc00018d6c0, 0xc000636380)
        /home/tomasl/src/pkg/mod/k8s.io/kubectl@v0.29.0/pkg/cmd/apply/apply.go:731 +0xf3d
k8s.io/kubectl/pkg/cmd/apply.(*ApplyOptions).Run(0xc00018d6c0)
        /home/tomasl/src/pkg/mod/k8s.io/kubectl@v0.29.0/pkg/cmd/apply/apply.go:532 +0x2b6
sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/applier.(*directApplier).Run(0xc000636380?, 0x0?)
        /home/tomasl/src/pkg/mod/sigs.k8s.io/kubebuilder-declarative-pattern@v0.15.0-beta.1.0.20231025132032-bdee2d107c42/pkg/patterns/declarative/pkg/applier/direct.go:45 +0x16
sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/applier.(*DirectApplier).Apply(0x45ecd10, {0xc000600c00?, 0x2f7f26b?}, {{0xc00068a508, 0x1, 0x1}, 0xc0006b2b40, {0x348c5c0, 0xc00063fa00}, {0xc0006d41d0, ...}, ...})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/kubebuilder-declarative-pattern@v0.15.0-beta.1.0.20231025132032-bdee2d107c42/pkg/patterns/declarative/pkg/applier/direct.go:183 +0x1025
sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative.(*Reconciler).reconcileExists(0xc00057cdb0, {0x3485018, 0xc000693b60}, {{0xc0006d41d0?, 0xc0006ba780?}, {0xc0006d41c6?, 0xc0006d41d0?}}, {0x34a0f50?, 0xc0006ba780?})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/kubebuilder-declarative-pattern@v0.15.0-beta.1.0.20231025132032-bdee2d107c42/pkg/patterns/declarative/reconciler.go:406 +0x1908
sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative.(*Reconciler).Reconcile(0xc00057cdb0, {0x3485018?, 0xc000693b60}, {{{0xc0006d41d0?, 0x0?}, {0xc0006d41c6?, 0x5?}}})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/kubebuilder-declarative-pattern@v0.15.0-beta.1.0.20231025132032-bdee2d107c42/pkg/patterns/declarative/reconciler.go:183 +0x56b
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile(0x348a6d0?, {0x3485018?, 0xc000693b60?}, {{{0xc0006d41d0?, 0xb?}, {0xc0006d41c6?, 0x0?}}})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:119 +0xb7
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler(0xc00041f9a0, {0x3485050, 0xc00002fa40}, {0x2d27aa0?, 0xc0001fe1a0?})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:316 +0x3cc
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem(0xc00041f9a0, {0x3485050, 0xc00002fa40})
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:266 +0x1c9
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2()
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:227 +0x79
created by sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2 in goroutine 168
        /home/tomasl/src/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:223 +0x565
FAIL    github.com/tomasaschan/kubectl-nilref-repro/internal/controller 4.732s
FAIL
make: *** [Makefile:65: test] Error 1
```
