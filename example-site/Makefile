#!/usr/bin/make -f

dev:
	echo "Navigate a browser to http://localhost:7070/" && \
		huguinho-dev -content "./content" -templates "./templates" -target "./rendered" -with-drafts -with-future

generate:
	huguinho -content "./content" -templates "./templates" -target "./rendered"

clean:
	rm -rf "./rendered" && mkdir "./rendered"
