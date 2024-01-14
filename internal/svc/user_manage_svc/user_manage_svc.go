package user_manage_svc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kilicmu/user-service/internal/database"
	"github.com/kilicmu/user-service/internal/utils"
)

const (
	USER_INFO_BY_UID = "USER_INFO_BY_UID"
	USER_ID_SEED     = 100000
)

func getUserInfoCacheKey(uid string) string {
	return fmt.Sprintf("USER_INFO_BY_UID/%s", uid)
}

func getCachedByKey(key string, value any) (bool, error) {
	rdb := database.GetRedisUtils()
	if uInfoString, err := rdb.Get(key); err == nil {
		if err := json.Unmarshal([]byte(uInfoString), value); err == nil {
			return uInfoString != "", nil
		}
	}
	return false, nil
}

func setCacheByKV(key string, value any, expr int32) (bool, error) {
	rdb := database.GetRedisUtils()
	if jsonBytes, err := json.Marshal(value); err == nil {
		rdb.SetEx(key, string(jsonBytes), 2*60*60)
		return true, nil
	}
	return false, fmt.Errorf("set cached error")
}

func deleteCacheByKey(key ...string) (bool, error) {
	rdb := database.GetRedisUtils()
	if _, err := rdb.Delete(key...); err == nil {
		return true, nil
	}
	return false, fmt.Errorf("delete cached error")
}

func (s *UserManageService) SetCache(info *database.UserInfo) {
	if info.UId != "" {
		setCacheByKV(getUserInfoCacheKey("uid=?"+info.UId), info, 2*60*60)
	}
	if info.Email != "" {
		setCacheByKV(getUserInfoCacheKey("email=?"+info.Email), info, 2*60*60)
	}
	if info.Name != "" {
		setCacheByKV(getUserInfoCacheKey("name=?"+info.Name), info, 2*60*60)
	}
}

func (s *UserManageService) DeleteCache(info *database.UserInfo) {
	cachekeys := make([]string, 0)
	if info.UId != "" {
		cachekeys = append(cachekeys, getUserInfoCacheKey("uid=?"+info.UId))
	}
	if info.Email != "" {
		cachekeys = append(cachekeys, getUserInfoCacheKey("email=?"+info.Email))
	}
	if info.Name != "" {
		cachekeys = append(cachekeys, getUserInfoCacheKey("name=?"+info.Name))
	}
	deleteCacheByKey(cachekeys...)
}

func (s *UserManageService) EncodePassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s *UserManageService) CreateUserIfNotExist(info *database.UserInfo) (*database.UserInfo, error) {
	db := database.GetDB()
	tx := db.Begin()
	if info.UId == "" {
		var totalCountOfUser int64 = 0
		tx.Model(info).Count(&totalCountOfUser)
		info.UId = fmt.Sprintf("%d", totalCountOfUser+USER_ID_SEED)
	}

	userinfo := &database.UserInfo{}
	res := tx.Where("email=?", info.Email).Take(userinfo)
	if res.RowsAffected != 0 {
		tx.Rollback()
		return userinfo, fmt.Errorf("email already exist")
	}

	result := tx.Where("uid=?", info.UId).Attrs(&info).FirstOrCreate(info)

	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create or query user error: %s", err)
	}
	tx.Commit()
	return info, nil
}

func (s *UserManageService) GetUserByUid(uid string) (*database.UserInfo, error) {
	db := database.GetDB()
	queryString := "uid=?"
	cacheId := getUserInfoCacheKey(queryString + uid)

	userinfo := &database.UserInfo{}
	if matched, err := getCachedByKey(cacheId, userinfo); matched && err == nil {
		return userinfo, nil
	}

	result := db.Where(queryString, uid).First(userinfo)

	if err := result.Error; err != nil || result.RowsAffected == 0 {
		return &database.UserInfo{}, fmt.Errorf("get user by uid error: %s", err)
	}
	s.SetCache(userinfo)
	return userinfo, nil
}

func (s *UserManageService) GetUserInfoByEmailOrPhoneOrName(emailOrPhoneOrName string) (*database.UserInfo, error) {
	db := database.GetDB()
	userinfo := &database.UserInfo{}
	var queryString string
	switch {
	case utils.EmailRegex.MatchString(emailOrPhoneOrName):
		queryString = "email=?"
		break
	case utils.PhoneRegex.MatchString(emailOrPhoneOrName):
		queryString = "phone=?"
		break
	default:
		queryString = "name=?"
	}

	cacheId := getUserInfoCacheKey(queryString + emailOrPhoneOrName)
	if matched, err := getCachedByKey(cacheId, userinfo); matched && err == nil {
		return userinfo, nil
	}

	result := db.Where(queryString, emailOrPhoneOrName).First(userinfo)
	if err := result.Error; err != nil || result.RowsAffected == 0 {
		return &database.UserInfo{}, fmt.Errorf("get user error: %s", err)
	}
	s.SetCache(userinfo)
	return userinfo, nil
}

func (s *UserManageService) ValidateUserPassword(rawPassword string, userinfo *database.UserInfo) bool {
	if rawPassword == "" || userinfo.Password == "" {
		return false
	}

	return s.EncodePassword(rawPassword) == userinfo.Password
}

func (s *UserManageService) UpdateUserInfo(userinfo *database.UserInfo) (*database.UserInfo, error) {
	db := database.GetDB()

	if userinfo.Password != "" {
		userinfo.Password = s.EncodePassword(userinfo.Password)
	}

	s.DeleteCache(userinfo)
	result := db.UpdateColumns(userinfo)
	if err := result.Error; err != nil || result.RowsAffected == 0 {
		return &database.UserInfo{}, fmt.Errorf("get user by uid error: %s", err)
	}
	uinfo, err := s.GetUserByUid(userinfo.UId)
	go func() {
		time.Sleep(1 * time.Second)
		s.SetCache(uinfo)
	}()
	return uinfo, err
}

type UserManageService struct {
}

func NewUserManageService() *UserManageService {
	return &UserManageService{}
}
