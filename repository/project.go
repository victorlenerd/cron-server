package repository

import (
	"encoding/json"
	_ "errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/hashicorp/raft"
	"log"
	"net/http"
	"scheduler0/constants"
	"scheduler0/fsm"
	"scheduler0/marsher"
	"scheduler0/models"
	"scheduler0/protobuffs"
	"scheduler0/utils"
	"strconv"
	"time"
)

type Project interface {
	CreateOne(project models.ProjectModel) (int64, *utils.GenericError)
	GetOneByName(project *models.ProjectModel) *utils.GenericError
	GetOneByID(project *models.ProjectModel) *utils.GenericError
	List(offset int64, limit int64) ([]models.ProjectModel, *utils.GenericError)
	Count() (int64, *utils.GenericError)
	UpdateOneByID(project models.ProjectModel) (int64, *utils.GenericError)
	DeleteOneByID(project models.ProjectModel) (int64, *utils.GenericError)
}

type projectRepo struct {
	store   *fsm.Store
	jobRepo Job
	logger  *log.Logger
}

const (
	ProjectsTableName         = "projects"
	ProjectsIdColumn          = "id"
	NameColumn                = "name"
	DescriptionColumn         = "description"
	ProjectsDateCreatedColumn = "date_created"
)

func NewProjectRepo(logger *log.Logger, store *fsm.Store, jobRepo Job) Project {
	return &projectRepo{
		store:   store,
		jobRepo: jobRepo,
		logger:  logger,
	}
}

// CreateOne creates a single project
func (projectRepo *projectRepo) CreateOne(project models.ProjectModel) (int64, *utils.GenericError) {
	if len(project.Name) < 1 {
		return -1, utils.HTTPGenericError(http.StatusBadRequest, "name field is required")
	}

	if len(project.Description) < 1 {
		return -1, utils.HTTPGenericError(http.StatusBadRequest, "description field is required")
	}

	projectWithName := models.ProjectModel{
		ID:   -1,
		Name: project.Name,
	}

	_ = projectRepo.GetOneByName(&project)
	if projectWithName.ID != -1 {
		return -1, utils.HTTPGenericError(http.StatusBadRequest, "Another project exist with the same name")
	}

	query, params, err := sq.Insert(ProjectsTableName).
		Columns(
			NameColumn,
			DescriptionColumn,
			ProjectsDateCreatedColumn,
		).
		Values(
			project.Name,
			project.Description,
			time.Now().String(),
		).ToSql()
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	res, applyErr := projectRepo.applyToFSM(query, params)
	if applyErr != nil {
		return -1, applyErr
	}

	insertedId := res.Data[0].(int64)
	project.ID = insertedId

	return insertedId, nil

}

// GetOneByName returns a project with a matching name
func (projectRepo *projectRepo) GetOneByName(project *models.ProjectModel) *utils.GenericError {
	selectBuilder := sq.Select(
		ProjectsIdColumn,
		NameColumn,
		DescriptionColumn,
		ProjectsDateCreatedColumn,
	).
		From(ProjectsTableName).
		Where(fmt.Sprintf("%s = ?", NameColumn), project.Name).
		RunWith(projectRepo.store.SQLDbConnection)

	rows, err := selectBuilder.Query()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DateCreated,
		)
		if err != nil {
			return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}
		count += 1
	}
	if rows.Err() != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	if count == 0 {
		return utils.HTTPGenericError(http.StatusNotFound, "project with name : "+project.Name+" does not exist")
	}
	return nil
}

// GetOneByID returns a project that matches the uuid
func (projectRepo *projectRepo) GetOneByID(project *models.ProjectModel) *utils.GenericError {
	selectBuilder := sq.Select(
		ProjectsIdColumn,
		NameColumn,
		DescriptionColumn,
		ProjectsDateCreatedColumn,
	).
		From(ProjectsTableName).
		Where(fmt.Sprintf("%s = ?", ProjectsIdColumn), project.ID).
		RunWith(projectRepo.store.SQLDbConnection)

	rows, err := selectBuilder.Query()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DateCreated,
		)
		if err != nil {
			return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}
		count += 1
	}
	if rows.Err() != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	if count == 0 {
		return utils.HTTPGenericError(http.StatusNotFound, "project with name : "+project.Name+" does not exist")
	}
	return nil
}

// List returns a paginated set of results
func (projectRepo *projectRepo) List(offset int64, limit int64) ([]models.ProjectModel, *utils.GenericError) {
	selectBuilder := sq.Select(
		ProjectsIdColumn,
		NameColumn,
		DescriptionColumn,
		ProjectsDateCreatedColumn,
	).
		From(ProjectsTableName).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		RunWith(projectRepo.store.SQLDbConnection)

	projects := []models.ProjectModel{}
	rows, err := selectBuilder.Query()
	defer rows.Close()
	for rows.Next() {
		project := models.ProjectModel{}
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DateCreated,
		)
		if err != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}
		projects = append(projects, project)
	}
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return projects, nil
}

// Count return the number of projects
func (projectRepo *projectRepo) Count() (int64, *utils.GenericError) {
	countQuery := sq.Select("count(*)").From(ProjectsTableName).RunWith(projectRepo.store.SQLDbConnection)
	rows, err := countQuery.Query()
	if err != nil {
		return 0, utils.HTTPGenericError(500, err.Error())
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		err = rows.Scan(
			&count,
		)
		if err != nil {
			return 0, utils.HTTPGenericError(500, err.Error())
		}
	}
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	return int64(count), nil
}

// UpdateOneByID updates a single project
func (projectRepo *projectRepo) UpdateOneByID(project models.ProjectModel) (int64, *utils.GenericError) {
	updateQuery := sq.Update(ProjectsTableName).
		Set(NameColumn, project.Name).
		Set(DescriptionColumn, project.Description).
		Where(fmt.Sprintf("%s = ?", ProjectsIdColumn), project.ID)

	query, params, err := updateQuery.ToSql()
	if err != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	res, applyErr := projectRepo.applyToFSM(query, params)
	if err != nil {
		return -1, applyErr
	}

	count := res.Data[1].(int64)

	return count, nil
}

// DeleteOneByID deletes a single project
func (projectRepo *projectRepo) DeleteOneByID(project models.ProjectModel) (int64, *utils.GenericError) {
	projectJobs, getAllErr := projectRepo.jobRepo.GetAllByProjectID(project.ID, 0, 1, "id")
	if getAllErr != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, getAllErr.Error())
	}

	if len(projectJobs) > 0 {
		return -1, utils.HTTPGenericError(http.StatusBadRequest, "cannot delete project with jobs")
	}

	deleteQuery := sq.
		Delete(ProjectsTableName).
		Where(fmt.Sprintf("%s = ?", ProjectsIdColumn), project.ID)

	query, params, deleteErr := deleteQuery.ToSql()
	if deleteErr != nil {
		return -1, utils.HTTPGenericError(http.StatusInternalServerError, deleteErr.Error())
	}

	res, applyErr := projectRepo.applyToFSM(query, params)
	if applyErr != nil {
		return -1, applyErr
	}

	count := res.Data[1].(int64)

	return count, nil
}

func (projectRepo *projectRepo) applyToFSM(sqlString string, params []interface{}) (*fsm.Response, *utils.GenericError) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	createCommand := &protobuffs.Command{
		Type: protobuffs.Command_Type(constants.CommandTypeDbExecute),
		Sql:  sqlString,
		Data: data,
	}

	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	createCommandData, err := marsher.MarshalCommand(createCommand)
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	configs := utils.GetScheduler0Configurations(projectRepo.logger)

	timeout, err := strconv.Atoi(configs.RaftApplyTimeout)
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	af := projectRepo.store.Raft.Apply(createCommandData, time.Second*time.Duration(timeout)).(raft.ApplyFuture)
	if af.Error() != nil {
		if af.Error() == raft.ErrNotLeader {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, "raft leader not found")
		}
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, af.Error().Error())
	}

	r := af.Response().(fsm.Response)

	if r.Error != "" {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, r.Error)
	}

	return &r, nil
}
