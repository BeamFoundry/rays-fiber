# Changelog

## [0.4.1-rc.1](https://github.com/BeamFoundry/rays-fiber/compare/v0.4.0-rc.1...v0.4.1-rc.1) (2026-11-27)


### Bug Fixes

* **release:** fix release-please configuration for module tagging ([#25](https://github.com/BeamFoundry/rays-fiber/issues/25)) ([7949372](https://github.com/BeamFoundry/rays-fiber/commit/7949372cf51af622f333cab6978dc0d026fe5350))

## [0.4.0-rc.1](https://github.com/BeamFoundry/rays-fiber/compare/v0.3.2...v0.4.0-rc.1) (2026-11-25)


### Features

* **otel,controllers:** add contrib otel middleware, Optional provider injection, and TransactionManager context binding ([18e6ba1](https://github.com/BeamFoundry/rays-fiber/commit/18e6ba16d69ae335a52df0635ad42e8edd256459))
* **otel:** bind TransactionManager request context and support rays-otel providers ([cc7f52e](https://github.com/BeamFoundry/rays-fiber/commit/cc7f52e68568922f46d130856b074a8fd677010e))


### Bug Fixes

* **controllers,otel:** keep DefaultAppProvider route-free and verify Otel pre-tx ordering ([3adc656](https://github.com/BeamFoundry/rays-fiber/commit/3adc656d32d46eef675b87335614c9bc0ddfcfa5))
* **controllers,otel:** run controller-level pre-tx handlers before HttpSecurity ([d5a2245](https://github.com/BeamFoundry/rays-fiber/commit/d5a22453e2e90a7c54d6216c1962a9633604c7b4))


### Miscellaneous Chores

* **rays:** update "rays" to version v0.5.0-rc3 ([a7a2715](https://github.com/BeamFoundry/rays-fiber/commit/a7a2715a19640624243a6ae16b070eab1e3cbce3))
* **release:** prepare version v0.4.0-rc.1 ([fc368b5](https://github.com/BeamFoundry/rays-fiber/commit/fc368b553d7263b9576fcfa28585b2f1d6ba396f))
* **release:** prepare version v0.4.0-rc.1 ([0210b53](https://github.com/BeamFoundry/rays-fiber/commit/0210b53b4bcedf59a10bae47ecf51a6c036c0377))

## [0.3.2](https://github.com/BeamFoundry/rays-fiber/compare/v0.3.1...v0.3.2) (2026-11-20)


### Bug Fixes

* **controllers:** add QUERY HTTP method support and update default configuration ([c2bc62c](https://github.com/BeamFoundry/rays-fiber/commit/c2bc62ce58ba54c9dfcb34b0a8b960d245c5d47b))
* **security:** cast data to fiber.Map in RenderView to pass locals to views ([17f01b4](https://github.com/BeamFoundry/rays-fiber/commit/17f01b4c1562992effd9ff1367b602b74f7536f3))
