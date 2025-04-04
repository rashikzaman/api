package services

import (
	"bytes"
	"context"
	"fmt"
	"rashikzaman/api/models"
	"rashikzaman/api/utils"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/pkg/errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3BucketName = "act-local"
	S3Region     = "ap-southeast-1"
	S3FolderPath = "uploads/images"
)

type Filter struct {
	CategoryIDs        []string
	Skills             []string
	Latitude           float32
	Longitude          float32
	SubscribeUserID    uuid.UUID
	SearchTerm         string
	CreatedByUserID    uuid.UUID
	SubscribedByUserID uuid.UUID
	ApplyBlockFilter   bool
	Distance           int
}

func CreateTask(
	ctx context.Context, db bun.IDB, taskBody *models.Task, userID uuid.UUID,
	awsAccessKey, awsSecretAccessKey string,
) error {
	taskBody.Location = models.PostgisGeometry{Geometry: orb.Point{taskBody.Longitude, taskBody.Latitude}, SRID: 4326}
	taskBody.UserID = userID

	_, err := db.NewInsert().Model(taskBody).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	cfg := aws.Config{
		Region:      "ap-southeast-1",
		Credentials: credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretAccessKey, ""),
	}

	if awsAccessKey != "" && awsSecretAccessKey != "" {
		for _, media := range taskBody.Media {
			decoded, err := utils.DecodeBase64Image(media.Base64)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			mimetype, err := utils.GetMimeTypeFromBase64(media.Base64)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}

			fileExt := ""
			switch mimetype {
			case "image/jpeg":
				fileExt = ".jpg"
			case "image/png":
				fileExt = ".png"
			case "image/gif":
				fileExt = ".gif"
			default:
				fileExt = ".bin"
			}

			filename := fmt.Sprintf("%s/%s%s", S3FolderPath, uuid.New().String(), fileExt)

			client := s3.NewFromConfig(cfg)

			uploadParams := &s3.PutObjectInput{
				Bucket:        aws.String(S3BucketName),
				Key:           aws.String(filename),
				Body:          bytes.NewReader(decoded),
				ContentType:   aws.String(mimetype),
				ContentLength: aws.Int64(int64(len(decoded))),
			}

			// Upload to S3
			_, err = client.PutObject(ctx, uploadParams)
			if err != nil {
				fmt.Printf("Error uploading to S3: %v\n", err)
				return errors.Wrap(err, "failed to upload file to S3")
			}

			// Generate the S3 URL for the uploaded file
			s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3BucketName, S3Region, filename)

			newMedia := &models.TaskMedia{}
			newMedia.Link = s3URL
			newMedia.MimeType = mimetype
			newMedia.TaskID = taskBody.ID

			// Save media info to database
			_, err = db.NewInsert().Model(newMedia).Exec(ctx)
			if err != nil {
				return errors.Wrap(err, err.Error())
			}

		}
	}

	return nil
}

func FetchTasks(
	ctx context.Context, db bun.IDB, queryParam models.QueryParam,
	filter Filter,
) ([]models.Task, int, error) {
	tasks := []models.Task{}

	query := db.NewSelect().
		Model(&tasks)

	query.WhereGroup("AND", func(sq *bun.SelectQuery) *bun.SelectQuery {
		if len(filter.CategoryIDs) > 0 {
			sq.Where("category_id IN (?)", bun.In(filter.CategoryIDs))
		}

		if len(filter.Skills) > 0 {
			sq.WhereOr("required_skills && ?", pgdialect.Array(filter.Skills))
		}

		return sq
	})

	// find posts within specified meter's radius
	if filter.Latitude != 0 && filter.Longitude != 0 {
		query.Where(
			"ST_DWithin(ST_MakePoint(?, ?)::geography, location::geography, ?)",
			filter.Longitude, filter.Latitude, filter.Distance*1000,
		)
	}

	//check if the current user is subscribed to this Task
	if filter.SubscribeUserID != uuid.Nil {
		query.Column("task.*")
		query.ColumnExpr(
			`EXISTS (
			SELECT 1 FROM user_tasks ut
			WHERE ut.task_id = task.id AND ut.user_id = ?
		) AS is_subscribed`, filter.SubscribeUserID,
		)
	}

	if filter.CreatedByUserID != uuid.Nil {
		query.Where("user_id = ?", filter.CreatedByUserID)
	}

	if filter.SubscribedByUserID != uuid.Nil {
		fmt.Println(filter.SubscribedByUserID)
		query.Join("INNER JOIN user_tasks ON user_tasks.task_id = task.id").
			Where("user_tasks.user_id = ?", filter.SubscribedByUserID)
	}

	if filter.ApplyBlockFilter {
		query.Where("task.blocked = FALSE")
	}

	//case-insensitive search
	if filter.SearchTerm != "" {
		query.Where("title ILIKE ?", "%"+filter.SearchTerm+"%")
	}

	count, err := queryParam.Pagination.BuildPaginationQuery(ctx, query)
	if err != nil {
		return tasks, 0, errors.Wrap(err, err.Error())
	}

	if len(queryParam.Relations) != 0 {
		for _, relation := range queryParam.Relations {
			query.Relation(relation)
		}
	}

	err = query.Scan(ctx)
	if err != nil {
		return tasks, 0, errors.Wrap(err, err.Error())
	}

	return tasks, count, nil
}

func FetchTaskByID(ctx context.Context, db bun.IDB, id uuid.UUID, queryParam models.QueryParam) (models.Task, error) {
	task := models.Task{}

	err := models.SelectByID(ctx, db, id, &task, queryParam)
	if err != nil {
		return task, errors.Wrap(err, err.Error())
	}

	return task, nil
}

func UpdateTask(
	ctx context.Context, db bun.IDB, existingTask models.Task, taskBody models.Task,
) (models.Task, error) {
	existingTask.Title = taskBody.Title
	existingTask.Description = taskBody.Description
	existingTask.Latitude = taskBody.Latitude
	existingTask.Longitude = taskBody.Longitude
	existingTask.Location = models.PostgisGeometry{
		Geometry: orb.Point{taskBody.Longitude, taskBody.Latitude}, SRID: 4326,
	}
	existingTask.FormattedAddress = taskBody.FormattedAddress
	existingTask.CategoryID = taskBody.CategoryID
	existingTask.RequiredSkills = taskBody.RequiredSkills
	existingTask.RequiredVolunteersCount = taskBody.RequiredVolunteersCount

	err := models.Update(ctx, db, &existingTask)
	if err != nil {
		return existingTask, errors.Wrap(err, err.Error())
	}

	return existingTask, nil
}

func DeleteTask(ctx context.Context, db bun.IDB, taskID uuid.UUID) error {
	task := models.Task{}
	task.ID = taskID

	err := models.Delete(ctx, db, &task)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func ApplyToTask(ctx context.Context, db bun.IDB, taskID, userID uuid.UUID) error {
	userTask := &models.UserTask{UserID: userID, TaskID: taskID}

	return models.Create(ctx, db, userTask)
}

func WithdrawFromTask(ctx context.Context, db bun.IDB, taskID, userID uuid.UUID) error {
	userTask := &models.UserTask{}

	_, err := db.NewDelete().Model(userTask).Where("user_id = ?", userID).Where("task_id = ?", taskID).Exec(ctx)

	return err
}

func FetchSubscribersForTask(ctx context.Context, db bun.IDB, taskID uuid.UUID) ([]models.UserTask, error) {
	var userTasks []models.UserTask

	query := db.NewSelect().
		Model(&userTasks).
		Join("JOIN tasks ON tasks.id = user_task.task_id").
		Where("user_task.task_id = ?", taskID).
		Relation("User").
		Relation("Task")

	err := query.Scan(ctx)

	return userTasks, err
}

func ApplyActionToTask(ctx context.Context, db bun.IDB, taskID uuid.UUID, action string) (*models.Task, error) {
	task, err := FetchTaskByID(ctx, db, taskID, models.QueryParam{})
	if err != nil {
		return nil, err
	}

	if action == "block" {
		task.Blocked = true
	} else if action == "unblock" {
		task.Blocked = false
	}

	err = models.Update(ctx, db, &task)

	return &task, err
}

func FetchNearbyUsersOfTask(ctx context.Context, db bun.IDB, taskID uuid.UUID, latitude, longitude float32, distance int) ([]models.UserLocation, error) {
	userLocations := []models.UserLocation{}
	query := db.NewSelect().Model(&userLocations)
	query.Where(
		"ST_DWithin(ST_MakePoint(?, ?)::geography, location::geography, ?)",
		longitude, latitude, distance*1000,
	)

	query.Relation("User")

	err := query.Scan(ctx)

	return userLocations, err
}
