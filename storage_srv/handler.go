package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"baymax/errors"
	protocol "baymax/storage_srv/protocol/storage"
	"golang.org/x/image/webp"
)

type storageHandler int

func (s *storageHandler) StorePhoto(req *protocol.StorePhotoArgs, reply *protocol.StorePhotoReply) error {

	var (
		err         error
		info        image.Config
		fileName    string
		fileKey     string
		contentType string
	)

	if req.Photo == nil {
		log.Errorf("upload photo: %v", req.UserId)
		return errors.BadRequest("file_empty", "照片为空")
	}

	fileName = newFileName()

	r := bytes.NewReader(req.Photo)

	if r.Size() > Config.Storage.MaxSize {
		log.Errorf("photo big: %v", req.UserId)
		return errors.BadRequest("too_big_file", "照片文件太大")
	}

	contentType = "image/" + req.FileType

	if req.FileType == "png" {
		info, err = png.DecodeConfig(r)
	} else if req.FileType == "jpg" || req.FileType == "jpeg" {
		contentType = "image/jpeg"
		info, err = jpeg.DecodeConfig(r)
	} else if req.FileType == "gif" {
		info, err = gif.DecodeConfig(r)
	} else if req.FileType == "webp" {
		info, err = webp.DecodeConfig(r)
	} else {
		err = errors.BadRequest("unsupport_file_type", "不支持的文件类型")
	}

	if err != nil {
		return err
	}

	fileKey = fileName + "." + req.FileType
	err = storePicture(fileKey, bytes.NewReader([]byte(req.Photo)), contentType)

	if err != nil {
		return err
	}

	reply.Filekey = fileKey
	reply.Url = getUrl(fileKey)
	reply.Width = info.Width
	reply.Height = info.Height
	reply.Suffixes = Config.Url.Suffixes

	return nil
}
