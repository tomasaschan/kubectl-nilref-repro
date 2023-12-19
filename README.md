# kubectl-nilref-repro

This repository exists as a reproduction of a nilref exception in k8s.io/kubectl@v0.29.0.

Some context can be found here: https://github.com/kubernetes/kubernetes/commit/eb32969ab8f3d1455c33a595bba5b3ce4d8a03b8#r135086674

In order to run the repro locally, clone this repo and then run `make generate manifests test`.

The failing test tries to ensure that a resources is successfully created by the `kubebuilder-declarative-pattern`-based controller, but the resource is never created because the controller runs into the nilref panic. The panic, including a full stack trace, is available in the test logs.
