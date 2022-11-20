package firebase

import (
	"context"

	xfirestore "cloud.google.com/go/firestore"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"google.golang.org/api/iterator"
)

type (
	OrderFs interface {
		Add(collection string, ID string, data entity.OrderTrxFs) error
		Update(collection string, doc string, data entity.OrderTrxFs) (err error)
	}

	orderFs struct {
		client *xfirestore.Client
	}
)

func NewOrderFs(client *xfirestore.Client) OrderFs {
	return &orderFs{client}
}

func (f *orderFs) Add(collection string, ID string, data entity.OrderTrxFs) error {
	err := f.client.RunTransaction(context.Background(), func(ctx context.Context, tx *xfirestore.Transaction) error {
		// check limit
		sanpShot, err := f.client.Collection(collection).Documents(ctx).GetAll()
		if err != nil {
			return err
		}

		// if len snapshot >= max data
		// before we add data, we need to delete first data
		if len(sanpShot) >= constant.FIRESTORE_MAX_DATA {
			query := f.client.Collection(collection).OrderBy("created", xfirestore.Asc).Limit(1)
			docs := tx.Documents(query)
			for {
				doc, err := docs.Next()
				if err != nil {
					if err == iterator.Done {
						break
					}
					if err != nil {
						return err
					}
				}
				// delete data
				err = tx.Delete(doc.Ref)
				if err != nil {
					return err
				}
			}
		}

		err = tx.Create(f.client.Collection(collection).Doc(ID), data)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *orderFs) Update(collection string, doc string, data entity.OrderTrxFs) (err error) {
	err = f.client.RunTransaction(context.Background(), func(ctx context.Context, tx *xfirestore.Transaction) (err error) {
		xcol := f.client.Collection(collection)
		xdoc := xcol.Doc(doc)
		_, err = xdoc.Set(ctx, data)

		return
	})

	return
}
