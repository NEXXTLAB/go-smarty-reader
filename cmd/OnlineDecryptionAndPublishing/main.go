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
		Example on how to use the OnlineDecryptor with publishing the different OBIS codes over MQTT using a secure
		channel. Remember even though the transmission is secure, the iot.eclipse MQTT broker is not!
        It is a public broker, everyone is able to read your published messages.
        If you wish to use this example as basis for your own implementation, you should switch to a trusted MQTT broker.
*/

package main

import (
	"fmt"

	"github.com/NEXXTLAB/go-smarty-reader/cmd/util"
	"github.com/NEXXTLAB/go-smarty-reader/smarty"
	"github.com/basvdlei/gotsmart/dsmr"
)

func main() {

	// Function defined in cmd/util/CommonFlagParsing.go
	flags := util.StartupFlagParsing()
	fmt.Printf("Device to read from: %s\n", *flags.Device)
	fmt.Printf("Decryption key: %s\n", *flags.Key)

	// Preparing the MQTT connection
	// Functions defined in cmd/util/CommonMqttSetup.go
	client := util.MqttSetup(util.GetHostname(), flags.Mqtt)

	// Create a new smarty reader which will decrypt the telegrams after reading them
	// The serial connection is established right away
	// smartyObj is the object you may invoke methods on
	smartyObj := smarty.NewOnlineDecryptor(*flags.Device, *flags.Key)

	// Read until 100 telegrams could be successfully decrypted and published
	for telegramCounter := 0; telegramCounter < 100; {
		// Wait, get and decrypt the next telegram
		plainText, ok := smartyObj.GetTelegram()
		// If the decryption was successful, print the payload to the console
		if ok {
			frame, err2 := dsmr.ParseFrame(string(plainText))
			if err2 == nil {
				// Publish all present OBIS codes.
				// Take note that timestamp, equipmentID, header and version are stored separately in this package.
				// frame.Objects includes only measured data.
				for _, token := range frame.Objects {
					client.Publish(token.ID, token.Value, token.Unit, false, true)
				}
				telegramCounter++
			} else {
				fmt.Println(err2)
			}
		}
	}

	// Close the MQTT connection
	client.Disconnect(250)

	// After use, remember to close to serial port!
	smartyObj.Disconnect()
}
