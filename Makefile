build:
	go build .

cf:
	cf uninstall-plugin CreateServicePush || true
	cf install-plugin create-services-cliplugin

it: build cf
