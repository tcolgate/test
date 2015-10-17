obj-m += bloomtab.o

#KPATH := /lib/modules/$(shell uname -r)/build
KPATH := /home/tristan/src/linux-4.2.2

#ccflags-y := $(KPATH)/include

all:
	make -C $(KPATH) M=$(PWD) modules

clean:
	make -C $(KPATH) /build M=$(PWD) clean
