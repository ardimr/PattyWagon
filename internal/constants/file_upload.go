package constants

var (
	AllowedExtensions []string = []string{".jpg", ".jpeg"}

	MaxUploadSizeInBytes int64 = 2097152
	MinUploadSizeInBytes int64 = 102400
	MaxUploadForm              = 3 << 20
)
