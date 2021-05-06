package image

import (
	"context"

	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
)

type blobService struct {
	descriptors map[digest.Digest]distribution.Descriptor
}

func (bs *blobService) Get(context.Context, digest.Digest) ([]byte, error) {
	return []byte{}, nil
}

func (bs *blobService) Stat(_ context.Context, dgst digest.Digest) (distribution.Descriptor, error) {
	if descriptor, ok := bs.descriptors[dgst]; ok {
		return descriptor, nil
	}
	return distribution.Descriptor{}, distribution.ErrBlobUnknown
}

func (bs *blobService) Open(context.Context, digest.Digest) (distribution.ReadSeekCloser, error) {
	return nil, nil
}

func (bs *blobService) Put(_ context.Context, mediaType string, p []byte) (distribution.Descriptor, error) {
	d := distribution.Descriptor{
		Digest:    digest.FromBytes(p),
		Size:      int64(len(p)),
		MediaType: mediaType,
	}
	bs.descriptors[d.Digest] = d
	return d, nil
}

func (bs *blobService) Create(context.Context, ...distribution.BlobCreateOption) (distribution.BlobWriter, error) {
	return nil, nil
}

func (bs *blobService) Resume(context.Context, string) (distribution.BlobWriter, error) {
	return nil, nil
}
