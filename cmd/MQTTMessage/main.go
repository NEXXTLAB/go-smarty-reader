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
   Example on how to use the MQTT wrapper
*/

package main

import (
    "fmt"
    "os"

    "github.com/NEXXTLAB/go-smarty-reader/cmd/util"
    "github.com/NEXXTLAB/go-smarty-reader/share"
    "github.com/eclipse/paho.mqtt.golang"
)

func main() {

    util.StartupFlagParsing()

    // Get the machine hostname
    hostname := getHostname()

    // MQTT Setup extracted in a separate function.
    client := mqttSetup(hostname)

    // Publish "Hello" to the ""nexxtlab/dev/smarty/go/<hostname>/World" topic, without unit.
    client.Publish("World", "Hello", "", false, false)

    // Close the MQTT connection
    client.Disconnect(250)
}

func getHostname() string {
    // Retrieve the machine hostname
    hostname, err := os.Hostname()
    if err != nil {
        // If getting a host name fails, let the user know and assign a fixed backup string
        fmt.Println(err)
        hostname = "BackupSystemHost-123"
    }
    fmt.Printf("Deteced hostname: %s\n", hostname)
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
