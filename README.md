# Go Smarty Reader

The smart meter (called _Smarty_) deployed in Luxembourg allows the user to access a predefined set of data over the P1 port, a serial read-only interface, like its counterpart in the Netherlands.
However the data stream in Luxembourg is encrypted using the AES128-GCM algorithm (see the [specification](https://www.nexxtlab.lu/download/453/)), which raises the barrier for anyone interested in their own energy consumption.
This is where Go Smarty Reader comes into play: handling the serial connection, decrypting the Smarty telegrams, optionally publishing the measurements via MQTT, and with Golang as basis it easily compiles for different platforms.

For additional information please visit our [blog post](https://www.nexxtlab.lu/smarty-dongle/) at NEXXTLAB
  

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. If the explanations here are insufficient to your needs, please take a look at our [blog post](https://www.nexxtlab.lu/smarty-dongle/), where you will find additional details.  

### Prerequisites

Have [Go installed](https://golang.org/doc/install) and properly set up.

In order to connect Smarty to your machine you will need a P1 Cable (Dutch name: _Slimme Meter Kabel P1_), or build it yourself as seen on [weigu.lu](http://weigu.lu/microcontroller/smartyreader/index.html).
Furthermore ask your electrical grid operator for your decryption key.


### Installing

Run the following command to get a local copy of the project
```
go get github.com/NEXXTLAB/go-smarty-reader
```
After this you need to get the project dependencies. Either you 'go get' all four of the third party libraries listed at the bottom, or you use [dep](https://github.com/golang/dep) with your console points to the project directory.
```
dep ensure
```

### Running the examples

Please find all prepared examples inside of [cmd/](https://github.com/NEXXTLAB/go-smarty-reader/tree/master/cmd). In order to run any example you should take note of following command-line arguments:
* key: your device specific decryption key
* device: the interface on which the smart meter is connected to your target platform
* stderrthreshold: allows to adjust the console output of glog, see [here](https://svn.apache.org/repos/asf/incubator/mesos/trunk/third_party/glog-0.3.1/doc/glog.html?p=1197837) for more details

The key should be, as mentioned earlier, requested from your electricity grid operator. To set the device argument, you will need to find the correct interface. Here a quick How-To:
* Windows: open your *Device Manager*, expand the *Ports* section, find the correct device and write down the COM port (eg. COM8)
* Linux: open your terminal and run *dmesg*. Plug in your P1 cable and write down the device name (eg. /dev/ttyUSB2)

Navigate to the project main directory, then run:
```
go run ./cmd/OnlineDecryption/main.go -key yourKey -device yourInterface -stderrthreshold=INFO
```
Now you should see every 10 seconds the result of a decrypted Smarty telegram in your console. Please find the meaning of the OBIS codes in the [specification](https://www.nexxtlab.lu/download/453/)  

You may swap the *OnlineDecryption* part of the path to any other example found in the [cmd/](https://github.com/NEXXTLAB/go-smarty-reader/tree/master/cmd) folder. Not every example requires all arguments.


## Running the tests

Currently there is only one test, additions welcome!
```
go test github.com/NEXXTLAB/go-smarty-reader/smarty/
```

## Build

In order to build the project simply run the [go build](https://golang.org/cmd/go/) command inside the project directory.
[The Go Cookbook](https://golangcookbook.com/chapters/running/cross-compiling/) offers a short overview on how to cross compile for a target platform different from yours. 
After you get your platform specific binary, you may run it in your console using the arguments found in the [Running the examples](#running-the-examples) section.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Sébastien Thill** - *Initial work* - [Netskeh](https://github.com/Netskeh)

See also the list of [contributors](https://github.com/NEXXTLAB/go-smarty-reader/contributors) who participated in this project.

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

### Third Party Libraries
* [Eclipse Paho MQTT Go client](https://github.com/eclipse/paho.mqtt.golang)
* [Google glog](https://github.com/golang/glog)
* [GotSmart](https://github.com/basvdlei/gotsmart)
* [Serial](https://github.com/tarm/serial)

### Other Smarty Projects
* [Smarty DSMR Proxy in Python](https://github.com/mweimerskirch/smarty_dsmr_proxy)
* [SmartyReader for Arduino](http://weigu.lu/microcontroller/smartyreader/index.html)

### Additional Resources
* [README-Template.md](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)

