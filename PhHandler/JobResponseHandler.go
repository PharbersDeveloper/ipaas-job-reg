package PhHandler

import (
	"errors"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhChannel"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhJobManager"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhModel"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhPanic"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhThirdHelper"
	"strings"
)

func JobResponseHandler(kh *PhChannel.PhKafkaHelper, mh *PhThirdHelper.PhMqttHelper, rh *PhThirdHelper.PhRedisHelper) func(_ interface{}) {
	return func(receive interface{}) {
		model := receive.(*PhModel.JobResponse)
		switch strings.ToUpper(model.Status) {
		case "RUNNING":
			// TODO: 协议标准化
			_ = mh.Send("Job 执行进度: " + model.Progress)
		case "FINISH":
			err := PhJobManager.JobExecSuccess(model.JobId, rh)
			PhPanic.MqttPanicError_del(err, mh)
			go PhJobManager.JobExec(model.JobId, kh, mh, rh)
		case "ERROR":
			// TODO: 错误处理, 对redis信息的标识未做
			// TODO: 协议标准化
			PhPanic.MqttPanicError_del(errors.New("Job 执行出错: " + model.Message), mh)
		default:
			// TODO: 协议标准化
			PhPanic.MqttPanicError_del(errors.New("Job Response 返回状态异常: " + model.Message), mh)
		}
	}
}
