package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/doopush/doopush/api/internal/database"
	"github.com/doopush/doopush/api/internal/models"
	"github.com/doopush/doopush/api/pkg/utils"
	"gorm.io/gorm"
)

var ErrAppSecretNotFound = errors.New("App Secret不存在")

type CredentialService struct{}

func NewCredentialService() *CredentialService { return &CredentialService{} }

func normalizeScopes(scopes []string) (models.StringList, error) {
	valid := make(map[string]struct{}, len(models.ValidAppSecretScopes))
	for _, scope := range models.ValidAppSecretScopes {
		valid[scope] = struct{}{}
	}
	unique := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		if _, ok := valid[scope]; !ok {
			return nil, fmt.Errorf("无效的权限范围: %s", scope)
		}
		unique[scope] = struct{}{}
	}
	if len(unique) == 0 {
		return nil, errors.New("请至少选择一个权限范围")
	}
	result := make(models.StringList, 0, len(unique))
	for scope := range unique {
		result = append(result, scope)
	}
	sort.Strings(result)
	return result, nil
}

func (s *CredentialService) ListAppSecrets(appID uint) ([]models.AppSecret, error) {
	var secrets []models.AppSecret
	err := database.DB.Where("app_id = ?", appID).Order("created_at DESC").Find(&secrets).Error
	return secrets, err
}

func (s *CredentialService) CreateAppSecret(appID, userID uint, name string, scopes []string, expiresAt *time.Time) (*models.AppSecret, string, error) {
	normalizedScopes, err := normalizeScopes(scopes)
	if err != nil {
		return nil, "", err
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		return nil, "", errors.New("过期时间必须晚于当前时间")
	}
	var activeCount int64
	if err := database.DB.Model(&models.AppSecret{}).
		Where("app_id = ? AND status = 1 AND revoked_at IS NULL", appID).Count(&activeCount).Error; err != nil {
		return nil, "", errors.New("检查App Secret数量失败")
	}
	if activeCount >= 20 {
		return nil, "", errors.New("每个应用最多保留20个有效App Secret")
	}

	plain := utils.GenerateSecureToken(models.AppSecretPrefix)
	secret := &models.AppSecret{
		AppID: appID, Name: name, SecretHash: utils.HashCredential(plain),
		Prefix: models.AppSecretPrefix, Suffix: plain[len(plain)-8:], Scopes: normalizedScopes,
		Status: 1, ExpiresAt: expiresAt, CreatedBy: userID,
	}
	if err := database.DB.Create(secret).Error; err != nil {
		return nil, "", errors.New("App Secret创建失败")
	}
	return secret, plain, nil
}

func (s *CredentialService) UpdateAppSecretScopes(appID, secretID uint, scopes []string) (*models.AppSecret, models.StringList, error) {
	normalizedScopes, err := normalizeScopes(scopes)
	if err != nil {
		return nil, nil, err
	}

	var secret models.AppSecret
	if err := database.DB.Where("id = ? AND app_id = ?", secretID, appID).First(&secret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrAppSecretNotFound
		}
		return nil, nil, errors.New("查询App Secret失败")
	}
	if secret.Status != 1 || secret.RevokedAt != nil {
		return nil, nil, errors.New("已撤销的App Secret不能修改权限")
	}
	if secret.ExpiresAt != nil && !secret.ExpiresAt.After(time.Now()) {
		return nil, nil, errors.New("已过期的App Secret不能修改权限")
	}

	previousScopes := append(models.StringList(nil), secret.Scopes...)
	result := database.DB.Model(&models.AppSecret{}).
		Where("id = ? AND app_id = ? AND status = 1 AND revoked_at IS NULL", secretID, appID).
		Update("scopes", normalizedScopes)
	if result.Error != nil {
		return nil, nil, errors.New("更新App Secret权限失败")
	}
	if result.RowsAffected == 0 {
		return nil, nil, errors.New("App Secret状态已变化，请刷新后重试")
	}
	secret.Scopes = normalizedScopes
	secret.UpdatedAt = time.Now()
	return &secret, previousScopes, nil
}

func (s *CredentialService) RevokeAppSecret(appID, secretID uint) error {
	now := time.Now()
	result := database.DB.Model(&models.AppSecret{}).
		Where("id = ? AND app_id = ? AND status = 1", secretID, appID).
		Updates(map[string]interface{}{"status": 0, "revoked_at": now})
	if result.Error != nil {
		return errors.New("撤销App Secret失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("App Secret不存在或已撤销")
	}
	return nil
}

func (s *CredentialService) IssueInstallationCredential(appID, deviceID uint) (string, string, error) {
	token := utils.GenerateSecureToken(models.InstallationTokenPrefix)
	credential := models.InstallationCredential{}
	err := database.DB.Unscoped().Where("device_id = ?", deviceID).First(&credential).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		credential = models.InstallationCredential{AppID: appID, DeviceID: deviceID}
	} else if err != nil {
		return "", "", errors.New("查询安装凭证失败")
	}
	credential.AppID = appID
	credential.TokenHash = utils.HashCredential(token)
	credential.Status = 1
	credential.ExpiresAt = nil
	credential.DeletedAt = gorm.DeletedAt{}
	if credential.ID == 0 {
		err = database.DB.Create(&credential).Error
	} else {
		err = database.DB.Save(&credential).Error
	}
	if err != nil {
		return "", "", errors.New("安装凭证保存失败")
	}
	return fmt.Sprintf("ins_%d", credential.ID), token, nil
}
