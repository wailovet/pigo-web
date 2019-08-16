package pigoutil

import (
	"bytes"
	"github.com/esimov/pigo/core"
	"image"
	"log"
)

var InstanceClassifier *pigo.Pigo

func InitCascade(cascadeData []byte) {
	var err error
	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	InstanceClassifier, err = pigo.Unpack(cascadeData)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

}

// GetImage retrieves and decodes the image file to *image.NRGBA type.
func getImage(input []byte) (*image.NRGBA, error) {
	src, _, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	img := pigo.ImgToNRGBA(src)

	return img, nil
}

func DetectFaces(imgData []byte) ([]pigo.Detection, error) {
	src, err := getImage(imgData)
	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	minSize := 20
	maxSize := 1000
	cParams := pigo.CascadeParams{
		MinSize:     minSize,
		MaxSize:     maxSize,
		ShiftFactor: 0.15,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	dets := InstanceClassifier.RunCascade(cParams, 0)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = InstanceClassifier.ClusterDetections(dets, 0.2)

	return dets, nil
}

type DetectFace struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// drawMarker mark the detected face region with the provided
// marker (rectangle or circle) and write it to io.Writer.
func FilterMarker(detections []pigo.Detection, qThresh float32) []DetectFace {
	var detectFaces []DetectFace
	for i := 0; i < len(detections); i++ {
		if detections[i].Q > qThresh {
			detectFaces = append(detectFaces, DetectFace{
				float64(detections[i].Col - detections[i].Scale/2),
				float64(detections[i].Row - detections[i].Scale/2),
				float64(detections[i].Scale),
				float64(detections[i].Scale),
			})
		}
	}
	return detectFaces
}
