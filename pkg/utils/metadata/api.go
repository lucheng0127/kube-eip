package metadata

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func ListMD() ([]*EipMetadata, error) {
	var mdfs []string
	var mds []*EipMetadata

	err := filepath.Walk(MDDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), "metadata") {
			mdfs = append(mdfs, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, filename := range mdfs {
		filename := fmt.Sprintf("%s/%s", MDDir, filename)
		md, err := ParseMD(filename)
		if err != nil {
			return nil, err
		}

		mds = append(mds, md)
	}

	return mds, nil
}
