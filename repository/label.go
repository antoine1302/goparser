package repository

import (
	"context"

	"totoro1302/goparser/entity"

	"github.com/jackc/pgx/v5"
)

type LabelRepository struct {
	db *DbPool
}

const insertLabelStatement string = `INSERT INTO label (id_discogs, name, profile, contactinfo, urls)
values (@idDiscogs, @name, @profile, @contactInfo, @urls)
ON CONFLICT (id_discogs)
DO UPDATE SET name = @name, profile = @profile, contactinfo = @contactInfo, urls = @urls;`

const insertLabelBatchSize = 50

func GetLabelRepository() *LabelRepository {
	pool := GetDbPoolConnection("postgres://discogs_rw:motdepasse@myPostgres:5432/discogs?pool_max_conns=20")

	return &LabelRepository{pool}
}

func (r *LabelRepository) SaveFromBuffer(buffer chan *entity.Label, quit chan bool) {

	batch := &pgx.Batch{}
	counter := 0

	for label := range buffer {

		counter++

		args := pgx.NamedArgs{
			"idDiscogs":   label.Id,
			"name":        label.Name,
			"profile":     label.Profile,
			"contactInfo": label.ContactInfo,
			"urls":        label.Urls,
		}
		batch.Queue(insertLabelStatement, args)

		if counter%insertLabelBatchSize == 0 {
			batchResult := r.db.Pool.SendBatch(context.Background(), batch)
			batchResult.Close()
			batch = &pgx.Batch{}
			counter = 0
		}
	}

	if batch.Len() > 0 {
		batchResult := r.db.Pool.SendBatch(context.Background(), batch)
		batchResult.Close()
	}

	quit <- true
}
