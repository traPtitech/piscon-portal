package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"

	isuxportalResources "github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/mattn/go-shellwords"
	"github.com/traPtitech/piscon-portal/model"
	"google.golang.org/protobuf/proto"
)

type Result struct {
	Pass     bool   `json:"pass"`
	Score    int64  `json:"score"`
	Reason   string `json:"reason"`
	Language string `json:"language"`
}

func RunBenchmark(task *model.Task) *model.Result {
	log.Println("run benchmark")
	defer log.Println("end benchmark")

	args, err := shellwords.Parse(task.CmdStr)
	if err != nil {
		log.Println(err)
		return resultFromError(task, err)
	}

	log.Println(args)

	output, err := runBenchmarkCommand(args)
	if err != nil {
		log.Println(err)
		return resultFromError(task, err)
	}

	return resultFromOutput(task, output)
}

// TODO: 最終的に BenchmarkResult -> Output -> Result と変換されているので, 設計を見直す
func runBenchmarkCommand(args []string) (*model.Output, error) {
	// ISUCON11のベンチマーカーはディレクトリの移動が必要
	if err := os.Chdir("/bench"); err != nil {
		return nil, err
	}

	// パイプを使ってベンチマーカーのプロセスから結果を取得する
	// pipeRead, pipeWrite, err := os.Pipe()
	// if err != nil {
	// 	return nil, err
	// }
	// defer pipeRead.Close()
	// defer pipeWrite.Close()

	cmd := exec.Command(args[0], args[1:]...)
	// cmd.ExtraFiles = []*os.File{pipeWrite}
	// 子プロセスの3番のfdの先がパイプの書き口になる
	// cmd.Env = append(os.Environ(), "ISUXBENCH_REPORT_FD=3")
	// cmd.Stderr = os.Stderr

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	// r := io.TeeReader(cmdOut, os.Stdout)

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// readをブロックしないように, 不要なファイルは閉じる
	// pipeWrite.Close()

	var messages []model.OutputMessage
	scanner := bufio.NewScanner(cmdOut)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		messages = append(messages, model.OutputMessage{Text: text})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	// wire-formatのデータが連続するバイト列
	// wires, err := ioutil.ReadAll(pipeRead)
	// if err != nil {
	// 	return nil, err
	// }

	// result, err := lastBenchmarkResultFromBinary(wires)
	// if err != nil {
	// 	return nil, err
	// }

	result := &Result{}

	json.Unmarshal([]byte(messages[len(messages)-1].Text), result)

	output := &model.Output{
		Pass:     result.Pass,
		Score:    result.Score,
		Reason:   result.Reason,
		Language: result.Language,
		Messages: messages,
	}
	return output, nil
}

func lastBenchmarkResultFromBinary(wires []byte) (*isuxportalResources.BenchmarkResult, error) {
	head := 0
	size := 0
	// 最後のデータの先頭まで走査
	for {
		// 各データの先頭2byteはデータのサイズ
		size = int(binary.BigEndian.Uint16(wires[head : head+2]))
		next := head + 2 + size
		if next == len(wires) {
			break
		}
		head = next
	}

	result := &isuxportalResources.BenchmarkResult{}
	if err := proto.Unmarshal(wires[head+2:head+2+size], result); err != nil {
		return nil, err
	}
	if !result.Finished {
		return nil, errors.New("not reported final result")
	}
	return result, nil
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
