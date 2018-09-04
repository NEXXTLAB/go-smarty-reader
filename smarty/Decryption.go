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
   This file contains two Object definitions. First the Decryptor which is able to decrypt a provided smarty
   telegram and return the original text. Second is the OnlineDecryptor which will listen to the serial device,
   read, split it into its components and return the plain text.
*/

package smarty

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/golang/glog"
)

// Struct allowing simple decryption of existing telegrams
type Decryptor struct {
	key, aad []byte
}

// Struct allowing live capture of smarty telegrams with decryption
type OnlineDecryptor struct {
	decryptor Decryptor
	deviceInfo
}

// Creation of a new OnlineDecryptor
// Parameter:
// * deviceName: the port to listen to
// * decryptionKey: your smarty key
// Return:
// * OnlineDecryptor: a new object to execute methods on
func NewOnlineDecryptor(deviceName, decryptionKey string) OnlineDecryptor {
	reader, port := openSerialConnection(deviceName)
	glog.Infoln("Serial connection established")
	return OnlineDecryptor{
		decryptor: NewDecryptor(decryptionKey),
		deviceInfo: deviceInfo{
			reader:     reader,
			port:       port,
			deviceName: deviceName,
		},
	}
}

// Creation of a new Decryptor
// Parameter:
// * decryptionKey: your smarty key
// Return:
// * Decryptor: a new object to execute methods on
func NewDecryptor(decryptionKey string) Decryptor {
	decodedKey, err := hex.DecodeString(decryptionKey)
	if err != nil {
		glog.Errorf("Error parsing the decryption key: %s", err.Error())
	}
	if len(decryptionKey) != 32 {
		glog.Errorf("Invalid decryption key length\n"+
			"Required 32 characters\n"+
			"Found %v characters", len(decryptionKey))
	}
	decodedAad, _ := hex.DecodeString("3000112233445566778899AABBCCDDEEFF")
	return Decryptor{
		key: decodedKey,
		aad: decodedAad,
	}
}

// Waits for the next telegram and decrypts it
// Return:
// * plaintText: the decrypted text
// * ok: true if the decryption was successful
func (od *OnlineDecryptor) GetTelegram() (plainText []byte, ok bool) {
	readTelegram(od.deviceInfo.reader)
	return od.decryptor.Decrypt(prepareCipherComponents())
}

func (od *OnlineDecryptor) getDeviceInfo() deviceInfo {
	return od.deviceInfo
}

// Decrypt a smarty telegram
// Parameter:
// * initialValue: the initial value as specified in the smarty documentation
// * cipherText: this expects the payload with appended gcm tag at the end (payload + gcmTag)!
// Return:
// * plaintText: the decrypted text
// * ok: true if the decryption was successful
func (d Decryptor) Decrypt(initialValue, cipherText []byte) (plainText []byte, ok bool) {
	cipherBlock, err := aes.NewCipher(d.key)
	if err != nil {
		glog.Fatalln(err.Error())
	}

	aesgcm, err := cipher.NewGCMWithTagSize(cipherBlock, 12)
	if err != nil {
		glog.Fatalln(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, initialValue, cipherText, d.aad)
	if err != nil {
		glog.Errorln(err.Error())
	}

	return plaintext, err == nil
}

// Disconnect the serial connection
func (od *OnlineDecryptor) Disconnect() {
	err := od.port.Close()
	if err == nil {
		glog.Infoln("Serial connection closed")
	} else {
		glog.Errorln("Unable to close serial connection")
	}

}
