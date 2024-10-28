package post

const (
	fileTypeJPG = "image/jpeg"
	fileTypeBMP = "image/bmp"
	fileTypePNG = "image/png"

	extJpeg = ".jpeg"

	imgNameFmt    = "%s_%s%s"
	displayPrefix = "display"
	srcPrefix     = "src"

	srcImgDir     = "src_images"
	displayImgDir = "display_images"

	maxFileSize       = int64(100 * 1024 * 1024) // 100 mb
	defaultWidthRatio = 600
	defaultHeightRatio

	maxPosts           = 10
	maxCommentsPerPost = 2
)

var (
	allowedFormats = map[string]bool{
		fileTypeJPG: true,
		fileTypeBMP: true,
		fileTypePNG: true,
	}
)

type CreatePostInput struct {
	Caption        string `json:"caption"`
	Creator        string `json:"creator" validate:"required"`
	SrcImgPath     string `form:"-"`
	DisplayImgPath string `form:"-"`
}
