package PhHandler

import (
	"errors"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhChannel"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhJobManager"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhModel"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhPanic"
	"github.com/PharbersDeveloper/ipaas-job-reg/PhThirdHelper"
	"strconv"
	"strings"
)

func ConnectResponseHandler(kh *PhChannel.PhKafkaHelper, mh *PhThirdHelper.PhMqttHelper, rh *PhThirdHelper.PhRedisHelper) func(interface{}) {
	return func(receive interface{}) {
		model := receive.(*PhModel.ConnectResponse)
		switch strings.ToUpper(model.Status) {
		case "RUNNING":
			// TODO: 协议标准化
			_ = mh.Send("Channel 执行进度: " + strconv.FormatInt(model.Progress, 10))
		case "FINISH":
			err := PhJobManager.JobExecSuccess(model.JobId, rh)
			PhPanic.MqttPanicError_del(err, mh)
			go PhJobManager.JobExec(model.JobId, kh, mh, rh)
		case "ERROR":
			// TODO: 协议标准化
			PhPanic.MqttPanicError_del(errors.New("Channel 执行出错: "+model.Message), mh)
		default:
			// TODO: 协议标准化
			PhPanic.MqttPanicError_del(errors.New("Channel Response 返回状态异常: "+model.Message), mh)
		}
	}
}
