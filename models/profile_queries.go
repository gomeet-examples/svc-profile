package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func diffTime(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func ageFull(birthday time.Time) (year, month, day, hour, min, sec int) {
	return diffTime(birthday, time.Now())
}

func age(birthday time.Time) (year int) {
	y, _, _, _, _, _ := ageFull(birthday)
	return y
}

func castProfileParams(gender uint16, email, name, birthday string) (error, *uint16, *string, *string, *time.Time) {
	switch {
	case gender == 0:
		return errors.New("gender is required"), nil, nil, nil, nil
	case gender > 2:
		return fmt.Errorf("%s bad gender value [1|2]", gender), nil, nil, nil, nil
	case len(email) > 100:
		return fmt.Errorf("%s bad email value", email), nil, nil, nil, nil
	case len(name) > 255:
		return fmt.Errorf("%s bad name value", name), nil, nil, nil, nil
		// TODO more business validation rules
	}

	tBirthday, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return fmt.Errorf("%s bad birthday - %s", birthday, err), nil, nil, nil, nil
	}
	tBirthday = tBirthday.UTC()
	a := age(tBirthday)
	if a < 17 || a > 99 {
		return fmt.Errorf("%s bad birthday - 17 < age < 99", birthday), nil, nil, nil, nil
	}

	return nil, &gender, &email, &name, &tBirthday
}

// CreateDevice inserts a new profile in the table.
func CreateProfile(db *gorm.DB, gender uint16, email, name, birthday string) (*Profile, error) {
	err, cGender, cEmail, cName, cBirthday := castProfileParams(gender, email, name, birthday)
	if err != nil {
		return nil, err
	}
	profile := &Profile{
		Gender:   *cGender,
		Email:    *cEmail,
		Name:     *cName,
		Birthday: cBirthday.UTC(),
	}

	err = db.Create(profile).Error

	return profile, err
}

// FindProfileByUUID returns the existing profile with the specified UUID.
func FindProfileByUUID(db *gorm.DB, sUuid string) (*Profile, error) {
	uuid, err := uuid.Parse(sUuid)
	if err != nil {
		return nil, fmt.Errorf("%s bad UUID %s", err)
	}

	profile := &Profile{
		UUID: uuid,
	}

	// Unscoped() enables the retrieval of logically-deleted rows
	if db.Unscoped().First(profile).RecordNotFound() {
		return nil, fmt.Errorf("cannot find profile UUID %s", sUuid)
	}

	return profile, nil
}

// ListProfiles returns a list of profiles matching a set of criteria.
func ListProfiles(db *gorm.DB, offset uint, limit uint, order string, criteria map[string]interface{}, excludeSoftDeleted bool, softDeletedOnly bool) ([]Profile, uint, bool, error) {
	var (
		profiles      []Profile
		resultSetSize int
	)

	// Unscoped() enables the retrieval of logically-deleted rows
	db = db.Unscoped().Limit(limit).Offset(offset).Order(order).Where(criteria)

	if excludeSoftDeleted && softDeletedOnly {
		return profiles, 0, false, errors.New("excludeSoftDeleted and softDeletedOnly are true")
	}

	if excludeSoftDeleted {
		db = db.Where("deleted_at IS NULL")
	} else if softDeletedOnly {
		db = db.Where("deleted_at IS NOT NULL")
	}

	db = db.Find(&profiles)
	err := db.Error

	// the offset must be cancelled to get the total count
	db.Offset(-1).Count(&resultSetSize)
	hasMore := false
	if int(offset)+len(profiles) < resultSetSize {
		hasMore = true
	}

	return profiles, uint(resultSetSize), hasMore, err
}

// UpdateProfile updates the existing profile with the specified UUID using the modifications
// provided in the map argument.
func UpdateProfile(db *gorm.DB, sUuid string, gender uint16, email, name, birthday string) (*Profile, error) {
	uuid, err := uuid.Parse(sUuid)
	if err != nil {
		return nil, fmt.Errorf("%s bad UUID %s", err)
	}

	err, cGender, cEmail, cName, cBirthday := castProfileParams(gender, email, name, birthday)
	if err != nil {
		return nil, err
	}

	// set changes fields
	changes := map[string]interface{}{
		"gender":   *cGender,
		"email":    *cEmail,
		"name":     *cName,
		"birthday": cBirthday.UTC(),
	}

	profile := &Profile{
		UUID: uuid,
	}

	// Unscoped() enables the retrieval of logically-deleted rows
	if db.Unscoped().First(profile).RecordNotFound() {
		return nil, fmt.Errorf("cannot find profile UUID %s", sUuid)
	}

	err = db.Model(profile).Updates(changes).Error

	return profile, err
}

// DeleteProfileLogically performs the logical deletion of the profile with the specified ID.
func DeleteProfileLogically(db *gorm.DB, sUuid string) (*Profile, error) {
	var profile *Profile

	uuid, err := uuid.Parse(sUuid)
	if err != nil {
		return nil, fmt.Errorf("%s bad UUID %s", err)
	}

	profile = &Profile{
		UUID: uuid,
	}

	// Unscoped() enables the retrieval of logically-deleted rows
	if db.First(profile).RecordNotFound() {
		return nil, fmt.Errorf("cannot find profile UUID %s", sUuid)
	}

	err = db.Set("gorm:delete_option", "LIMIT 1").Delete(profile).Error
	if err == nil && profile.DeletedAt == nil {
		t := time.Now()
		profile.DeletedAt = &t
	}

	return profile, err
}

// DeleteProfilePhysically performs the physical deletion of the device with the specified ID.
func DeleteProfilePhysically(db *gorm.DB, sUuid string) error {
	if _, err := uuid.Parse(sUuid); err != nil {
		return fmt.Errorf("%s bad UUID %s", err)
	}

	return db.Exec("DELETE FROM profiles WHERE uuid = ? LIMIT 1", sUuid).Error
}
