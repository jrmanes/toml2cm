# Toml2CM

---

## Description
Toml 2 ConfigMap is a tool that help us to migrate Toml files to Kubernetes
ConfigMaps.

This tools was created with the goal of migrating `toml` configurations to
ConfigMaps in an easy way and removing the toil.

---

## How does it work?

It has some flags (currently one), where we can specify the file to the path
that we would like to convert.
It reads the `toml` file and parse it line by line, generating a new one with
a ConfigMap format.

---

## TODO

- [ ] Add ConfigMap header
- [ ] Add tests 👀

---

Jose Ramon Mañes
