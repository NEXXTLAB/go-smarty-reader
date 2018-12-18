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
    "os"

    "github.com/NEXXTLAB/go-smarty-reader/cmd/util"
    "github.com/NEXXTLAB/go-smarty-reader/share"
    "github.com/NEXXTLAB/go-smarty-reader/smarty"
    "github.com/basvdlei/gotsmart/dsmr"
    "github.com/eclipse/paho.mqtt.golang"
)

func main() {

    // Function defined in cmd/util/CommonFlagParsing.go
    deviceFlag, keyFlag := util.StartupFlagParsing()
    fmt.Printf("Device to read from: %s\n", *deviceFlag)
    fmt.Printf("Decryption key: %s\n", *keyFlag)

    // Get the machine hostname
    hostname := getHostname()

    // Preparing the MQTT connection
    client := mqttSetup(hostname)

    // Create a new smarty reader which will decrypt the telegrams after reading them
    // The serial connection is established right away
    // smartyObj is the object you may invoke methods on
    smartyObj := smarty.NewOnlineDecryptor(*deviceFlag, *keyFlag)

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

func getHostname() string {
    // Retrieve the machine hostname
    hostname, err := os.Hostname()
    if err != nil {
        // If getting a host name fails, let the user know and assign a fixed backup string
        fmt.Println(err)
        hostname = "BackupSystemHost-123"
    }
    fmt.Printf("Detected hostname: %s\n", hostname)
    return hostname
}

func mqttSetup(hostname string) share.Connection {
    // The topic root serves as a common root for all published messages.
    // In order to avoid interference of other users who might publish to the same topic
    // the hostname is part of the topic root.
    topicRoot := "nexxtlab/dev/smarty/go/" + hostname
    qualityOfService := 2

    // Setting some options, such as the broker to connect to
    // Further settings can be found on the "paho.mqtt.golang" project
    // Using the encrypted broker port
    // Alternatively you could connect to "iot.eclipse.org:1883" for an unencrypted connection
    opts := mqtt.NewClientOptions()
    opts.AddBroker("ssl://iot.eclipse.org:8883")

    // Create the client on which publishing operations can be executed
    return share.NewConnection(topicRoot, qualityOfService, opts)
}
