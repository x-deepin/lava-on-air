BOARD_TAG=uno
#BOARD_TAG=pro328
CXXFLAGS+=-I../include
include /usr/share/arduino/Arduino.mk

deploy: all
	scp build-uno/firmware_slaver.hex 10.0.4.226:~/lava_slaver.hex
	ssh 10.0.4.226 ~/upload_lava_slaver

test:
	ssh 10.0.4.226 'cat /etc/wgetrc > /dev/ttyUSB1'
