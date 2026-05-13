package webServer

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"github.com/oklog/ulid"
	"github.com/the-suleiman/nla_framework/utils"
)

const IMAGE_DIR = "../image"
const STAT_IMAGE_PATH = "/stat-img"

func uploadImage(c *gin.Context) {
	// извлекаем название таблицы и id записи, к которой крепится файл
	tableName, _ := c.GetPostForm("tableName")
	{
		if len(tableName) == 0 {
			utils.HttpError(c, http.StatusBadRequest, "missed tableName")
			return
		}
	}
	tableId, _ := c.GetPostForm("tableId")
	{
		if len(tableId) == 0 {
			utils.HttpError(c, http.StatusBadRequest, "missed tableId")
			return
		}
	}
	// извлекаем минимальную ширину фото для сжатия. Если 0, то не сжимаем
	width := 0
	if widthStr, ok := c.GetPostForm("width"); ok {
		w, err := strconv.Atoi(widthStr)
		if err == nil {
			width = w
		}
	}
	// извлекаем параметры для crop. Например, 300x400
	crop := []int{}
	if cropStr, ok := c.GetPostForm("crop"); ok {
		widthStr := strings.Split(cropStr, "x")
		// перекладываем []string -> []int
		if len(widthStr) == 2 {
			for _, v := range widthStr {
				w, err := strconv.Atoi(v)
				if err == nil {
					crop = append(crop, w)
				}
			}
		}
	}
	path := fmt.Sprintf("%s/%s/%s", IMAGE_DIR, tableName, tableId)
	// optional post form: when truthy, save under a sanitized original basename (see createOutputFile); default is random ulid name
	preserveOriginalFileName := false
	if v, ok := c.GetPostForm("preserveOriginalFileName"); ok {
		v = strings.TrimSpace(strings.ToLower(v))
		if v == "1" || v == "true" || v == "yes" || v == "y" || v == "on" {
			preserveOriginalFileName = true
		}
	}
	saveImage(c, path, "", width, crop, preserveOriginalFileName)
}

func saveImage(c *gin.Context, path, filePrefix string, width int, crop []int, preserveOriginalFileName bool) {
	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName string

	if len(form.File) == 0 {
		utils.HttpError(c, http.StatusBadRequest, "list of files is empty")
		return
	}
	// берем первое имя файла из присланного списка
	for key := range form.File {
		if len(fileName) > 0 {
			continue
		}
		fileName = key
	}
	// извлекаем содержание присланного файла по названию файла
	file, fileHeader, err := c.Request.FormFile(fileName)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadFile c.Request.FormFile error: %s", err.Error()))
		return
	}
	defer file.Close()

	// извлекаем расширение файла
	imgExt := "jpeg"
	isImgExtInContentType := false
	contentType, _ := c.GetPostForm("Content-Type")
	if len(contentType) > 0 {
		arr := strings.Split(contentType, "/")
		if len(arr) > 1 {
			imgExt = arr[1]
			isImgExtInContentType = true
		}
	}

	// в случае если расширение файла не найдено в Content-Type, то извлекаем его из названия файла
	if !isImgExtInContentType {
		arr := strings.Split(fileName, ".")
		ext := arr[len(arr)-1]
		if ext == "png" || ext == "jpeg" || ext == "jpg" || ext == "gif" {
			imgExt = ext
		}
	}

	var isSaveAsIs bool
	// для png, gif если не указаны параметры преобразования, то сохраняем их без декодирования. Иначе анимация gif теряется, а у png теряется прозрачный фон
	if imgExt == "png" || imgExt == "gif" {
		if (crop == nil || len(crop) != 2) && width == 0 {
			isSaveAsIs = true
		}
	}

	// перекодируем файл в картинку
	var img image.Image
	var resizedImg image.Image
	if !isSaveAsIs {
		switch imgExt {
		case "jpeg", "jpg":
			img, err = jpeg.Decode(file)
		case "png":
			img, err = png.Decode(file)
		case "gif":
			img, err = gif.Decode(file)
		default:
			err = errors.New("Unsupported file type")
		}
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage jpeg.Decode error: %s", err.Error()))
			return
		}

		// сжатие размеров картинки до минимума - 500 или фактический размер
		imgWidth := uint(utils.MinInt(width, img.Bounds().Max.X))
		resizedImg = resize.Resize(imgWidth, 0, img, resize.Lanczos3)

		// если необходимо обрезать
		if crop != nil && len(crop) == 2 {
			analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
			topCrop, _ := analyzer.FindBestCrop(resizedImg, crop[0], crop[1])
			type SubImager interface {
				SubImage(r image.Rectangle) image.Image
			}
			resizedImg = resizedImg.(SubImager).SubImage(topCrop)
		}
	}

	// создаем директорию, если еще не создана
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.MkdirAll error: %s", err))
		return
	}

	// открываем файл для сохранения картинки
	fullFileName, fileOnDisk, err := createOutputFile(path, filePrefix, imgExt, preserveOriginalFileName, fileHeader)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage createOutputFile err: %s", err))
		return
	}
	defer fileOnDisk.Close()

	// сохранение файла
	// два варианта
	if isSaveAsIs {
		// без перкодировки - сохраняем как есть, только заменяем имя
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage io.ReadAll(file) error: %s", err))
			return
		}
		_, err = fileOnDisk.Write(fileBytes)
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage fileOnDisk.Write error: %s", err))
			return
		}
	} else {
		// с перекодировкой
		err = jpeg.Encode(fileOnDisk, resizedImg, nil)
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage jpeg.Encode err: %s", err))
			return
		}
	}

	// возвращаем ссылку на файл
	utils.HttpSuccess(c, map[string]string{"file": fmt.Sprintf("%s/%s", strings.Replace(path, IMAGE_DIR, STAT_IMAGE_PATH, 1), fullFileName)})
}

// загрузка аватарки
func uploadProfileImage(c *gin.Context) {
	if userId, ok := utils.ExtractUserIdString(c); ok {
		path := fmt.Sprintf("%s/profile", IMAGE_DIR)
		prefix := fmt.Sprintf("id_%s_", userId)
		saveImage(c, path, prefix, 200, []int{200, 200}, false)
	}
}

// json body uses ContextJsonParamFldParam: params.filename is /stat-img/... or resolvable disk path; one 270° rotation in place (overwrites file)
func rotateImage(c *gin.Context) {
	params, _ := c.Get(utils.ContextJsonParamFldParam)

	var filename string
	if filenameValue, ok := params.(map[string]interface{})["filename"].(string); ok {
		filename = filenameValue
	} else {
		utils.HttpError(c, http.StatusBadRequest, "filename is nil")
		return
	}

	filename = imageStatPathToFilePath(filename)
	if !fileExists(filename) {
		utils.HttpError(c, http.StatusBadRequest, "file not found. filename: "+filename)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, "image open error:"+err.Error())
		return
	}
	img, format, err := image.Decode(file)
	file.Close()
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, "image decode error:"+err.Error())
		return
	}

	out, err := os.Create(filename)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, "image create error:"+err.Error())
		return
	}
	defer out.Close()

	rotated := rotate270(img)
	switch format {
	case "png":
		err = png.Encode(out, rotated)
	case "gif":
		err = gif.Encode(out, rotated, nil)
	default:
		err = jpeg.Encode(out, rotated, nil)
	}
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, "image save error:"+err.Error())
		return
	}

	utils.HttpSuccess(c, "done")
}

// maps web or absolute URL paths to a path under IMAGE_DIR when possible
func imageStatPathToFilePath(filename string) string {
	if strings.HasPrefix(filename, STAT_IMAGE_PATH+"/") {
		return IMAGE_DIR + strings.TrimPrefix(filename, STAT_IMAGE_PATH)
	}
	if strings.HasPrefix(filename, "/") {
		parts := strings.Split(filename, "/")
		if len(parts) > 2 {
			return IMAGE_DIR + "/" + strings.Join(parts[2:], "/")
		}
	}
	return filename
}

// single 270° rotation (counter-clockwise if y grows downward); used for in-place photo adjust
func rotate270(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, height, width))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(y-bounds.Min.Y, width-1-(x-bounds.Min.X), src.At(x, y))
		}
	}
	return dst
}

// генерим случаный uid для названия файла
func randomFilename() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return strings.ToLower(fmt.Sprintf("%v", ulid.MustNew(ulid.Timestamp(t), entropy)))
}

// picks output filename and opens the file; ulid mode vs sanitized original + O_EXCL collision loop
func createOutputFile(path, filePrefix, imgExt string, preserveOriginalFileName bool, fileHeader *multipart.FileHeader) (string, *os.File, error) {
	if !preserveOriginalFileName || fileHeader == nil {
		fullFileName := fmt.Sprintf("%s%s.%s", filePrefix, randomFilename(), imgExt)
		f, err := os.Create(filepath.Join(path, fullFileName))
		return fullFileName, f, err
	}

	// in preserveOriginalFileName mode:
	// - sanitize the original basename
	// - keep/derive extension safely
	// - avoid overwriting by adding a numeric suffix on collision
	base, extFromName := sanitizeOriginalName(fileHeader.Filename)
	finalExt := imgExt
	if len(extFromName) > 0 {
		finalExt = extFromName
	}
	if len(base) == 0 {
		base = randomFilename()
	}

	const maxAttempts = 1000
	for i := 0; i <= maxAttempts; i++ {
		candidateBase := base
		if i > 0 {
			candidateBase = fmt.Sprintf("%s-%d", base, i)
		}
		fullFileName := fmt.Sprintf("%s%s.%s", filePrefix, candidateBase, finalExt)
		fullPath := filepath.Join(path, fullFileName)

		f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
		if err == nil {
			return fullFileName, f, nil
		}
		if os.IsExist(err) {
			continue
		}
		return "", nil, err
	}

	// fallback: too many collisions (or hostile inputs)
	fullFileName := fmt.Sprintf("%s%s.%s", filePrefix, randomFilename(), imgExt)
	f, err := os.Create(filepath.Join(path, fullFileName))
	return fullFileName, f, err
}

func sanitizeOriginalName(original string) (base string, ext string) {
	// strip any path components and normalize
	name := strings.TrimSpace(filepath.Base(original))
	name = strings.ReplaceAll(name, " ", "_")

	// keep only simple safe characters to avoid filesystem edge cases
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			b.WriteRune(r)
		}
	}
	safe := b.String()
	safe = strings.Trim(safe, "._-")
	if len(safe) == 0 {
		return "", ""
	}

	// cap length to avoid very long filenames
	if len(safe) > 80 {
		safe = safe[:80]
		safe = strings.Trim(safe, "._-")
	}
	if len(safe) == 0 {
		return "", ""
	}

	// split extension
	lastDot := strings.LastIndex(safe, ".")
	if lastDot <= 0 || lastDot == len(safe)-1 {
		return safe, ""
	}
	basePart := safe[:lastDot]
	extPart := strings.ToLower(safe[lastDot+1:])
	if extPart == "png" || extPart == "jpeg" || extPart == "jpg" || extPart == "gif" {
		return basePart, extPart
	}
	return basePart, ""
}
