package image

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/alibaba/sealer/common"
	imageUtils "github.com/alibaba/sealer/image/utils"
	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/mount"
	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
)

func buildBlobs(dig digest.Digest, size int64, mediaType string) distribution.Descriptor {
	return distribution.Descriptor{
		Digest:    dig,
		Size:      size,
		MediaType: mediaType,
	}
}

// GetImageHashList return image hash list
func GetImageHashList(image *v1.Image) (res []string, err error) {
	baseLayer, err := GetImageAllLayers(image)
	if err != nil {
		return res, fmt.Errorf("get base image failed error is :%v", err)
	}
	for _, layer := range baseLayer {
		if layer.Hash != "" {
			res = append(res, filepath.Join(common.DefaultLayerDir, layer.Hash))
		}
	}
	return
}

// GetImageAllLayers return all image layers, TODO need to refactor, do need to cut first one
func GetImageAllLayers(image *v1.Image) (res []v1.Layer, err error) {
	for {
		res = append(res, image.Spec.Layers[1:]...)
		if image.Spec.Layers[0].Value == common.ImageScratch {
			break
		}
		if len(res) > 128 {
			return nil, fmt.Errorf("current layer is exceed 128 layers")
		}
		i, err := imageUtils.GetImage(image.Spec.Layers[0].Value)
		if err != nil {
			return []v1.Layer{}, err
		}
		image = i
	}
	return
}

// GetClusterFileFromImage retrieve ClusterFile From image
func GetClusterFileFromImage(imageName string) string {
	clusterfile := GetClusterFileFromImageManifest(imageName)
	if clusterfile != "" {
		return clusterfile
	}

	clusterfile = GetClusterFileFromBaseImage(imageName)
	if clusterfile != "" {
		return clusterfile
	}
	return ""
}

// GetClusterFileFromImageManifest retrieve ClusterFile from image manifest(image yaml)
func GetClusterFileFromImageManifest(imageName string) string {
	//  find cluster file from image manifest
	imageMetadata, err := NewImageMetadataService().GetRemoteImage(imageName)
	if err != nil {
		return ""
	}
	if imageMetadata.Annotations == nil {
		return ""
	}
	clusterFile, ok := imageMetadata.Annotations[common.ImageAnnotationForClusterfile]
	if !ok {
		return ""
	}

	return clusterFile
}

// GetClusterFileFromBaseImage retrieve ClusterFile from base image, TODO need to refactor
func GetClusterFileFromBaseImage(imageName string) string {
	mountTarget, _ := utils.MkTmpdir()
	mountUpper, _ := utils.MkTmpdir()
	defer func() {
		err := utils.CleanDirs(mountTarget, mountUpper)
		if err != nil {
			logger.Warn(err)
		}
	}()

	if err := NewImageService().PullIfNotExist(imageName); err != nil {
		return ""
	}
	driver := mount.NewMountDriver()
	image, err := imageUtils.GetImage(imageName)
	if err != nil {
		return ""
	}

	layers, err := GetImageHashList(image)
	if err != nil {
		return ""
	}

	err = driver.Mount(mountTarget, mountUpper, layers...)
	if err != nil {
		return ""
	}
	defer func() {
		err := driver.Unmount(mountTarget)
		if err != nil {
			logger.Warn(err)
		}
	}()

	clusterFile := filepath.Join(mountTarget, "etc", common.DefaultClusterFileName)
	data, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		return ""
	}
	return string(data)
}
