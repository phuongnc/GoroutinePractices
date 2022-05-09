package src

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	aws "resizeimage/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/h2non/bimg" //export CGO_CFLAGS_ALLOW=-Xpreprocessor
)

func (obj *ResizeImagesReq) ResizeImages() []*ResizeImageItem {
	var wg sync.WaitGroup
	chanels := make(chan *ResizeImageItem, len(obj.Data))

	for _, item := range obj.Data {
		wg.Add(1)
		go ResizeImageWithProfile(item, chanels, &wg)
	}

	wg.Wait()
	close(chanels)

	var result []*ResizeImageItem
	for c := range chanels {
		result = append(result, c)
	}
	return result
}

//resize image with profiles from configuration
func ResizeImageWithProfile(obj ResizeImageItem, channel chan<- *ResizeImageItem, wg *sync.WaitGroup) {
	defer wg.Done()

	//Download origin image
	err := downloadFile(obj.Key)
	if err != nil {
		obj.Message = append(obj.Message, err.Error())
		LogWithField(obj, err).Error(err.Error())
		channel <- &obj
		return
	}

	//Get all profiles from config base on type
	var imageProfiles []int
	if obj.Type == "post" {
		imageProfiles = GetConfig().PostImageProfiles
	} else if obj.Type == "avatar" {
		imageProfiles = GetConfig().AvatarImageProfiles
	}

	//resize earch profile
	originBuffer, err := bimg.Read(obj.Key)
	for _, size := range imageProfiles {
		var err error
		//retry 3 times if failed
		for i := 0; i < GetConfig().RetryTimes; i++ {
			newImageProfile, err := resizeImageWithSize(obj.Key, originBuffer, size)
			if err == nil {
				obj.Profile = append(obj.Profile, newImageProfile)
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			obj.Message = append(obj.Message, err.Error())
			LogWithField(obj, err).Error(err.Error())
		}
	}

	//Upload to s3
	resObj := obj.uploadFileToStorage(obj.Profile)
	channel <- resObj
	return
}

//resize image function
func resizeImageWithSize(key string, orgImageBuffer []byte, size int) (string, error) {
	orgImage := bimg.NewImage(orgImageBuffer)
	originImageSize, _ := orgImage.Size()
	newWidth := 0
	newHeight := 0
	//determine new size
	if originImageSize.Width > size || originImageSize.Height > size {
		if originImageSize.Width >= originImageSize.Height {
			newWidth = size
			newHeight = (size * originImageSize.Height) / originImageSize.Width
		} else {
			newHeight = size
			newWidth = (size * originImageSize.Width) / originImageSize.Height
		}
	}
	//prepare new image name
	folder := path.Dir(key)
	orgImageName := path.Base(key)
	newImageName := strings.ReplaceAll(orgImageName, "original", strconv.Itoa(size))
	if ext := path.Ext(key); ext != ".jpg" {
		newImageName = strings.ReplaceAll(newImageName, ext, ".jpg")
	}
	//resize image
	if newWidth != 0 || newHeight != 0 {
		newImageBuffer, err := orgImage.Resize(newWidth, newHeight)
		if err != nil {
			return newImageName, err
		}
		err = bimg.Write(path.Join(folder, newImageName), newImageBuffer)
		if err != nil {
			return newImageName, err
		}
	} else {
		err := bimg.Write(path.Join(folder, newImageName), orgImageBuffer)
		if err != nil {
			return newImageName, err
		}
	}
	return path.Join(folder, newImageName), nil
}

//download file
func downloadFile(key string) error {
	imageUrl := GetConfig().StorageUrl + "/" + key
	response, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Error download image status: " + response.Status)
	}
	//Create a empty file
	if err := os.MkdirAll(filepath.Dir(key), 0770); err != nil {
		return err
	}
	file, err := os.Create(key)
	if err != nil {
		return err
	}
	defer file.Close()
	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

//upload array file
func (obj *ResizeImageItem) uploadFileToStorage(filePaths []string) *ResizeImageItem {
	var wg sync.WaitGroup
	chanels := make(chan UploaderRes, len(filePaths))

	for _, filePath := range filePaths {
		wg.Add(1)

		maps := make(map[string]interface{})
		maps["bucket_name"] = obj.Bucket
		maps["key"] = filePath

		go uploadFile(maps, chanels, &wg)
	}

	wg.Wait()
	close(chanels)

	//read from channels
	uploadedProfile := make([]string, 0)
	for c := range chanels {
		if c.err != nil {
			obj.Message = append(obj.Message, c.fileName+":"+c.err.Error())
			LogWithField(obj, c.err).Error(c.err.Error())
		} else {
			uploadedProfile = append(uploadedProfile, c.fileName)
		}
	}
	obj.Profile = uploadedProfile
	return obj
}

//upload single file
func uploadFile(maps map[string]interface{}, channel chan<- UploaderRes, wg *sync.WaitGroup) {
	defer wg.Done()

	file, _ := os.Open(maps["key"].(string))
	defer file.Close()

	uploadedFileName, err := aws.UploadFileToS3(file, maps)

	response := UploaderRes{fileName: uploadedFileName, err: err}
	channel <- response
}

//removeLocalFiles(obj.Data[0].Key)
func removeLocalFiles(key string) {
	os.RemoveAll(filepath.Dir(key))
}
