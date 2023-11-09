package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/rickywei/sparrow/project/cache"
	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/dao"
	"github.com/rickywei/sparrow/project/po"
	"github.com/rickywei/sparrow/project/vo"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Create
//
//	@Tags			user
//	@Description	create
//	@Accept			json
//	@Produce		json
//	@Param			body	body		vo.ReqCreateUser	true	" "
//	@Success		200		{object}	responseWithoutData
//	@Failure		200		{object}	responseWithoutData
//	@Router			/api/v1/user [POST]
func (h *UserHandler) Create(ctx *gin.Context) {
	req := &vo.ReqCreateUser{}
	if err := ctx.BindJSON(req); err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	data := &po.User{}
	if err := copier.Copy(data, req); err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	if err := dao.Q.User.WithContext(ctx).
		Create(data); err != nil {
		ctx.JSON(
			lo.Ternary(errors.Is(err, gorm.ErrDuplicatedKey), http.StatusBadRequest, http.StatusInternalServerError),
			err,
		)
		return
	}

	ctxJSONWithoutData(ctx)
}

// Delete
//
//	@Tags			user
//	@Description	delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	" "
//	@Success		200	{object}	responseWithoutData
//	@Failure		200	{object}	responseWithoutData
//	@Router			/api/v1/user/:id [DELETE]
func (h *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}

	if err = dao.Q.User.WithContext(ctx).
		DeleteByID(id); err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	if _, err = cache.RC.Del(ctx, fmt.Sprintf("user-%d", id)).Result(); err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	ctxJSONWithoutData(ctx)
}

// Update
//
//	@Tags			user
//	@Description	update
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64				true	" "
//	@Param			body	body		vo.ReqCreateUser	true	" "
//	@Success		200		{object}	responseWithoutData
//	@Failure		200		{object}	responseWithoutData
//	@Router			/api/v1/user/:id [PUT]
func (h *UserHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	data := &po.User{}
	if err := ctx.BindJSON(data); err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	data.ID = id

	if _, err = dao.Q.User.WithContext(ctx).
		Updates(data); err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	if _, err = cache.RC.Del(ctx, fmt.Sprintf("user-%d", id)).Result(); err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	ctxJSONWithoutData(ctx)
}

// Query
//
//	@Tags			user
//	@Description	query
//	@Produce		json
//	@Param			id	path		int64	true	" "
//	@Success		200	{object}	response[po.User]
//	@Failure		200	{object}	responseWithoutData
//	@Router			/api/v1/user/:id [GET]
func (h *UserHandler) Query(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	data, err := dao.Q.User.WithContext(ctx).
		FindByID(id)
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}

	if bs, err := json.Marshal(data); err == nil {
		cache.RC.Del(ctx, fmt.Sprintf("user-%d", bs))
	}

	ctxJSONWithData(ctx, data)
}

// QueryList
//
//	@Tags			user
//	@Description	query list
//	@Produce		json
//	@Param			offset	query		int	true	" "
//	@Param			limit	query		int	true	" "
//	@Success		200		{object}	response[[]po.User]
//	@Failure		200		{object}	responseWithoutData
//	@Router			/api/v1/user/:id [GET]
func (h *UserHandler) QueryList(ctx *gin.Context) {
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}

	data, total, err := dao.Q.User.WithContext(ctx).
		FindByPage(offset, limit)
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	ctxJSONWithListData(ctx, total, data)
}

// Login
//
//	@Tags			user
//	@Description	login
//	@Accept			json
//	@Produce		json
//	@Param			body	body		vo.ReqLogin	true	" "
//	@Success		200		{object}	responseWithoutData
//	@Failure		200		{object}	responseWithoutData
//	@Router			/api/v1/user/login [POST]
func (h *UserHandler) Login(ctx *gin.Context) {
	data := &vo.ReqLogin{}
	if err := ctx.BindJSON(data); err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}

	q := dao.Q.User
	u, err := dao.Q.User.WithContext(ctx).
		Where(q.UserName.Eq(data.UserName), q.Password.Eq(data.Password)).
		First()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if setJwtCookie(ctx, u.ID) != nil {
		return
	}

	ctxJSONWithoutData(ctx)
}

// Refresh
//
//	@Tags			user
//	@Description	refresh
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responseWithoutData
//	@Failure		200	{object}	responseWithoutData
//	@Router			/api/v1/user/refresh [POST]
func (h *UserHandler) Refresh(ctx *gin.Context) {
	uid, _ := ctx.Value("uid").(int64)
	if _, err := dao.Q.User.WithContext(ctx).
		FindByID(uid); err != nil {
		ctxJSONWithError(ctx, ERRNO_INVALID_ARGUMENT, err)
		return
	}

	if setJwtCookie(ctx, uid) != nil {
		return
	}

	ctxJSONWithoutData(ctx)
}

// setJwtCookie sets jwt token in cookie
func setJwtCookie(ctx *gin.Context, uid int64) (err error) {
	exp := conf.Int("jwt.exp")
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject: strconv.FormatInt(uid, 10),
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Second * time.Duration(exp)),
			},
			NotBefore: &jwt.NumericDate{
				Time: now,
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
		},
	)
	token, err := t.SignedString([]byte(conf.String("jwt.secret")))
	if err != nil {
		ctxJSONWithError(ctx, ERRNO_INTERNAL_ERROR, err)
		return
	}

	ctx.SetCookie("token", token, exp, "/", "localhost", false, true)

	return
}
