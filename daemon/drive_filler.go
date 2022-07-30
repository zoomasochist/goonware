package daemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	types "goonware/types"
)

// TODO: This should be configurable by the package
var filenameMessages []string = []string{
	"GOONWARE",
	"OWNEDBYGOONWARE",
	"PORNLOVESYOU",
	"GOON",
	"GIVEIN",
	"GIVEUP",
	"NEVERSTOPGOONING",
}

var allSubdirectories []string
var imageUrlCache []string // Up to 100 image urls
var imageDataCache []Image // Up to 5 images

type Image struct {
	Ext   string
	Bytes []byte
}

type E621File struct {
	Url string `json:"url"`
}

type E621Post struct {
	File E621File `json:"file"`
}

type E621Response struct {
	Posts []E621Post `json:"posts"`
}

// TODO: This function panics and I suppose it shouldn't
func DoDriveFiller(c *types.Config, pkg *types.EdgewarePackage) {
	var image Image

	if len(allSubdirectories) == 0 {
		if err := GetSubdirectories(c.DriveFillerBase); err != nil {
			panic(err)
		}
	}
	// This is interesting lol
	savePathNoExt := fmt.Sprintf("%s/%s%d",
		allSubdirectories[rand.Intn(len(allSubdirectories))],
		filenameMessages[rand.Intn(len(filenameMessages))],
		rand.Intn(int(time.Now().Unix())))

	if c.DriveFillerImageSource == types.DriveFillerImageSourceBooru {
		if len(imageUrlCache) == 0 {
			err := FillImageUrlCache(c)
			if err != nil {
				panic(err)
			}
		}

		if len(imageDataCache) == 0 {
			err := FillImageDataCache()
			if err != nil {
				panic(err)
			}
		}

		image, imageDataCache = imageDataCache[0], imageDataCache[1:]
	} else if c.DriveFillerImageSource == types.DriveFillerImageSourcePackage {

	}

	err := os.WriteFile(savePathNoExt+"."+image.Ext, image.Bytes, 0644)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(c.DriveFillerDelay) * time.Millisecond)
}

// Fills imageUrlCache with up to 100 image URLs
func FillImageUrlCache(c *types.Config) error {
	tags := ""
	if c.DriveFillerImageUseTags {
		tags = c.DriveFillerTags[rand.Intn(len(c.DriveFillerTags))]
	}

	if c.DriveFillerBooru == "https://e621.net/" {
		query := "https://e621.net/posts.json?limit=100&tags=rating:e+order:random+-animated" + tags
		if c.DriveFillerDownloadMinimumScoreToggle {
			query += fmt.Sprintf("+score:>=%d", c.DriveFillerDownloadMinimumScoreThreshold)
		}

		resp, err := MakeHttpRequest(query)
		if err != nil {
			return err
		}

		var marshalledResponse E621Response
		if err = json.Unmarshal(resp, &marshalledResponse); err != nil {
			return err
		}

		for _, response := range marshalledResponse.Posts {
			if response.File.Url != "" {
				imageUrlCache = append(imageUrlCache, response.File.Url)
			}
		}
	}

	return nil
}

func FillImageDataCache() error {
	imageUrls := imageUrlCache[:5]
	imageUrlCache = imageUrlCache[5:]

	for _, image := range imageUrls {
		imageBytes, err := MakeHttpRequest(image)
		if err != nil {
			return err
		}

		components := strings.Split(image, ".")
		imageDataCache = append(imageDataCache, Image{
			Ext:   components[len(components)-1],
			Bytes: imageBytes,
		})

		// E621s rate limit * 2. I assume most sites' are similar
		time.Sleep(1 * time.Second)
	}

	return nil
}

func GetSubdirectories(path string) error {
	return filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				return nil
			}

			allSubdirectories = append(allSubdirectories, path)
			return nil
		})
}

func MakeHttpRequest(uri string) ([]byte, error) {
	// TODO: I think we can reuse the client
	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "Goonware")

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
