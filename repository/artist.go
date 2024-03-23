package repository

import (
	"context"
	"log"
	"totoro1302/goparser/entity"

	"github.com/jackc/pgx/v5"
)

type ArtistRepository struct {
	db *DbPool
}

const insertArtistStatement string = `INSERT INTO artist (id_discogs, name, realname, namevariations, urls, profile)
values (@idDiscogs, @name, @realName, @nameVariations, @urls, @profile)
ON CONFLICT (id_discogs)
DO UPDATE SET name = @name, realname = @realName, namevariations = @nameVariations, urls = @urls, profile = @profile;`

const insertArtistBatchSize = 20

func GetArtistRepository() *ArtistRepository {
	pool := GetDbPoolConnection("postgres://discogs_rw:motdepasse@myPostgres:5432/discogs?pool_max_conns=20")

	return &ArtistRepository{pool}
}

func (r *ArtistRepository) SaveFromBuffer(buffer chan *entity.Artist, quit chan bool) {

	batch := &pgx.Batch{}
	counter := 0

	for artist := range buffer {

		counter++

		args := pgx.NamedArgs{
			"idDiscogs":      artist.Id,
			"name":           artist.Name,
			"realName":       artist.Realname,
			"nameVariations": artist.NameVariations,
			"urls":           artist.Urls,
			"profile":        artist.Profile,
		}

		batch.Queue(insertArtistStatement, args)

		if counter >= insertArtistBatchSize {
			batchResult := r.db.Pool.SendBatch(context.Background(), batch)
			err := batchResult.Close()
			if err != nil {
				log.Printf("Failed to send batch: %s, artistId: %d, artistName: %s\n", err, artist.Id, artist.Name)
			}
			batch = &pgx.Batch{}
			counter = 0
		}
	}

	// Send the last batch if any queue
	if batch.Len() > 0 {
		batchResult := r.db.Pool.SendBatch(context.Background(), batch)
		batchResult.Close()
	}

	quit <- true
}
