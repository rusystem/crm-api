package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/skip2/go-qrcode"
	"image/png"
)

// GenerateQRCodePNG генерация QR-кода из структуры QRInfo
func GenerateQRCodePNG(info domain.CodeInfo) ([]byte, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CodeInfo to JSON: %v", err)
	}

	qrCode, err := qrcode.Encode(string(data), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("error generating QR code: %v", err)
	}

	return qrCode, nil
}

// GenerateBarcode генерация штрих-кода из строки
func GenerateBarcode(info domain.CodeInfo, width, height int) ([]byte, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CodeInfo to JSON: %v", err)
	}

	fmt.Println(string(data))
	fmt.Println(data)

	barCode, err := code128.EncodeWithoutChecksum(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to generate barcode: %v", err)
	}

	barCode, err = barcode.Scale(barCode, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %v", err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, barCode)
	if err != nil {
		return nil, fmt.Errorf("failed to encode barcode to png: %v", err)
	}

	return buf.Bytes(), nil
}
