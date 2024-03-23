package services

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"runtime"
	"totoro1302/goparser/entity"
	"totoro1302/goparser/repository"
)

func ParseArtist(filePath string) {

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)

	var (
		v                xml.StartElement
		ok               bool
		poolSize         int                          = runtime.NumCPU() * 5
		artistRepository *repository.ArtistRepository = repository.GetArtistRepository()
	)

	buffer := make(chan *entity.Artist, poolSize)
	quit := make(chan bool)

	for i := 0; i < poolSize; i++ {
		go artistRepository.SaveFromBuffer(buffer, quit)
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error getting token: %v\n", err)
			break
		}

		if v, ok = token.(xml.StartElement); !ok {
			continue
		}

		if v.Name.Local == "artist" {
			log.Println("{\"type\": \"debug\", \"timestamp\": \"2020-03-04 21:16:20+00:00 GMT\", \"message\": \"OK\"}")
			var artist entity.Artist
			err := decoder.DecodeElement(&artist, &v)
			if err != nil {
				log.Println("Error decoding artist element:", err)
				continue

			}
			buffer <- &artist
		}
	}

	close(buffer)

	for i := 0; i < poolSize; i++ {
		<-quit
	}
}
