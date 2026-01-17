package media

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
)

type File struct {
	Name      string
	FileName  string
	Path      string
	Extension string
	Size      int64
}

func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errs.BuildError(err, "error opening file")
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", errs.BuildError(err, "error calculating MD5 hash")
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func GetRelativePath(root, path string) string {
	return strings.Replace(path, root, "", 1)
}

func GetTitleOfFile(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return filename
	}

	parts = parts[:len(parts)-1]

	return strings.Join(parts, ".")
}

func GetFileSize(path string) (int64, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return int64(math.Abs(float64(fileinfo.Size()))), nil
}

func GetFilesByExtensions(root string, extensions []string) (ret []File, reterr error) {
	reterr = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if slices.Contains(extensions, strings.ToLower(filepath.Ext(d.Name()))) {
				file, err := GetFileInformation(path)
				if err != nil {
					return errors.Join(reterr, errs.BuildError(err, "GetFilesByExtensions"))
				}

				ret = append(ret, *file)
			}
		}

		return nil
	})

	return ret, reterr
}

func FindNonExistentMedia(existingVideos []model.Media, files []File) []model.Media {
	nonExsistentVideos := []model.Media{}
	for _, v := range existingVideos {
		if !slices.ContainsFunc(files, func(mediaFile File) bool {
			return mediaFile.Path == v.Path
		}) {
			nonExsistentVideos = append(nonExsistentVideos, v)
		}
	}
	return nonExsistentVideos
}

func GetFileInformation(p string) (*File, error) {
	fileSize, err := GetFileSize(p)
	if err != nil {
		return nil, errs.BuildError(err, "could not determine file size for: %v", p)
	}
	base := filepath.Base(p)
	file := File{
		Name:     GetTitleOfFile(base),
		FileName: base,
		Path:     p,
		Size:     fileSize,
	}
	return &file, nil
}
