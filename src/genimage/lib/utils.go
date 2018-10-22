package lib

import(
	"math/rand"
	"time"
	"image"
	"net/http"
	_ "image/png"
	_ "image/jpeg"
	"log"
	"github.com/fogleman/gg"
)

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func Imagefetch(url string) (Image image.Image, err error) {
	imageHTTPdata, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer imageHTTPdata.Body.Close()

	Image, _, err = image.Decode(imageHTTPdata.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return Image, nil
}

func DrawRoundedRectangles(image *gg.Context, x float64, y float64, width float64, height float64, bartype string, chamfer float64) (err error) {
	switch colorInput := bartype; colorInput {
	case "box":
		image.SetRGB(1, 1, 1)
	case "black":
		image.SetRGB(0, 0, 0)
	case "level":
		image.SetRGB(0, 0, 1)
	case "health":
		image.SetRGB(1, 0, 0)
	case "hunger":
		image.SetRGB(0, 1, 0)
	case "poop":
		image.SetRGB(1, 1, 1)
	default:
		image.SetRGB(0, 0, 0)
	}

	image.DrawRoundedRectangle(x, y, width, height, chamfer)
	image.Fill()

	return nil
}
