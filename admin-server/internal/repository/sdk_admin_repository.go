package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"postapocgame/admin-server/internal/model"
)

type SdkAdminRepository struct {
	repo *Repository
}

func NewSdkAdminRepository(repo *Repository) *SdkAdminRepository {
	return &SdkAdminRepository{repo: repo}
}

// -------- API Key --------

func (r *SdkAdminRepository) FindSdkKey(ctx context.Context, id uint64) (*model.SdkKey, error) {
	return r.repo.SdkKeyModel.FindOne(ctx, id)
}

func (r *SdkAdminRepository) FindSdkKeyByApiKey(ctx context.Context, apiKey string) (*model.SdkKey, error) {
	return r.repo.SdkKeyModel.FindOneByApiKey(ctx, apiKey)
}

func (r *SdkAdminRepository) FindSdkKeyByApiSecret(ctx context.Context, apiSecret string) (*model.SdkKey, error) {
	return r.repo.SdkKeyModel.FindOneByApiSecret(ctx, apiSecret)
}

func (r *SdkAdminRepository) ListSdkKeys(ctx context.Context, page, pageSize int64, name string, status int64) ([]model.SdkKey, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	where := []string{"deleted_at = 0"}
	args := []interface{}{}

	if name != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if status != 0 {
		where = append(where, "status = ?")
		args = append(args, status)
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM sdk_key WHERE %s", whereClause)
	if err := r.repo.DB.QueryRowCtx(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("SELECT * FROM sdk_key WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	argsWithPage := append(args, pageSize, offset)
	var list []model.SdkKey
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, listSQL, argsWithPage...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *SdkAdminRepository) CreateSdkKey(ctx context.Context, key *model.SdkKey) (uint64, error) {
	res, err := r.repo.SdkKeyModel.Insert(ctx, key)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *SdkAdminRepository) UpdateSdkKey(ctx context.Context, key *model.SdkKey) error {
	return r.repo.SdkKeyModel.Update(ctx, key)
}

func (r *SdkAdminRepository) DeleteSdkKey(ctx context.Context, id uint64) error {
	return r.repo.SdkKeyModel.Delete(ctx, id)
}

// -------- SDK Interface --------

func (r *SdkAdminRepository) FindInterface(ctx context.Context, id uint64) (*model.SdkInterface, error) {
	return r.repo.SdkInterfaceModel.FindOne(ctx, id)
}

func (r *SdkAdminRepository) FindInterfaceByCode(ctx context.Context, apiCode string) (*model.SdkInterface, error) {
	return r.repo.SdkInterfaceModel.FindOneByApiCode(ctx, apiCode)
}

func (r *SdkAdminRepository) ListInterfaces(ctx context.Context, page, pageSize int64, name, apiCode string, status int64) ([]model.SdkInterface, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	where := []string{"deleted_at = 0"}
	args := []interface{}{}

	if name != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if apiCode != "" {
		where = append(where, "api_code LIKE ?")
		args = append(args, "%"+apiCode+"%")
	}
	if status != 0 {
		where = append(where, "status = ?")
		args = append(args, status)
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM sdk_interface WHERE %s", whereClause)
	if err := r.repo.DB.QueryRowCtx(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("SELECT * FROM sdk_interface WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	argsWithPage := append(args, pageSize, offset)
	var list []model.SdkInterface
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, listSQL, argsWithPage...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *SdkAdminRepository) CreateInterface(ctx context.Context, iface *model.SdkInterface) (uint64, error) {
	res, err := r.repo.SdkInterfaceModel.Insert(ctx, iface)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *SdkAdminRepository) UpdateInterface(ctx context.Context, iface *model.SdkInterface) error {
	return r.repo.SdkInterfaceModel.Update(ctx, iface)
}

func (r *SdkAdminRepository) DeleteInterface(ctx context.Context, id uint64) error {
	return r.repo.SdkInterfaceModel.Delete(ctx, id)
}

// -------- 绑定关系 --------

// ListBindings 返回接口列表及绑定状态
func (r *SdkAdminRepository) ListBindings(ctx context.Context, sdkKeyId uint64) ([]SdkBindingView, error) {
	args := []interface{}{sdkKeyId}
	query := `
SELECT
    i.id AS sdk_interface_id,
    i.api_code,
    i.name,
    i.path,
    i.method,
    i.rate_limit_default,
    IFNULL(k.id, 0) AS bound,
    IFNULL(k.custom_rate_limit, 0) AS custom_rate_limit
FROM sdk_interface i
LEFT JOIN sdk_key_api k ON k.sdk_interface_id = i.id AND k.sdk_key_id = ? AND k.deleted_at = 0
WHERE i.deleted_at = 0
ORDER BY i.id DESC`
	var list []SdkBindingView
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

type SdkBindingView struct {
	SdkInterfaceId  uint64 `db:"sdk_interface_id"`
	ApiCode         string `db:"api_code"`
	Name            string `db:"name"`
	Path            string `db:"path"`
	Method          string `db:"method"`
	RateLimit       int64  `db:"rate_limit_default"`
	Bound           int64  `db:"bound"`
	CustomRateLimit int64  `db:"custom_rate_limit"`
}

// SaveBindings 先软删除旧绑定，再插入新绑定
func (r *SdkAdminRepository) SaveBindings(ctx context.Context, sdkKeyId uint64, bindings []model.SdkKeyApi) error {
	now := time.Now().Unix()
	// 软删除旧绑定
	delSQL := "UPDATE sdk_key_api SET deleted_at = ?, updated_at = ? WHERE sdk_key_id = ? AND deleted_at = 0"
	if _, err := r.repo.DB.ExecCtx(ctx, delSQL, now, now, sdkKeyId); err != nil {
		return err
	}

	for _, b := range bindings {
		// 保底 updated_at/created_at
		if b.CreatedAt == 0 {
			b.CreatedAt = now
		}
		if b.UpdatedAt == 0 {
			b.UpdatedAt = now
		}
		if b.DeletedAt != 0 {
			b.DeletedAt = 0
		}
		_, err := r.repo.SdkKeyApiModel.Insert(ctx, &b)
		if err != nil {
			return err
		}
	}
	return nil
}

// -------- 调用记录 --------

func (r *SdkAdminRepository) ListCallLogs(ctx context.Context, page, pageSize int64, sdkKeyId uint64, apiCode string, respCode int64, ip string, startTime, endTime int64) ([]model.SdkCallLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	where := []string{"deleted_at = 0"}
	args := []interface{}{}

	if sdkKeyId > 0 {
		where = append(where, "sdk_key_id = ?")
		args = append(args, sdkKeyId)
	}
	if apiCode != "" {
		where = append(where, "api_code LIKE ?")
		args = append(args, "%"+apiCode+"%")
	}
	if respCode != 0 {
		where = append(where, "resp_code = ?")
		args = append(args, respCode)
	}
	if ip != "" {
		where = append(where, "ip LIKE ?")
		args = append(args, "%"+ip+"%")
	}
	if startTime > 0 {
		where = append(where, "created_at >= ?")
		args = append(args, startTime)
	}
	if endTime > 0 {
		where = append(where, "created_at <= ?")
		args = append(args, endTime)
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM sdk_call_log WHERE %s", whereClause)
	if err := r.repo.DB.QueryRowCtx(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("SELECT * FROM sdk_call_log WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	argsWithPage := append(args, pageSize, offset)
	var list []model.SdkCallLog
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, listSQL, argsWithPage...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// ExportCallLogs 返回最多 maxRows 行
func (r *SdkAdminRepository) ExportCallLogs(ctx context.Context, maxRows int64, sdkKeyId uint64, apiCode string, respCode int64, ip string, startTime, endTime int64) ([]model.SdkCallLog, error) {
	if maxRows <= 0 {
		maxRows = 2000
	}
	where := []string{"deleted_at = 0"}
	args := []interface{}{}

	if sdkKeyId > 0 {
		where = append(where, "sdk_key_id = ?")
		args = append(args, sdkKeyId)
	}
	if apiCode != "" {
		where = append(where, "api_code LIKE ?")
		args = append(args, "%"+apiCode+"%")
	}
	if respCode != 0 {
		where = append(where, "resp_code = ?")
		args = append(args, respCode)
	}
	if ip != "" {
		where = append(where, "ip LIKE ?")
		args = append(args, "%"+ip+"%")
	}
	if startTime > 0 {
		where = append(where, "created_at >= ?")
		args = append(args, startTime)
	}
	if endTime > 0 {
		where = append(where, "created_at <= ?")
		args = append(args, endTime)
	}

	whereClause := strings.Join(where, " AND ")
	query := fmt.Sprintf("SELECT * FROM sdk_call_log WHERE %s ORDER BY id DESC LIMIT ?", whereClause)
	args = append(args, maxRows)

	var list []model.SdkCallLog
	if err := r.repo.DB.QueryRowsCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

// -------- 字典 --------

func (r *SdkAdminRepository) GetRateLimitDefault(ctx context.Context) (int64, error) {
	// 从字典 sdk_rate_limit_default 读取第一条值
	typeId, err := r.repo.AdminDictTypeModel.FindIdByCode(ctx, "sdk_rate_limit_default")
	if err != nil {
		return 0, err
	}
	if typeId == 0 {
		return 0, sql.ErrNoRows
	}
	items, _, err := r.repo.AdminDictItemModel.FindPageByTypeId(ctx, typeId, 1, 1)
	if err != nil {
		return 0, err
	}
	if len(items) == 0 {
		return 0, sql.ErrNoRows
	}
	v, err := strconv.ParseInt(items[0].Value, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
