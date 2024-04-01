.PHONY: build init config docs clean run stop stat front
default: clean init config docs front build

# set variables for build information.
buildTime=$(shell date -u "+%Y%m%d%H%M")
revVersion=$(shell git rev-parse --short HEAD)
tagVersion=$(shell git describe --tags)

# set variables for build
output=bin
app=template

# for build mode
mode=$(MODE)
ifeq ("$(mode)","")
mode=dev
endif

build: init
	# build $(app)
	go build -ldflags "-X stove-gitlab.sginfra.net/backend/template/constant.RevVersion=$(revVersion).$(buildTime) -X stove-gitlab.sginfra.net/backend/template/constant.TagVersion=$(tagVersion)" -o $(output)/$(app) main.go

init:
	# initialize output directory.
	@if [ -d $(output) ]; then rm -rf $(output); fi;
	@if [ ! -d $(output) ]; then mkdir -p $(output); mkdir -p $(output)/logs; mkdir -p $(output)/config; fi;

run: stop
	# run in background mode
	- cd $(output) && ./$(app) -c=config.yml &
stat:
	- @pgrep $(app)
stop:
	- @pkill -9 $(app)

test: run
	# test
	# If the test succeeds, kill the test process and then do the next build,
	# if it fails, kill the test process and exit the build.
	go test stove-gitlab.sginfra.net/backend/$(app)/controllers -timeout 30s \
		&& { echo test success; make stop; } \
		|| { echo test failure; make stop; exit 1;}

deploy: test
	# If the test is successful, deploy it to the development environment.

config: init
	cp etc/conf/config.idc.$(mode).yml bin/config.yml
	cp etc/conf/internal_api.yml bin/internal_api.yml
	cp etc/conf/external_api.yml bin/external_api.yml
	cp -R etc/conf/* bin/config

docs: init
	# copy the API docs to the output directory.
	@if [ ! -d $(output)/docs/int ]; then mkdir -p $(output)/docs/int; fi;
	@if [ ! -d $(output)/docs/ext ]; then mkdir -p $(output)/docs/ext; fi;
	@if [ ! -d $(output)/web ]; then mkdir -p $(output)/web; fi;
	cp -R etc/swagger/int/* bin/docs/int/
	cp -R etc/swagger/ext/* bin/docs/ext/
	cp -R etc/web/* bin/web/

clean:
	# clean output
	rm -rf $(output)