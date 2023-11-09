package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/logger"
)

var (
	zhT  = zh.New()
	enT  = en.New()
	uni  = ut.New(enT, enT, zhT)
	locs = []string{"en", "zh"}
)

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		logger.L().Fatal("cannot init validate")
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	trans, _ := uni.GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		logger.L().Fatal("cannot init validate", zap.String("translator", "en"))
	}
	trans, _ = uni.GetTranslator("zh")
	if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		logger.L().Fatal("cannot init validate", zap.String("translator", "en"))
	}

}

type listData[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

type response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type responseWithoutData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ctxJSONWithData
func ctxJSONWithData[T any](ctx *gin.Context, data T) {
	ctx.JSON(http.StatusOK, &response[T]{
		Code: ERRNO_OK,
		Msg:  errno2str[ERRNO_OK],
		Data: data,
	})
}

// ctxJSONWithListData
func ctxJSONWithListData[T any](ctx *gin.Context, total int64, list []T) {
	ctxJSONWithData(ctx, &listData[T]{
		Total: total,
		List:  list,
	})
}

// ctxJSONWithoutData
func ctxJSONWithoutData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &responseWithoutData{
		Code: ERRNO_OK,
		Msg:  errno2str[ERRNO_OK],
	})
}

// ctxJSONWithoutError
func ctxJSONWithError(ctx *gin.Context, code int, err error) {
	ctx.JSON(http.StatusOK, &responseWithoutData{
		Code: code,
		Msg: lo.TernaryF(err == nil,
			func() string { return errno2str[code] },
			func() (msg string) {
				msg = err.Error()
				switch e := err.(type) {
				case validator.ValidationErrors:
					trans, _ := uni.GetTranslator(GetLocFromAcceptLanguage(ctx))
					msg = getMsg(e.Translate(trans))
				}
				return
			},
		),
	})
}

func GetLocFromAcceptLanguage(ctx *gin.Context) string {
	al := strings.ToLower(ctx.GetHeader("Accept-Language"))
	als := strings.Split(al, ",")
	found, ok := lo.Find(locs, func(loc string) bool {
		for _, l := range als {
			if strings.Contains(l, loc) {
				return true
			}
		}
		return false
	})

	return lo.Ternary(ok, found, "en")
}

func getMsg(fields map[string]string) string {
	data := []string{}
	for field, err := range fields {
		data = append(data, fmt.Sprintf("%s: %s", field[strings.Index(field, ".")+1:], err))
	}

	return strings.Join(data, "\n")
}
