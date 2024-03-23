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

func ParseLabel(filePath string) {

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)

	var (
		v               xml.StartElement
		ok              bool
		poolSize        int                         = runtime.NumCPU() * 5
		labelRepository *repository.LabelRepository = repository.GetLabelRepository()
	)

	buffer := make(chan *entity.Label, poolSize)
	quit := make(chan bool)

	for i := 0; i < poolSize; i++ {
		go labelRepository.SaveFromBuffer(buffer, quit)
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

		if v.Name.Local == "label" {
			var label entity.Label
			err := decoder.DecodeElement(&label, &v)
			if err != nil {
				log.Println("Error decoding label element:", err)
				continue

			}
			buffer <- &label
		}
	}

	close(buffer)

	for i := 0; i < poolSize; i++ {
		<-quit
	}
}
