
include ./etc/.config

all:

bin:
	mkdir bin

bin/exif: bin src/exif.c
	gcc -lexif -o bin/exif src/exif.c

bin/exif-dumper: bin src/exif-dumper.go src/file2.go
	go build -o bin/exif-dumper src/exif-dumper.go src/file2.go

bin/sha512sum: bin src/sha512sum.go
	go build -o bin/sha512sum src/sha512sum.go

u:
	mkdir u

u/filesup: etc/.config etc/site.conf.in src/filesup.go src/file2.go u
	@sed \
		-e 's@%PREFIX%@$(shell pwd)@g' \
		-e 's@%FILESUP_UPLOADED_DIR%@$(FILESUP_UPLOADED_DIR)@g' \
		-e 's@%VIRTUAL_HOST_NAME%@$(VIRTUAL_HOST_NAME)@g' \
		etc/site.conf.in > etc/site.conf
	mkdir -p u var/log uploaded
	chmod 777 var/log uploaded
	go build -o u/filesup src/filesup.go src/file2.go

clean:
	rm -rf bin/exif-dumper bin/handler-photo-loader bin/sha512sum u sha512sum.go.sha512sum var etc/site.conf uploaded 

test:
	./bin/sha512sum src/sha512sum.go > sha512sum.go.sha512sum
	shasum -c sha512sum.go.sha512sum

