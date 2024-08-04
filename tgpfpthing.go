package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"gorm.io/driver/sqlite"

	_ "github.com/joho/godotenv/autoload"
)

// https://core.telegram.org/method/photos.uploadProfilePhoto
// https://core.telegram.org/method/photos.updateProfilePhoto

func main() {
	appId, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		panic(err)
	}

	client, err := gotgproto.NewClient(
		appId,
		os.Getenv("API_HASH"),
		gotgproto.ClientTypePhone(os.Getenv("PHONE_NUMBER")),
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("../session.db")),
		},
	)
	if err != nil {
		panic(err)
	}

	id := client.Self.ID
	fmt.Printf("Client: @%v with ID #%v has been started...\n", client.Self.Username, id)
	ctx := client.CreateContext()

	// get all user photos

	userPhotos, err := client.Client.API().PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: &tg.InputUserSelf{}, Offset: 0, MaxID: 0, Limit: 999})

	if err != nil {
		panic(err)
	}

	// fmt.Printf("User photos: %v\n", userPhotos)

	for i, photo := range userPhotos.GetPhotos() {
		photo := photo.(*tg.Photo)

		// telegram doesn't store the original profile pictures, so determine the best available size to download

		// https://core.telegram.org/api/files#image-thumbnail-types
		// s 	box 	100x100
		// m 	box 	320x320
		// x 	box 	800x800
		// y 	box 	1280x1280
		// w 	box 	2560x2560
		// a 	crop 	160x160
		// b 	crop 	320x320
		// c 	crop 	640x640
		// d 	crop 	1280x1280

		// therefore w > y > d > x > c > b > m > a > s

		photoSizes := make([]string, len(photo.Sizes))
		for _, size := range photo.Sizes {
			photoSizes = append(photoSizes, size.GetType())
		}

		var biggestSize string
		for _, size := range []string{"w", "y", "d", "x", "c", "b", "m", "a", "s"} {
			if slices.Contains(photoSizes, size) {
				biggestSize = size
				break
			}
		}

		if biggestSize == "" {
			fmt.Printf("No available sizes for photo %v\n", i)
			continue
		}

		// TODO: download video as video
		// TODO: goroutine for downloading multiple photos, downloads quite quickly as is but why not

		downloader.NewDownloader().Download(ctx.Raw, &tg.InputPhotoFileLocation{
			ID:            photo.ID,
			AccessHash:    photo.AccessHash,
			FileReference: photo.FileReference,
			ThumbSize:     biggestSize,
		}).ToPath(ctx, fmt.Sprintf("./downloaded/photo_%v.jpg", i))
	}

	fmt.Print("Client is idle...\n")
	client.Idle()
}
