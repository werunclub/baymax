package validator

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
)

type TestForm struct {
	Old       interface{} `json:"-"`
	Content   string      `json:"content"`
	Title     string      `json:"title"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
}

func (form *TestForm) Validate_Content(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("类型错误")
	}

	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("字符串不能为空")
	}
	return nil
}

func stringToTime(s string) *time.Time {
	layout := time.RFC3339
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		fmt.Printf("根据字符串解析时间失败 error: %v\n", err)
	}
	return &t
}

type ValidatorTestSuite struct {
	suite.Suite
}

func (s *ValidatorTestSuite) TestForm_Validate() {
	//j := `{"content":"test content","title":"test title","team":[{"name":"team1"},{"name":"team2"}]}`
	//j := `{"content":"test content","title":"title","team":[], "end_time": "2016-10-01T00:00:01+08:00"}`
	j := `{"content":"test content","title":"test title","team":{"name":"team1"}}`
	incomingFields := map[string]interface{}{}
	if err := json.Unmarshal([]byte(j), &incomingFields); err != nil {
		fmt.Println(err)
	}
	old := TestForm{
		Content:   "this is content",
		Title:     "a title",
		StartTime: *stringToTime("2016-09-01T00:00:01+08:00"),
		EndTime:   *stringToTime("2016-10-02T00:00:02+08:00"),
	}

	form := TestForm{Old: old}
	if err := json.Unmarshal([]byte(j), &form); err != nil {
		fmt.Println("Error: ", err)
	}

	out := make(map[string]interface{}, len(incomingFields))
	err := ValidateJSONStruct(&form, incomingFields, &out)
	assert.Nil(s.T(), err, "校验失败")
	if err == nil {
		assert.NotEmpty(s.T(), out, "校验后输出可用 map 为空")
		assert.NotEmpty(s.T(), form, "校验为返回错误, 但表单绑定数据为空")
	}
}

func (s *ValidatorTestSuite) Test_EmptyString() {
	j := `{"content":"","title":"test title","team":{"name":"team1"}}`
	incomingFields := map[string]interface{}{}
	if err := json.Unmarshal([]byte(j), &incomingFields); err != nil {
		s.Error(err, "初始化测试环境失败")
	}
	old := TestForm{
		Content:   "this is content",
		Title:     "a title",
		StartTime: *stringToTime("2016-09-01T00:00:01+08:00"),
		EndTime:   *stringToTime("2016-10-02T00:00:02+08:00"),
	}

	form := TestForm{Old: old}
	if err := json.Unmarshal([]byte(j), &form); err != nil {
		fmt.Println("Error: ", err)
	}

	out := make(map[string]interface{}, len(incomingFields))
	err := ValidateJSONStruct(&form, incomingFields, &out)
	assert.NotNil(s.T(), err, "未返回错误提示")
	assert.Equal(s.T(), "字符串不能为空", err.(*FieldError).Errors[0], "非预期的错误提示")
}

func (s *ValidatorTestSuite) Test_TypeError() {
	j := `{"content": 1234,"title":"test title","team":{"name":"team1"}}`
	incomingFields := map[string]interface{}{}
	if err := json.Unmarshal([]byte(j), &incomingFields); err != nil {
		s.Error(err, "初始化测试环境失败")
	}
	old := TestForm{
		Content:   "this is content",
		Title:     "a title",
		StartTime: *stringToTime("2016-09-01T00:00:01+08:00"),
		EndTime:   *stringToTime("2016-10-02T00:00:02+08:00"),
	}

	form := TestForm{Old: old}

	out := make(map[string]interface{}, len(incomingFields))
	err := ValidateJSONStruct(&form, incomingFields, &out)
	assert.NotNil(s.T(), err, "未返回错误提示")
	assert.Equal(s.T(), "类型错误", err.(*FieldError).Errors[0], "非预期的错误提示")
}

func TestValidateJSONStruct(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}
