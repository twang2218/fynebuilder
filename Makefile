.PHONY: fonts demo

fonts:
	mkdir -p ./theme/fonts
	wget -nc -P ./theme/fonts/ https://github.com/anthonyfok/fonts-wqy-microhei/raw/master/wqy-microhei.ttc
	xz ./theme/fonts/wqy-microhei.ttc

demo:
	cd cmd/demo && go run .

builder:
	cd cmd/builder && go run . -v
