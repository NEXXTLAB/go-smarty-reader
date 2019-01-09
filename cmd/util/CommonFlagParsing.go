/*
 * This file is part of go-smarty-reader
 *
 * go-smarty-reader is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * go-smarty-reader is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with go-smarty-reader. If not, see <https://www.gnu.org/licenses/>.
 */

/*
   CommonFlagParsing holds the common command-line flag operations in one file to avoid code duplicates.
*/

package util

import (
	"flag"

	"github.com/golang/glog"
)

// Struct holding startup flag values
type Flag struct {
	Device *string
	Key    *string
	Mqtt   MqttInfo
}

func StartupFlagParsing() (flags Flag) {

	const VERSION = "1.1.0"

	// Define flags
	flags = Flag{
		Device: flag.String("device", "", "Serial device to read P1 data from."),
		Key:    flag.String("key", "", "Decryption Key to use."),
		Mqtt: MqttInfo{
			Broker: flag.String("mqttBroker", "ssl://iot.eclipse.org:8883",
				"MQTT Broker Address including protocol and port."),
			TopicRoot: flag.String("mqttTopicRoot", "nexxtlab/dev/smarty/go/",
				"MQTT Base topic, extended by extensions (such as OBIS codes) during publish."),
			Qos: flag.Int("mqttQos", 2,
				"MQTT Quality of service level."),
		},
	}

	flag.Parse()

	// Print version info and warnings if either the device- or keyFlag is missing
	glog.Infoln("Smarty Reader " + VERSION)
	if *flags.Device == "" {
		glog.Warningln("Serial device parameter missing.\n\t" +
			"This program instance will not be able to access any serial devices.")
	}
	if *flags.Key == "" {
		glog.Warningln("No decryption key found.\n\t" +
			"Telegrams can not be decrypted in this program instance.")
	}
	if *flags.Mqtt.Broker == "" {
		glog.Warningln("No MQTT Broker address found.\n\t" +
			"Unable to publish or subscribe to MQTT in this program instance.")
	}
	return flags
}
