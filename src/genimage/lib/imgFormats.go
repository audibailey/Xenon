package lib

import(
	"image"
	"log"
	"fmt"

	"github.com/nfnt/resize"
)

func PokemonImage(pokemon string) (resizedPokemon image.Image, pokemonWidth int, err error){
	pokemonImgURL := fmt.Sprintf("https://assets.pokemon.com/assets/cms2/img/pokedex/full/%v.png", pokemon)
	pokemonImage, err := Imagefetch(pokemonImgURL)
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}

	pokemonImageBounds := pokemonImage.Bounds()
	pokemonWidth = pokemonImageBounds.Max.X/2

	resizedPokemon = resize.Resize(uint(pokemonWidth), 0, pokemonImage, resize.Lanczos3)

	return resizedPokemon, pokemonWidth, nil
}

func BackgroundImage(background string) (backgroundImage image.Image, err error){
	var backgroundImgURL string
	switch backgroundImgURLInput := background; backgroundImgURLInput {
	case "test-bg":
		backgroundImgURL = "https://st3.depositphotos.com/5590000/12620/v/950/depositphotos_126202658-stock-illustration-cartoon-vector-illustration-interior-library.jpg"
	default:
		return nil, err
	}

	backgroundImage, err = Imagefetch(backgroundImgURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return backgroundImage, nil
}

func IconImages(icon string) (resizedImage image.Image, err error) {
	var iconImgURL string
	switch iconImgURLInput := icon; iconImgURLInput {
	case "heart":
		iconImgURL = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/160/apple/118/heavy-black-heart_2764.png"
	case "poop":
		iconImgURL = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/120/apple/129/pile-of-poo_1f4a9.png"
	case "hunger":
		iconImgURL = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/120/apple/129/poultry-leg_1f357.png"
	default:
		return nil, err
	}

	iconImage, err := Imagefetch(iconImgURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	iconImageBounds := iconImage.Bounds()
	iconWidth := iconImageBounds.Max.X/12

	resizedImage = resize.Resize(uint(iconWidth), 0, iconImage, resize.Lanczos3)
	return resizedImage, nil
}
