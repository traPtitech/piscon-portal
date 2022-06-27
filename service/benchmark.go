package service

import (
	"errors"
	"github.com/golang/protobuf/proto"
	isuxportalResources "github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/mattn/go-shellwords"
	"github.com/traPtitech/piscon-portal/model"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunBenchmark(task *model.Task) *model.Result {
	log.Println("run benchmark")
	defer log.Println("end benchmark")

	args, err := shellwords.Parse(task.CmdStr)
	if err != nil {
		log.Println(err)
		return resultFromError(task, err)
	}

	output, err := runBenchmarkCommand(args)
	if err != nil {
		log.Println(err)
		return resultFromError(task, err)
	}

	return resultFromOutput(task, output)
}

func runBenchmarkCommand(args []string) (*model.Output, error) {
	// ISUCON11のベンチマーカーはディレクトリの移動が必要
	//if err := os.Chdir("/bench"); err != nil {
	//	return nil, err
	//}

	// パイプを使ってベンチマーカーのプロセスから結果を取得する
	pipeRead, pipeWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer pipeRead.Close()
	defer pipeWrite.Close()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.ExtraFiles = []*os.File{pipeWrite}
	// 子プロセスの3番のfdの先がパイプの書き口になる
	cmd.Env = append(os.Environ(), "ISUXBENCH_REPORT_FD=3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// readをブロックしないように, 不要なファイルは閉じる
	pipeWrite.Close()
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(pipeRead)
	if err != nil {
		return nil, err
	}
	if len(data) < 2 {
		log.Println(data)
		return nil, errors.New("invalid response")
	}
	// data[:2]はlen(wire)が書き込まれる
	// ここでは特にチェックしない
	wire := data[2:]
	result := &isuxportalResources.BenchmarkResult{}
	if err := proto.Unmarshal(wire, result); err != nil {
		return nil, err
	}

	// TODO: 最終的に BenchmarkResult -> Output -> Result と変換されているので, 設計を見直す
	var messages []model.OutputMessage
	for _, text := range strings.Split(result.Execution.Stdout, "\n") {
		if text != "" {
			messages = append(messages, model.OutputMessage{Text: text})
		}
	}
	output := &model.Output{
		Pass:     result.Passed,
		Score:    result.Score,
		Reason:   result.Execution.Reason,
		Language: result.SurveyResponse.Language,
		Messages: messages,
	}
	return output, nil
}

// TODO: ユーザーに見せるべきでないエラーが含まれ得る
func resultFromError(task *model.Task, err error) *model.Result {
	result := &model.Result{
		TeamID:    task.TeamID,
		TaskID:    task.ID,
		Pass:      false,
		Score:     0,
		Betterize: task.Betterize,
		Messages:  []*model.Message{{Text: err.Error()}},
	}
	return result
}

func resultFromOutput(task *model.Task, output *model.Output) *model.Result {
	messages := make([]*model.Message, len(output.Messages))
	for i, text := range output.Messages {
		messages[i] = &model.Message{Text: text.Text}
	}

	result := &model.Result{
		TeamID:    task.TeamID,
		TaskID:    task.ID,
		Pass:      output.Pass,
		Score:     output.Score,
		Betterize: task.Betterize,
		Messages:  messages,
	}
	return result
}
