package test

import (
	"douyin/config"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"path/filepath"
)

type APITestSuite struct {
	suite.Suite

	serverAddr string
	testUserA  string
	testUserB  string

	totalUserCount     int
	totalVideoCount    int
	testVideoFilesDir  string
	testUserPaths      []string
	testUsersName      []string
	testVideoFilesPath [][]string
	password           string

	conf *config.TotalConfig
}

func (s *APITestSuite) SetupSuite() {
	var err error

	s.serverAddr = "http://localhost:8080"
	s.testUserA = "douyinTestUserA"
	s.testUserB = "douyinTestUserB"

	if s.conf, err = config.Init("../config.test", "yaml"); err != nil {
		s.T().Fatal("初始化配置失败", err)
	}

	s.totalUserCount = 0
	s.totalVideoCount = 0
	s.testVideoFilesDir = filepath.Join("videogen", "test_videos")
	s.testUserPaths, err = filepath.Glob(filepath.Join(s.testVideoFilesDir, "[A-Z]*"))
	if nil != err {
		s.T().Fatal("获取测试视频路径失败：", err)
	}
	s.totalUserCount = len(s.testUserPaths)

	s.T().Logf("测试用户共 %d", s.totalUserCount)
	if s.totalUserCount <= 0 {
		s.T().Fatal("没有测试用户")
	}

	s.testUsersName = make([]string, len(s.testUserPaths))
	s.testVideoFilesPath = make([][]string, len(s.testUserPaths))
	for i, testUserPath := range s.testUserPaths {
		s.testUsersName[i] = filepath.Base(testUserPath)
		s.testVideoFilesPath[i], err = filepath.Glob(filepath.Join(
			s.testVideoFilesDir,
			s.testUsersName[i],
			"[0-9]*.mp4",
		))
		if nil != err {
			s.T().Fatal("获取测试视频文件路径失败：", err)
		}
		s.totalVideoCount += len(s.testVideoFilesPath[i])
	}
	s.T().Logf("测试视频共 %d", s.totalVideoCount)

	s.password = "123456"
}

func (s *APITestSuite) SetupTest() {
}

func (s *APITestSuite) BeforeTest(suiteName, testName string) {
}

func (s *APITestSuite) AfterTest(suiteName, testName string) {
}

func (s *APITestSuite) TearDownTest() {
}

func (s *APITestSuite) TearDownSuite() {
	// 删除测试视频
	// if err := os.RemoveAll(filepath.Join("..", s.conf.Storage.Local.Path, "videos")); err != nil {
	// 	s.T().Log("删除测试视频失败", err)
	// }
	// 删除测试数据库
	// if err := db.DropDatabase(s.conf.MySQL); err != nil {
	// 	s.T().Log("删除测试数据库失败", err)
	// }
}

func (s *APITestSuite) newHTTPExpect(printers ...httpexpect.Printer) *httpexpect.Expect {
	t := s.T()
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  s.serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: printers,
	})
}

func getTestUserToken(username string, password string, e *httpexpect.Expect) (int, string) {
	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", username).WithQuery("password", password).
		WithFormField("username", username).WithFormField("password", password).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	userId := 0
	token := registerResp.Value("token").String().Raw()
	if len(token) == 0 {
		loginResp := e.POST("/douyin/user/login/").
			WithQuery("username", username).WithQuery("password", password).
			WithFormField("username", username).WithFormField("password", password).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginToken := loginResp.Value("token").String()
		loginToken.Length().Gt(0)
		token = loginToken.Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
	} else {
		userId = int(registerResp.Value("user_id").Number().Raw())
	}
	return userId, token
}
