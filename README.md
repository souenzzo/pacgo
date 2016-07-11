# pacgo
Manage AUR packages via CLI or HTTP

## I's a draft/prototype. Do not use without read the code!

First, add this to `pacman.conf`
```
[pacgo]
SigLevel = Optional TrustAll
Server = file:///home/fulano/.local/pkg
```

Add a package src

` pacgo -add=https://aur.archlinux.org/package.git `

List avalible src's

` pacgo -list `

Build a src

` pacgo -make package `

Install the package

` pacman -Sy package `

