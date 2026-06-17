# Changelog

## [1.4.0](https://github.com/prdktntwcklr/cloudmetrics/compare/v1.3.0...v1.4.0) (2026-06-17)


### Features

* **app:** add css styling ([#40](https://github.com/prdktntwcklr/cloudmetrics/issues/40)) ([1441693](https://github.com/prdktntwcklr/cloudmetrics/commit/1441693b42562055888ab3a6fe5a24bd487c4f5f))
* **app:** add POST endpoint for sensor readings ([#44](https://github.com/prdktntwcklr/cloudmetrics/issues/44)) ([bdeb477](https://github.com/prdktntwcklr/cloudmetrics/commit/bdeb477669d6cc033fb3b4fbe63850d6830b5218))

## [1.3.0](https://github.com/prdktntwcklr/cloudmetrics/compare/v1.2.0...v1.3.0) (2026-06-14)


### Features

* **app:** spin up several sensors ([#37](https://github.com/prdktntwcklr/cloudmetrics/issues/37)) ([21713d0](https://github.com/prdktntwcklr/cloudmetrics/commit/21713d0ffa566311185a5748b5c93ee29993e5e3))
* **app:** use structured logging ([#34](https://github.com/prdktntwcklr/cloudmetrics/issues/34)) ([df896c9](https://github.com/prdktntwcklr/cloudmetrics/commit/df896c97a88569ff2dc2052488ea8127ea257f31))

## [1.2.0](https://github.com/prdktntwcklr/cloudmetrics/compare/v1.1.0...v1.2.0) (2026-06-11)


### Features

* **ci:** add new version to bot commit message ([#33](https://github.com/prdktntwcklr/cloudmetrics/issues/33)) ([3e0c06b](https://github.com/prdktntwcklr/cloudmetrics/commit/3e0c06be3406f65a007fba8affe642de16e17f1e))
* **ci:** auto-update chart on push to main ([#29](https://github.com/prdktntwcklr/cloudmetrics/issues/29)) ([f13a509](https://github.com/prdktntwcklr/cloudmetrics/commit/f13a509613df464cfdb1fe7ddf085d84288d46c5))


### Bug Fixes

* **ci:** fix script path ([#31](https://github.com/prdktntwcklr/cloudmetrics/issues/31)) ([0a9f29c](https://github.com/prdktntwcklr/cloudmetrics/commit/0a9f29cfae9e845efca5f328e14a7d90ce26c49e))
* **ci:** use env var for image tag ([#32](https://github.com/prdktntwcklr/cloudmetrics/issues/32)) ([3f9e7ef](https://github.com/prdktntwcklr/cloudmetrics/commit/3f9e7ef8cfbffe4b6bf57fc59193b887358c1430))

## [1.1.0](https://github.com/prdktntwcklr/cloudmetrics/compare/v1.0.1...v1.1.0) (2026-06-09)


### Features

* add argo cd to automate deployment ([#22](https://github.com/prdktntwcklr/cloudmetrics/issues/22)) ([b9f28a4](https://github.com/prdktntwcklr/cloudmetrics/commit/b9f28a4257dfc0d5b245e7fe823d9a93210c5f81))
* **app:** add git sha version info ([#26](https://github.com/prdktntwcklr/cloudmetrics/issues/26)) ([2d74dfd](https://github.com/prdktntwcklr/cloudmetrics/commit/2d74dfdbfb5451ec54eecd9b36da020f97b4c22f))
* **app:** add health check ([#25](https://github.com/prdktntwcklr/cloudmetrics/issues/25)) ([2c1273a](https://github.com/prdktntwcklr/cloudmetrics/commit/2c1273ad196ffe2491e412a1f49861a6134c77d5))

## [1.0.1](https://github.com/prdktntwcklr/cloudmetrics/compare/v1.0.0...v1.0.1) (2026-06-03)


### Bug Fixes

* **ci:** push image to ghcr ([#20](https://github.com/prdktntwcklr/cloudmetrics/issues/20)) ([f17fe30](https://github.com/prdktntwcklr/cloudmetrics/commit/f17fe3072f35df7d0417469bd9b378f59e4dfd18))

## 1.0.0 (2026-06-03)


### Features

* **app:** add better logging ([#10](https://github.com/prdktntwcklr/cloudmetrics/issues/10)) ([a231df1](https://github.com/prdktntwcklr/cloudmetrics/commit/a231df11e366b817366fff3ee456f432279baced))
* **app:** add landing page ([#8](https://github.com/prdktntwcklr/cloudmetrics/issues/8)) ([74a4a75](https://github.com/prdktntwcklr/cloudmetrics/commit/74a4a75af79a4e07dd567245385dbf56cb2b10ec))
* **ci:** add release-please ([#18](https://github.com/prdktntwcklr/cloudmetrics/issues/18)) ([47f894a](https://github.com/prdktntwcklr/cloudmetrics/commit/47f894a1499483c5beef6605c5f71f45d61b80d3))
* **deploy:** use Helmfile for observability stack ([#17](https://github.com/prdktntwcklr/cloudmetrics/issues/17)) ([63c19ae](https://github.com/prdktntwcklr/cloudmetrics/commit/63c19ae874e593e40d48e066520a84d8cd85f68b))
* initial commit ([1beaca5](https://github.com/prdktntwcklr/cloudmetrics/commit/1beaca538f3bdc30b2c29319947f1858dcde8bb6))
* **k8s:** add deployment and service manifests ([#5](https://github.com/prdktntwcklr/cloudmetrics/issues/5)) ([895048b](https://github.com/prdktntwcklr/cloudmetrics/commit/895048b5a395d95404fd4fe7889bcf17f1bac0b7))
* **k8s:** add grafana dashboard ([#11](https://github.com/prdktntwcklr/cloudmetrics/issues/11)) ([27f2eae](https://github.com/prdktntwcklr/cloudmetrics/commit/27f2eae3b5915d18c548144bf1b00e04657e0daf))
* **k8s:** add Loki and Alloy for logs ([#16](https://github.com/prdktntwcklr/cloudmetrics/issues/16)) ([52b3992](https://github.com/prdktntwcklr/cloudmetrics/commit/52b3992325cf27f05ee1a6773457a5e0282358d3))
* **k8s:** add Prometheus monitoring ([#6](https://github.com/prdktntwcklr/cloudmetrics/issues/6)) ([ae8781d](https://github.com/prdktntwcklr/cloudmetrics/commit/ae8781da793ccd17c4b043b25c3e7de213b5da06))
* **k8s:** create Helm chart ([#15](https://github.com/prdktntwcklr/cloudmetrics/issues/15)) ([72c1fd2](https://github.com/prdktntwcklr/cloudmetrics/commit/72c1fd2ce7214697b07775bb597dcf49028ed477))
* **k8s:** use table in grafana ([#13](https://github.com/prdktntwcklr/cloudmetrics/issues/13)) ([4633ece](https://github.com/prdktntwcklr/cloudmetrics/commit/4633eceec15c0e69d471c3f92c669c1c9aad2204))
* **monitoring:** use k8s secret for grafana login ([#14](https://github.com/prdktntwcklr/cloudmetrics/issues/14)) ([5200390](https://github.com/prdktntwcklr/cloudmetrics/commit/5200390dda9b040e6dac359229a6cba4b6b1c852))


### Bug Fixes

* **k8s:** put ServiceMonitor in default namespace ([#7](https://github.com/prdktntwcklr/cloudmetrics/issues/7)) ([b35c53a](https://github.com/prdktntwcklr/cloudmetrics/commit/b35c53a6111fe7ce9334c231bcd7ef172c07108f))
* **k8s:** set pull policy to always ([#9](https://github.com/prdktntwcklr/cloudmetrics/issues/9)) ([e02d910](https://github.com/prdktntwcklr/cloudmetrics/commit/e02d9105e3e28dabec77b02c8b125eec4cfedc48))
