OUTPUT_FOLDER = public
STATIC_FOLDER = static

.PHONY: full build clean

full: clean build

build:
	mkdir -p $(OUTPUT_FOLDER)
	go run generate_static.go
	cp -r $(STATIC_FOLDER)/* $(OUTPUT_FOLDER)

clean:
	rm -rf $(OUTPUT_FOLDER)
