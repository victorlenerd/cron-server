package managers

import (
	"cron-server/server/src/misc"
	"cron-server/server/src/models"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type ProjectManager models.ProjectModel

func (p *ProjectManager) CreateOne(pool *misc.Pool) (string, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return "", err
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	if len(p.Name) < 1 {
		err := errors.New("name field is required")
		return "", err
	}

	if len(p.Description) < 1 {
		err := errors.New("description field is required")
		return "", err
	}

	var projectWithName = ProjectManager{}
	c, e := projectWithName.GetOne(pool, "name = ?", strings.ToLower(p.Name))
	if c > 0 && e == nil {
		err := errors.New("projects exits with the same name")
		return "", err
	}

	p.ID = ksuid.New().String()
	p.DateCreated = time.Now().UTC()

	p.Name = strings.ToLower(p.Name)

	_, err = db.Model(p).Insert()
	if err != nil {
		return "", err
	}

	return p.ID, nil
}

func (p *ProjectManager) GetOne(pool *misc.Pool, query string, params interface{}) (int, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, err
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	baseQuery := db.Model(p).Where(query, params)

	count, err := baseQuery.Count()
	if count < 1 {
		return 0, nil
	}

	err = baseQuery.Select()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *ProjectManager) GetAll(pool *misc.Pool, query string, offset int, limit int, orderBy string, params ...string) (int, []interface{}, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, []interface{}{}, err
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	ip := make([]interface{}, len(params))

	for i := 0; i < len(params); i++ {
		ip[i] = params[i]
	}

	var projects []ProjectManager

	baseQuery := db.Model(&projects).Where(query, ip...)

	count, err := baseQuery.Count()
	if err != nil {
		return 0, []interface{}{}, err
	}

	err = baseQuery.
		Order(orderBy).
		Offset(offset).
		Limit(limit).
		Select()

	if err != nil {
		return 0, []interface{}{}, err
	}

	var results = make([]interface{}, len(projects))

	for i := 0; i < len(projects); i++ {
		results[i] = projects[i]
	}

	return count, results, nil
}

func (p *ProjectManager) UpdateOne(pool *misc.Pool) (int, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return 0, err
	}

	db := conn.(*pg.DB)
	defer pool.Release(conn)

	savedProject := ProjectManager{ID: p.ID}

	_, err = savedProject.GetOne(pool, "id = ?", savedProject.ID)
	if err != nil {
		return 0, err
	}

	if len(savedProject.ID) < 1 {
		return 0, errors.New("project does not exist")
	}

	if savedProject.Name != p.Name {

		fmt.Println("p.Name", p.Name)

		var projectWithSimilarName = ProjectManager{}

		c, err := projectWithSimilarName.GetOne(pool, "name = ?", strings.ToLower(p.Name))

		fmt.Println("projectWithSimilarName", projectWithSimilarName, err, c)

		if err != nil {
			return 0, err
		}

		if c > 0 {
			return 0, errors.New("project with same name exits")
		}
	}

	p.Name = strings.ToLower(p.Name)

	res, err := db.Model(p).Where("id = ?", p.ID).Update(p)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}

func (p *ProjectManager) DeleteOne(pool *misc.Pool) (int, error) {
	conn, err := pool.Acquire()
	if err != nil {
		return -1, err
	}
	db := conn.(*pg.DB)
	defer pool.Release(conn)

	var jobs []models.JobModel

	err = db.Model(&jobs).Where("project_id = ?", p.ID).Select()
	if err != nil {
		return -1, err
	}

	if len(jobs) > 0 {
		err = errors.New("cannot delete projects with jobs")
		return -1, err
	}

	r, err := db.Model(p).Where("id = ?", p.ID).Delete()
	if err != nil {
		return -1, err
	}

	return r.RowsAffected(), nil
}
