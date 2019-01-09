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
	CommonMqttSetup holds the common MQTT operations in one file to avoid code duplicates.
*/

package util

import (
    "fmt"
    "math/rand"
    "os"

    "github.com/NEXXTLAB/go-smarty-reader/share"
    "github.com/eclipse/paho.mqtt.golang"
)

// Struct holding startup mqtt values
type MqttInfo struct {
    Broker    *string
    TopicRoot *string
    Qos       *int
}

func GetHostname() string {
    // Retrieve the machine hostname
    hostname, err := os.Hostname()
    if err != nil {
        // If getting a host name fails, let the user know and assign a random string
        fmt.Println(err)
        hostname = randomString(16)
    }
    fmt.Printf("Detected hostname: %s\n", hostname)
    return hostname
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    str := make([]byte, n)
    for i := range str {
        str[i] = letters[rand.Int63()%int64(len(letters))]
    }
    return string(str)
}

func MqttSetup(hostname string, info MqttInfo) share.MqttConnection {
    // The topic root serves as a common root for all published messages.
    // In order to avoid interference of other users who might publish to the same topic
    // the hostname is part of the topic root.
    topicRoot := *info.TopicRoot + hostname

    // Setting some options, such as the Broker to connect to.
    // Further settings can be found on the "paho.mqtt.golang" project.
    // Using the encrypted Broker "ssl://iot.eclipse.org:8883" as default using the startup flags.
    // Alternatively you could connect to "iot.eclipse.org:1883" for an unencrypted connection.
    opts := mqtt.NewClientOptions()
    opts.AddBroker(*info.Broker)

    // Create the client on which publishing operations can be executed
    return share.NewMqttConnection(topicRoot, *info.Qos, opts)
}
