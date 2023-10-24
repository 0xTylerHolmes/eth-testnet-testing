package consensus_objects

import (
	"github.com/attestantio/go-eth2-client/spec/deneb"
	fuzz "github.com/google/gofuzz"
)

// RandomBlobSideCar creates a totally random but syntactically valid BlobSideCar
func RandomBlobSideCar() *deneb.BlobSidecar {
	var blobSideCar deneb.BlobSidecar
	f := fuzz.New().NilChance(0)
	for true {
		f.Fuzz(&blobSideCar)
		_, err := blobSideCar.MarshalSSZ()
		if err == nil {
			break
		}
	}
	return &blobSideCar
}
