BOARD_TAG=uno
CXXFLAGS+=-I../include
include /usr/share/arduino/Arduino.mk

deploy: all
	scp build-uno/firmware_slaver.hex 10.0.4.226:/dev/shm/lava_slaver.hex
	ssh 10.0.4.226 ~/upload_lava_slaver

test:
	ssh 10.0.4.226 'echo "Testing" > /dev/ttyUSB1'
