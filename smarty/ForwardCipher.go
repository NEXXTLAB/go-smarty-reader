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
   The CipherForwarder reads a smarty telegram, but instead of decrypting it, it returns the necessary parts
   for decryption at a later date, for instance if you transmit it over insecure channels.
*/

package smarty

import (
	"github.com/golang/glog"
)

// Struct allowing to retrieve smarty telegrams, split into initial value, cipher text and gcm tag
type CipherForwarder struct {
	deviceInfo
}

// Creation of a new CipherForwarder
// Parameter:
// * deviceName: the port to listen to
// Return:
// * CipherForwarder: a new object to execute methods on
func NewCipherForwarder(deviceName string) CipherForwarder {
	reader, port := openSerialConnection(deviceName)
	glog.Infoln("Serial connection established")
	return CipherForwarder{
		deviceInfo: deviceInfo{
			deviceName: deviceName,
			reader:     reader,
			port:       port,
		},
	}
}

// Waits for the next telegram and splits it into its tokens
// If you are to decrypt this result with the Decryptor, you will have to append the gcm tag to the cipher text!
// Return:
// * initialValue: the initial value as specified in the smarty documentation
// * cipherText: the payload
// * gcmTag: the aes-gcm tag
func (cf *CipherForwarder) GetTelegram() (initialValue, cipherText, gcmTag []byte) {
	readTelegram(cf.deviceInfo.reader)
	return cf.forwardTelegram()
}

func (cf *CipherForwarder) getDeviceInfo() deviceInfo {
	return cf.deviceInfo
}

func (cf *CipherForwarder) forwardTelegram() (iv, cipherText, gcmTag []byte) {
	iv, cipherText = prepareCipherComponents()
	return iv, cipherText[:len(cipherText)-GCMTagLength], cipherText[len(cipherText)-GCMTagLength:]
}

// Disconnect the serial connection
func (cf *CipherForwarder) Disconnect() {
	err := cf.port.Close()
	if err == nil {
		glog.Infoln("Serial connection closed")
	} else {
		glog.Errorln("Unable to close serial connection")
	}
}
