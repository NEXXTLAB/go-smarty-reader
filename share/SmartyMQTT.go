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
   SmartyMQTT serves as convenience, for wrapping the MQTT logic behind simple commands. It is not required and
   everything can be done using the third-party library. This wrapper offers simpler code and publish only if the
   value in question has not already been published, avoiding meaningless subscriber updates for specific OBIS
   codes
*/

package share

import (
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
)

var Client mqtt.Client

// Struct holding the MQTT client and user settings
type Connection struct {
	client   mqtt.Client
	settings Settings
}

// Struct holding user settings for the connection
type Settings struct {
	topicRoot string
	qos       int
	opts      *mqtt.ClientOptions
}

// Creating a new MQTT connection
// Parameter:
// * topicRoot: common prefix of a MQTT topic for this connection
// * qualityOfService: the MQTT quality of service for all operations
// * options:
// Return:
// * c: Connection struct, containing a connected MQTT client and the specified parameters
func NewConnection(topicRoot string, qualityOfService int, options *mqtt.ClientOptions) (c Connection) {
	lastCharacter := topicRoot[len(topicRoot)-1:]
	if lastCharacter != "/" {
		topicRoot = topicRoot + "/"
	}
	c = Connection{
		client: mqtt.NewClient(options),
		settings: Settings{
			topicRoot: topicRoot,
			qos:       qualityOfService,
			opts:      options,
		},
	}
	c.Reconnect()
	return c
}

// (Re)Connects the MQTT client
func (c Connection) Reconnect() {
	if !c.client.IsConnected() {
		token := c.client.Connect()
		if token.Wait() && token.Error() != nil {
			glog.Fatalln(token.Error())
		} else {
			glog.Infoln("MQTT Client connected")
		}
	}
}

// Publishes a message to the MQTT broker
// Parameter:
// * obis: the obis code, which will be used as topic suffix
// * value: the MQTT message containing the value
// * unit: if a unit is specified it will be appended to the MQTT message. To drop the unit use an empty string ""
// * retained: set to true if the message should be retained by the MQTT server
// * updateOnlyIfChanged: set to true to publish the message only if the value differs from the previous
//      This avoids short interval subscriber updates without any value change
// Return:
// * updated: true if a message was published (depending on updateOnlyIfChanged!)
func (c Connection) Publish(obis, value, unit string, retained, updateOnlyIfChanged bool) (updated bool) {
	formattedInput := formatValue(value)
	if unit != "" {
		formattedInput = formattedInput + " " + unit
	}
	// Implication of updateOnlyIfChanged => isNewValueFor
	if !updateOnlyIfChanged || isNewValueFor(obis, formattedInput) {
		tokenP := c.client.Publish(c.settings.topicRoot+obis,
			byte(c.settings.qos), retained, formattedInput)
		if tokenP.Wait() && tokenP.Error() == nil {
			updateValueFor(obis, formattedInput)
			glog.Infof("Successfully published %s for OBIS %s\n", formattedInput, obis)
		} else {
			glog.Errorf("Unable to publish %s for OBIS: %s\n", formattedInput, obis)
		}
		return tokenP.Error() == nil
	}
	return false
}

// Registers as subscriber to the specified topic
// Parameter:
// * obis: the topic extension to subscribe (topicRoot + extension)
// * callback: the callback function to call once a new message arrives (specify your own)
// Return:
// success: true is subscribing the topic was successful
func (c Connection) Subscribe(obis string, callback mqtt.MessageHandler) (success bool) {
	topic := c.settings.topicRoot + obis
	tokenP := c.client.Subscribe(topic, byte(c.settings.qos), callback)
	if tokenP.Wait() && tokenP.Error() == nil {
		glog.Errorf("Successfully subscribed to topic: %s\n", topic)
		return true
	} else {
		glog.Errorf("Unable to subscribe to topic: %s\n", topic)
		return false
	}
}

func formatValue(inputValue string) (formattedInput string) {
	formattedValue := strings.TrimLeft(inputValue, "0")
	if formattedValue == "" {
		formattedValue = "0"
	}
	if string(formattedValue[0]) == "." {
		formattedValue = "0" + formattedValue
	}
	return formattedValue
}

var lastReadValue = make(map[string]string)

func isNewValueFor(obis string, formattedInput string) bool {
	return lastReadValue[obis] != formattedInput
}

func updateValueFor(obis string, formattedInput string) {
	lastReadValue[obis] = formattedInput
}

// Disconnect from MQTT broker
// The does not invalidate the connection struct, Reconnect re-enables all functionality
// Parameter:
// * quiesce: amount of milliseconds to wait before closing
func (c Connection) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
	glog.Infoln("MQTT connection closed")
}
