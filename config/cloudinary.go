package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ShikharY10/gbAPI/models"
	"github.com/ShikharY10/gbAPI/utils"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Cloudinary struct {
	CloudName            string
	APIKey               string
	APISecret            string
	ChatImageFolder      string
	ChatVideoFOlder      string
	ProfilePictureFolder string
	cloudinary           *cloudinary.Cloudinary
}

func InitCloudinary() Cloudinary {
	cloudName, found := os.LookupEnv("CLOUDINARY_NAME")
	if !found {
		panic("key -> CLOUDINARY_NAME is not found in .env")
	}

	cldApiKey, found := os.LookupEnv("CLOUDINARY_API_KEY")
	if !found {
		panic("key -> CLOUDINARY_API_KEY is not found in .env")
	}

	cldApiSecret, found := os.LookupEnv("CLOUDINARY_API_SECRET")
	if !found {
		panic("key -> CLOUDINARY_API_SECRET is not found in .env")
	}

	cldProfilePicFolder, found := os.LookupEnv("CLOUDINARY_PROFILE_PIC_FOLDER")
	if !found {
		panic("key -> CLOUDINARY_PROFILE_PIC_FOLDER is not found in .env")
	}

	cldChatImageFolder, found := os.LookupEnv("CLOUDINARY_CHAT_IMAGE_FOLDER")
	if !found {
		panic("key -> CLOUDINARY_CHAT_IMAGE_FOLDER is not found in .env")
	}

	cldChatVidoeFolder, found := os.LookupEnv("CLOUDINARY_CHAT_VIDEO_FOLDER")
	if !found {
		panic("key -> CLOUDINARY_CHAT_IMAGE_FOLDER is not found in .env")
	}

	var cloud Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, cldApiKey, cldApiSecret)
	if err != nil {
		panic(err)
	}

	cloud.APIKey = cldApiKey
	cloud.APISecret = cldApiSecret
	cloud.CloudName = cloudName
	cloud.ChatImageFolder = cldChatImageFolder
	cloud.ProfilePictureFolder = cldProfilePicFolder
	cloud.ChatVideoFOlder = cldChatVidoeFolder
	cloud.cloudinary = cld

	return cloud
}

func (cloud *Cloudinary) UploadUserAvatar(tempName string, imageData string, extension string) (*models.Avatar, error) {
	var image []byte = utils.Decode(imageData)
	f, err := os.Create("temp/" + tempName + "." + extension)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(image)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	uploadParam, err := cloud.cloudinary.Upload.Upload(
		ctx,
		"temp/"+tempName+"."+extension,
		uploader.UploadParams{Folder: cloud.ProfilePictureFolder},
	)
	if err != nil {
		return nil, err
	}
	os.Remove("temp/" + tempName + "." + extension)

	var avatar models.Avatar
	// avatar.Height = strconv.Itoa(uploadParam.Height)
	avatar.PublicId = uploadParam.PublicID
	avatar.SecureUrl = uploadParam.SecureURL
	// avatar.Width = strconv.Itoa(uploadParam.Width)
	return &avatar, nil
}

func (cloud *Cloudinary) UploadChatImage(imageData string, extension string) (secureUrl string, publicID string, err error) {
	var image []byte = utils.Decode(imageData)
	name := string(utils.GenerateAesKey(10))
	f, err := os.Create("temp_" + name + extension)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(image)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	uploadParam, err := cloud.cloudinary.Upload.Upload(
		ctx,
		"temp_"+name+extension,
		uploader.UploadParams{Folder: cloud.ChatImageFolder},
	)
	if err != nil {
		return "", "", err
	}
	return uploadParam.SecureURL, uploadParam.PublicID, nil
}

func (cloud *Cloudinary) UploadChatVideo() {}

func (cloud *Cloudinary) UploadProfilePicture() {}

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cld, err := cloudinary.NewFromParams("", "", "")
	if err != nil {
		return "", err
	}

	//upload file
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: "config.EnvCloudUploadFolder()"})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
