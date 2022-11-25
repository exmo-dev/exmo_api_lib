package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/tag"
)

const (
	delimiter = 0x01
)

type Presign struct {
	SendingTime  int64
	MsgSeqNum    int
	SenderCompID string
	TargetCompID string
	Password     string
}

func SignLogonMsg(msg *quickfix.Message, secret ApiSecret) {
	if msg.IsMsgTypeOf(enum.MsgType_LOGON) {
		//set passPhrase
		msg.Body.SetString(tag.Password, strconv.FormatInt(time.Now().UTC().Unix(), 10))

		//extract presign struct
		presign, err := getMsgPresign(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		//make preSignByte from presign
		preSignByte, err := makePresignByte(presign)
		if err != nil {
			fmt.Println(err)
			return
		}

		//sign is 64 based sha512(preSignByte)
		sign := createSignFromBodyAndSecret(preSignByte, []byte(secret))

		//sign the logonMessage
		msg.Body.SetString(tag.RawData, sign)
	}
}

func getMsgPresign(msg *quickfix.Message) (Presign, error) {
	var senderCompID quickfix.FIXString
	err := msg.Header.GetField(tag.SenderCompID, &senderCompID)
	if err != nil {
		return Presign{}, err
	}
	var targetCompID quickfix.FIXString
	err = msg.Header.GetField(tag.TargetCompID, &targetCompID)
	if err != nil {
		return Presign{}, err
	}
	var msgSeqNum quickfix.FIXInt
	err = msg.Header.GetField(tag.MsgSeqNum, &msgSeqNum)
	if err != nil {
		return Presign{}, err
	}
	var sendingTime quickfix.FIXUTCTimestamp
	err = msg.Header.GetField(tag.SendingTime, &sendingTime)
	if err != nil {
		return Presign{}, err
	}

	var password quickfix.FIXString
	err = msg.Body.GetField(tag.Password, &password)
	if err != nil {
		return Presign{}, err
	}

	return Presign{
		SendingTime:  sendingTime.UTC().Unix(),
		MsgSeqNum:    msgSeqNum.Int(),
		SenderCompID: senderCompID.String(),
		TargetCompID: targetCompID.String(),
		Password:     password.String(),
	}, nil
}

func makePresignByte(msg Presign) ([]byte, error) {
	presignByte := new(bytes.Buffer)

	//sendingTime
	binaryWriteErr := addToPresign(presignByte, msg.SendingTime, true)
	if binaryWriteErr != nil {
		return nil, binaryWriteErr
	}
	//msgSeqNum
	binaryWriteErr = addToPresign(presignByte, int64(msg.MsgSeqNum), true)
	if binaryWriteErr != nil {
		return nil, binaryWriteErr
	}
	//senderCompID
	binaryWriteErr = addToPresign(presignByte, msg.SenderCompID, true)
	if binaryWriteErr != nil {
		return nil, binaryWriteErr
	}
	//targetCompID
	binaryWriteErr = addToPresign(presignByte, msg.TargetCompID, true)
	if binaryWriteErr != nil {
		return nil, binaryWriteErr
	}
	//password
	binaryWriteErr = addToPresign(presignByte, msg.Password, false)
	if binaryWriteErr != nil {
		return nil, binaryWriteErr
	}

	return presignByte.Bytes(), nil
}

func addToPresign[Field int64 | string](presignByte *bytes.Buffer, field Field, withDelimeter bool) (binaryWriteErr error) {
	switch f := any(field).(type) {
	case string:
		_, binaryWriteErr = presignByte.WriteString(f)
		if binaryWriteErr != nil {
			return binaryWriteErr
		}
	default:
		binaryWriteErr = binary.Write(presignByte, binary.LittleEndian, f)
		if binaryWriteErr != nil {
			return binaryWriteErr
		}
	}
	if withDelimeter {
		return addDelimiter(presignByte)
	}
	return nil
}

func addDelimiter(presignByte *bytes.Buffer) error {
	binaryWriteErr := presignByte.WriteByte(delimiter)
	if binaryWriteErr != nil {
		return binaryWriteErr
	}
	return nil
}

func createSignFromBodyAndSecret(body, secret []byte) string {
	mac := hmac.New(sha512.New, secret)
	_, _ = mac.Write(body)

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
