# Toml2CM

---

## Description
Toml 2 ConfigMap is a tool that help us to migrate Toml files to [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) in a 
[Helm template](https://helm.sh/docs/chart_best_practices/templates/) way.

This tools was created with the goal of migrating `toml` configurations to
ConfigMaps in an easy way and removing the toil.

---

## How does it work?

It has some flags (currently one), where we can specify the file to the path
that we would like to convert.
It reads the `toml` file and parse it line by line, generating a new one with
a ConfigMap format.

---

## Sources

- [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)
- [Helm template](https://helm.sh/docs/chart_best_practices/templates/)

---

## Requirements

- Go: Install go


## How to execute it?

```go
go run *.go --file=<fileName>
// Example
go run *.go --file=example.toml
```

---

## TODO

- [x] Add ConfigMap header
- [ ] Add tests 👀
- [ ] Read all files inside a folder

---

Jose Ramon Mañes
