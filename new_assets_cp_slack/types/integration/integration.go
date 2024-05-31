package type_integration

import "io/fs"

type SaveFile struct {
	FileName    string
	FileContent []byte
	Perm        fs.FileMode
}
