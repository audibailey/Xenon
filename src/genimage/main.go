package main

import(
	"log"
	"image/jpeg"
	_ "image/png"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"regexp"
	"context"
	"encoding/base64"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"testservice/src/genimage/lib"
)

// TODO: Struct the Datatypes

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Variable Assignments
	/*age, err := strconv.Atoi(params["age"])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}*/
	level, err := strconv.Atoi(request.PathParameters["level"])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}	
	
	XP := regexp.MustCompile("-").Split(request.PathParameters["xp"], 2)
	levelxp, err := strconv.Atoi(XP[0])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	levelxpmax, err := strconv.Atoi(XP[1])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	HP := regexp.MustCompile("-").Split(request.PathParameters["hp"], 2)
	health, err := strconv.Atoi(HP[0])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	healthmax, err := strconv.Atoi(HP[1])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	HUNGER := regexp.MustCompile("-").Split(request.PathParameters["hunger"], 2)
	hunger, err := strconv.Atoi(HUNGER[0])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	hungermax, err := strconv.Atoi(HUNGER[1])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	DIRT := regexp.MustCompile("-").Split(request.PathParameters["dirt"], 2)
	poop, err := strconv.Atoi(DIRT[0])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	poopmax, err := strconv.Atoi(DIRT[1])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	text := fmt.Sprintf("Level: %v", level)

	// get Images and set font
	resizedPokemon, pokemonWidth, err := functs.PokemonImage(request.PathParameters["pokemon"])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	backgroundImage, err := functs.BackgroundImage(request.PathParameters["background"])
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	resizedHeart, err := functs.IconImages("heart")
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	resizedLeg, err := functs.IconImages("hunger")
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	resizedPoop, err := functs.IconImages("poop")
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 10,
	})

	// Image Generation
	finalImage := gg.NewContextForImage(backgroundImage)
	finalImageWidth := finalImage.Width()
	finalImageHeight := finalImage.Height()

	// Add Pokemon to Image
	randomX := functs.Random(pokemonWidth, finalImageWidth-pokemonWidth)
	finalImage.DrawImage(resizedPokemon, randomX, 500)

	// Draw InfoBox
	infoboxwidth := float64(200)
	infoboxheight := float64(75)
	infoboxX := float64((finalImageWidth/2)-(int(infoboxwidth/2)))
	infoboxY := float64(finalImageHeight-int(infoboxheight))
	if err := functs.DrawRoundedRectangles(finalImage, infoboxX, infoboxY, infoboxwidth, infoboxheight, "box", 5); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Add Level Text
	finalImage.SetRGB(0, 0, 0)
	finalImage.SetFontFace(face)
	finalImage.DrawStringAnchored(text, float64((finalImageWidth/2)-(int(infoboxwidth/2))+5), float64(finalImageHeight-int(infoboxheight)+10), 0.0, 0.0)

	// Outer Level Bar
	levelmaxbar := infoboxwidth-60
	levebarmaxX := float64((finalImageWidth/2)-(int(infoboxwidth/2))+50)
	levelbarmaxY := float64(finalImageHeight-int(infoboxheight)+2)
	if err := functs.DrawRoundedRectangles(finalImage, levebarmaxX, levelbarmaxY, levelmaxbar, 10, "black", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Inner Level Bar
	levelbar := ((float64(levelxp)/float64(levelxpmax-2)))*float64(levelmaxbar-2)
	levelbarX := float64((finalImageWidth/2)-(int(infoboxwidth/2))+51)
	levelbarY := float64(finalImageHeight-int(infoboxheight)+3)
	if err := functs.DrawRoundedRectangles(finalImage, levelbarX, levelbarY, levelbar, 8, "level", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Common Icon Bar Values
	commonIconX := int((finalImageWidth/2)-(int(infoboxwidth/2))+5)
	commonIconBarOuterX := float64((finalImageWidth/2)-(int(infoboxwidth/2))+20)
	commonIconBarInnerX := float64((finalImageWidth/2)-(int(infoboxwidth/2))+21)
	othermaxbar := infoboxwidth-30

	// health
	finalImage.DrawImage(resizedHeart, commonIconX, int(finalImageHeight-int(infoboxheight))+20)
	healthbar := ((float64(health)/float64(healthmax)))*float64(othermaxbar-2)
	healthbarmaxY := float64(finalImageHeight-int(infoboxheight)+20)
	healthbarY := float64(finalImageHeight-int(infoboxheight)+21)

	if err := functs.DrawRoundedRectangles(finalImage, commonIconBarOuterX, healthbarmaxY, othermaxbar, 10, "black", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	if err := functs.DrawRoundedRectangles(finalImage, commonIconBarInnerX, healthbarY, healthbar, 8, "health", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Hunger
	finalImage.DrawImage(resizedLeg, commonIconX, int(finalImageHeight-int(infoboxheight))+40)
	hungerbar := ((float64(hunger)/float64(hungermax)))*float64(othermaxbar-2)
	hungerbarmaxY := float64(finalImageHeight-int(infoboxheight)+40)
	hungerbarY := float64(finalImageHeight-int(infoboxheight)+41)

	if err := functs.DrawRoundedRectangles(finalImage, commonIconBarOuterX, hungerbarmaxY, othermaxbar, 10, "black", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	if err := functs.DrawRoundedRectangles(finalImage, commonIconBarInnerX, hungerbarY, hungerbar, 8, "hunger", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Poop
	finalImage.DrawImage(resizedPoop, commonIconX, int(finalImageHeight-int(infoboxheight))+60)
	poopbarmaxY := float64(finalImageHeight-int(infoboxheight)+60)
	if err := functs.DrawRoundedRectangles(finalImage, commonIconBarOuterX, poopbarmaxY, othermaxbar, 10, "black", 2); err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	if poop != 0 {
		poopbar := ((float64(poop)/float64(poopmax)))*float64(othermaxbar-2)
		poopbarY := float64(finalImageHeight-int(infoboxheight)+61)
		if err := functs.DrawRoundedRectangles(finalImage, commonIconBarInnerX, poopbarY, poopbar, 8, "poop", 2); err != nil {
			log.Fatal(err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	}

	jpegprint := finalImage.Image()
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, jpegprint, &jpeg.Options{jpeg.DefaultQuality})

	resp := events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					IsBase64Encoded: true,
					Body: base64.StdEncoding.EncodeToString(buf.Bytes()),
					Headers: map[string]string{
									"Content-Type":           "image/jpeg",
									"X-MyCompany-Func-Reply": "hello-handler",
					},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

