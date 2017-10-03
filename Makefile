build:
	go build .

cf:
	cf uninstall-plugin MyBasicPlugin || true
	yes | cf install-plugin create-services-cliplugin

it: build cf
