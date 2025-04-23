package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type StudyRepo struct {
	Id               uuid.UUID   `json:"id"`
	Name             string      `json:"name"`
	UserID           uuid.UUID   `json:"user_id"`
	ImageID          uuid.UUID   `json:"image_id"`
	CreatedAt        time.Time   `json:"created_at"`
	ColoniesNumber   int         `json:"colonies_number"`
	PredictedClasses map[int]int `json:"predicted_classes"`
	IsDeleted        bool        `json:"is_deleted"`
}

func (s *Storage) CreateStudy(tx *StorageTx, userID, imageID string, coloniesNumber int, predictedClasses map[int]int) (StudyRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	classes, err := json.Marshal(predictedClasses)
	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("error serializing predicted classes: %w", err)}
	}
	classesJson := string(classes)

	query, args, err := builder.Insert(`"study" (id, user_id, name,image_id, created_at, colonies_number, predicted_classes)`).
		Values(uuid.New(), userID, "", imageID, time.Now().UTC(), coloniesNumber, classesJson).
		Suffix(`RETURNING id, user_id, name, image_id, created_at, colonies_number, predicted_classes, is_deleted`).
		ToSql()

	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("error creating study: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var study StudyRepo

	err = (*tx).tx.QueryRow(ctx, query, args...).Scan(
		&study.Id,
		&study.UserID,
		&study.Name,
		&study.ImageID,
		&study.CreatedAt,
		&study.ColoniesNumber,
		&study.PredictedClasses,
		&study.IsDeleted,
	)
	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("error creating study: %w", err)}
	}
	return study, nil
}

func (s *Storage) GetStudy(id string) (StudyRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.Select(`
		id,
		user_id,
		name,
		image_id,
		created_at,
		colonies_number,
		predicted_classes,
		is_deleted`).
		From(`study s`).
		Where(
			sq.And{
				sq.Eq{`s.id`: id},
				sq.Eq{`s.is_deleted`: false},
			}).
		ToSql()

	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("Error getting study: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var study StudyRepo
	err = s.db.QueryRow(ctx, query, args...).Scan(
		&study.Id,
		&study.UserID,
		&study.Name,
		&study.ImageID,
		&study.CreatedAt,
		&study.ColoniesNumber,
		&study.PredictedClasses,
		&study.IsDeleted,
	)
	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("Error getting study: %w", err)}
	}
	return study, nil
}

func (s *Storage) DeleteStudies(userID uuid.UUID, id []uuid.UUID) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Update(`study`).
		Set(`is_deleted`, true).
		Where(
			sq.And{
				sq.Eq{`id`: id},
				sq.Eq{`is_deleted`: false},
				sq.Eq{`user_id`: userID},
			}).
		ToSql()

	if err != nil {
		return er.BadRequestErr{Message: fmt.Errorf("Error deleting studies: %w", err)}
	}

	ctx, _ := utils.GetContext()

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return er.BadRequestErr{Message: fmt.Errorf("Error deleting studies: %w", err)}
	}

	return nil
}

func (s *Storage) UpdateStudy(id, name string) (StudyRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.Update(`study`).
		Set(`name`, name).
		Where(
			sq.And{
				sq.Eq{`id`: id},
				sq.Eq{`is_deleted`: false},
			}).
		Suffix(`RETURNING id, user_id, name, image_id, created_at, colonies_number, predicted_classes, is_deleted`).
		ToSql()

	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("error updating study: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var study StudyRepo
	err = s.db.QueryRow(ctx, query, args...).Scan(
		&study.Id,
		&study.UserID,
		&study.Name,
		&study.ImageID,
		&study.CreatedAt,
		&study.ColoniesNumber,
		&study.PredictedClasses,
		&study.IsDeleted,
	)
	if err != nil {
		return StudyRepo{}, er.InternalErr{Message: fmt.Errorf("error updating study: %w", err)}
	}
	return study, nil
}
