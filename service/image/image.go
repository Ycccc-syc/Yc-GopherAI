//go:build cgo

package image

import (
	"GopherAI/common/image"
	"io"
	"log/slog"
	"mime/multipart"
)

func RecognizeImage(file *multipart.FileHeader) (string, error) {

	modelPath := "/root/models/mobilenetv2/mobilenetv2-7.onnx"
	labelPath := "/root/imagenet_classes.txt"
	inputH, inputW := 224, 224

	recognizer, err := image.NewImageRecognizer(modelPath, labelPath, inputH, inputW)
	if err != nil {
		slog.Error("NewImageRecognizer failed", "error", err)
		return "", err
	}
	defer recognizer.Close()

	src, err := file.Open()
	if err != nil {
		slog.Error("file open failed", "error", err)
		return "", err
	}
	defer src.Close()

	buf, err := io.ReadAll(src)
	if err != nil {
		slog.Error("io.ReadAll failed", "error", err)
		return "", err
	}

	return recognizer.PredictFromBuffer(buf)
}
