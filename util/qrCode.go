package util

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	qrcodeGen "github.com/skip2/go-qrcode"
)

// QrCodeDecode 解码二维码
func QrCodeDecode(pngBytes []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		return "", err
	}

	// 准备 BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	// 解码图像
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// QrCodePrint 打印二维码
func QrCodePrint(content string) error {
	// 生成二维码并输出到命令行
	qr, err := qrcodeGen.New(content, qrcodeGen.Low)
	if err != nil {
		return err
	}
	fmt.Println(qr.ToSmallString(false))

	return nil
}

// ClearQrCode 清空二维码
func ClearQrCode() {
	// 清空26行
	for i := 0; i < 26; i++ {
		fmt.Print("\033[1A\033[K")
	}
}
