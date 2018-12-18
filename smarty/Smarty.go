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
   This file handles reading the smarty telegrams over a serial connection, and splitting in its different tokens
   using a state machine.
*/
package smarty

import (
    "bufio"

    "github.com/golang/glog"
    "github.com/tarm/serial"
)

const GCMTagLength = 12
type State int
const (
    waitingForStartByte State = iota + 1
    readSystemTitleLength
    readSystemTitle
    readSeparator82
    readPayloadLength
    readSeparator30
    readFrameCounter
    readPayload
    readGcmTag
    doneReadingTelegram
)
var (
    state                                                = waitingForStartByte
    currentBytePosition, changeToNextStateAt, dataLength int
    systemTitle, frameCounter, dataPayload, gcmTag       []byte
)

type Smarty interface {
    Disconnect()
}

type deviceInfo struct {
    deviceName string
    reader     *bufio.Reader
    port       *serial.Port
}

func processByteStream(input []byte, length int) (ready bool) {
    ready = false
    for i := 0; i < length && !ready; i++ {
        // Keep track of the position in the byte stream
        currentBytePosition++

        // Run the appropriate actions (nil if telegram not yet complete)
        ready = processStateActions(input[i])
    }
    return

}

func resetVariables() {
    state = waitingForStartByte
    currentBytePosition = 0
    changeToNextStateAt = 0
    systemTitle = []byte("")
    dataLength = 0
    frameCounter = []byte("")
    dataPayload = []byte("")
    gcmTag = []byte("")
}

func processStateActions(rawInput byte) (ready bool) {
    switch state {
    case waitingForStartByte:
        if rawInput == 0xDB {
            resetVariables()
            state = readSystemTitleLength
        }
    case readSystemTitleLength:
        state = readSystemTitle
        // 2 start bytes (position 0 and 1) + system title length
        changeToNextStateAt = 1 + int(rawInput)
    case readSystemTitle:
        systemTitle = append(systemTitle, rawInput)
        if currentBytePosition >= changeToNextStateAt {
            state = readSeparator82
            changeToNextStateAt++
        }
    case readSeparator82:
        if rawInput == 0x82 {
            state = readPayloadLength // Ignore separator byte
            changeToNextStateAt += 2
        } else {
            glog.Errorln("Missing separator (0x82). Dropping telegram.")
            state = waitingForStartByte
        }
    case readPayloadLength:
        dataLength <<= 8
        dataLength |= int(rawInput)
        if currentBytePosition >= changeToNextStateAt {
            state = readSeparator30
            changeToNextStateAt++
        }
    case readSeparator30:
        if rawInput == 0x30 {
            state = readFrameCounter
            // 4 bytes for frame counter
            changeToNextStateAt += 4
        } else {
            glog.Errorln("Missing separator (0x30). Dropping telegram.")
            state = waitingForStartByte
        }
    case readFrameCounter:
        frameCounter = append(frameCounter, rawInput)
        if currentBytePosition >= changeToNextStateAt {
            state = readPayload
            changeToNextStateAt += dataLength - 17
        }
    case readPayload:
        dataPayload = append(dataPayload, rawInput)
        if currentBytePosition >= changeToNextStateAt {
            state = readGcmTag
            changeToNextStateAt += GCMTagLength
        }
    case readGcmTag:
        // All input has been read.
        gcmTag = append(gcmTag, rawInput)
        if currentBytePosition >= changeToNextStateAt {
            state = doneReadingTelegram
        }
    }
    if state == doneReadingTelegram {
        state = waitingForStartByte
        return true
    }
    return false
}

func prepareCipherComponents() (iv, cipherText []byte) {
    iv = append(systemTitle, frameCounter...)
    cipherText = append(dataPayload, gcmTag...)
    return
}

func openSerialConnection(deviceName string) (reader *bufio.Reader, port *serial.Port) {
    config := &serial.Config{
        Name:     deviceName,
        Baud:     115200,
        Size:     8,
        Parity:   serial.ParityNone,
        StopBits: serial.StopBits(1),
    }
    port, err := serial.OpenPort(config)
    if err != nil {
        glog.Fatalln(err.Error())
        return nil, nil
    }
    return bufio.NewReader(port), port
}

func readTelegram(reader *bufio.Reader) {
    var ready = false
    buffer := make([]byte, 4096)
    for !ready {
        if length, err := reader.Read(buffer); err == nil && length > 0 {
            ready = processByteStream(buffer, length)
        }
    }
}

func ProcessTelegram(input []byte) (iv, cipherText []byte) {
    ok := processByteStream(input, len(input))
    if !ok {
        glog.Errorf("Telegram tokenization unable to complete.")
    }
    return prepareCipherComponents()
}
