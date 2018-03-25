DEP ?= dep
BUILD_DIR = ./build
NATURE_REMO_PLUGIN = $(BUILD_DIR)/mackerel-plugin-nature-remo

all: deps $(NATURE_REMO_PLUGIN)

deps:
	$(DEP) ensure -vendor-only

$(NATURE_REMO_PLUGIN):
	go build -o $(NATURE_REMO_PLUGIN) main.go

clean:
	rm -rf $(BUILD_DIR)
