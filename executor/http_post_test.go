package executor

//func Test_HTTPPost(t *testing.T) {
//
//	db.TeardownTestDB()
//	db.PrepareTestDB()
//
//	dbConn := db.GetTestDBConnection()
//	pendingJobs := make([]*models.JobModel, 0)
//
//	i := 1
//	for i < 5 {
//		job := fixtures.CreateJobFixture(dbConn)
//		job.CallbackUrl = "some-url"
//		pendingJobs = append(pendingJobs, job)
//		i++
//	}
//
//	httpExecutionHandler := NewHTTTPExecutor()
//	err := httpExecutionHandler.ExecuteHTTPJob(pendingJobs)
//	if err != nil {
//		return
//	}
//
//}
