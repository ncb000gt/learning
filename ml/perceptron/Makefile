all: clean install

deppkg:
	go get github.com/droundy/goopt

install: deppkg
	cd src/digitaltumbleweed.com/perceptron
	go install
	cd ../../../

clean:
	rm -Rf pkg

profile:
	/usr/local/go/pkg/tool/linux_amd64/6g -vx runner.go

.PHONY : all deppkg install clean test profile
