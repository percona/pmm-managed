// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: agentpb/agent.proto

package agentpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/duration"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "github.com/percona/pmm/api/inventorypb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Ping) Validate() error {
	return nil
}
func (this *Pong) Validate() error {
	if this.CurrentTime != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CurrentTime); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CurrentTime", err)
		}
	}
	return nil
}
func (this *QANCollectRequest) Validate() error {
	for _, item := range this.MetricsBucket {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("MetricsBucket", err)
			}
		}
	}
	return nil
}
func (this *QANCollectResponse) Validate() error {
	return nil
}
func (this *StateChangedRequest) Validate() error {
	return nil
}
func (this *StateChangedResponse) Validate() error {
	return nil
}
func (this *SetStateRequest) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *SetStateRequest_AgentProcess) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *SetStateRequest_BuiltinAgent) Validate() error {
	return nil
}
func (this *SetStateResponse) Validate() error {
	return nil
}
func (this *StartActionRequest) Validate() error {
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_PtSummaryParams); ok {
		if oneOfNester.PtSummaryParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PtSummaryParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PtSummaryParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_PtMysqlSummaryParams); ok {
		if oneOfNester.PtMysqlSummaryParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PtMysqlSummaryParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PtMysqlSummaryParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_MysqlExplainParams); ok {
		if oneOfNester.MysqlExplainParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.MysqlExplainParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("MysqlExplainParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_MysqlShowCreateTableParams); ok {
		if oneOfNester.MysqlShowCreateTableParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.MysqlShowCreateTableParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("MysqlShowCreateTableParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_MysqlShowTableStatusParams); ok {
		if oneOfNester.MysqlShowTableStatusParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.MysqlShowTableStatusParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("MysqlShowTableStatusParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_MysqlShowIndexParams); ok {
		if oneOfNester.MysqlShowIndexParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.MysqlShowIndexParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("MysqlShowIndexParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_PgDumpParams); ok {
		if oneOfNester.PgDumpParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PgDumpParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PgDumpParams", err)
			}
		}
	}
	if oneOfNester, ok := this.GetParams().(*StartActionRequest_PostgresqlShowIndexParams); ok {
		if oneOfNester.PostgresqlShowIndexParams != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PostgresqlShowIndexParams); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PostgresqlShowIndexParams", err)
			}
		}
	}
	return nil
}
func (this *StartActionRequest_ProcessParams) Validate() error {
	return nil
}
func (this *StartActionRequest_MySQLExplainParams) Validate() error {
	return nil
}
func (this *StartActionRequest_MySQLShowCreateTableParams) Validate() error {
	return nil
}
func (this *StartActionRequest_MySQLShowTableStatusParams) Validate() error {
	return nil
}
func (this *StartActionRequest_MySQLShowIndexParams) Validate() error {
	return nil
}
func (this *StartActionRequest_PostgreSQLShowIndexParams) Validate() error {
	return nil
}
func (this *StartActionResponse) Validate() error {
	return nil
}
func (this *StopActionRequest) Validate() error {
	return nil
}
func (this *StopActionResponse) Validate() error {
	return nil
}
func (this *ActionResultRequest) Validate() error {
	return nil
}
func (this *ActionResultResponse) Validate() error {
	return nil
}
func (this *CheckConnectionRequest) Validate() error {
	if this.Timeout != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Timeout); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Timeout", err)
		}
	}
	return nil
}
func (this *CheckConnectionResponse) Validate() error {
	return nil
}
func (this *AgentMessage) Validate() error {
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_Ping); ok {
		if oneOfNester.Ping != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Ping); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Ping", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_StateChanged); ok {
		if oneOfNester.StateChanged != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StateChanged); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StateChanged", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_QanCollect); ok {
		if oneOfNester.QanCollect != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.QanCollect); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("QanCollect", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_ActionResult); ok {
		if oneOfNester.ActionResult != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.ActionResult); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ActionResult", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_Pong); ok {
		if oneOfNester.Pong != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Pong); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Pong", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_SetState); ok {
		if oneOfNester.SetState != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.SetState); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("SetState", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_StartAction); ok {
		if oneOfNester.StartAction != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StartAction); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StartAction", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_StopAction); ok {
		if oneOfNester.StopAction != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StopAction); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StopAction", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*AgentMessage_CheckConnection); ok {
		if oneOfNester.CheckConnection != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.CheckConnection); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("CheckConnection", err)
			}
		}
	}
	return nil
}
func (this *ServerMessage) Validate() error {
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_Pong); ok {
		if oneOfNester.Pong != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Pong); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Pong", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_StateChanged); ok {
		if oneOfNester.StateChanged != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StateChanged); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StateChanged", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_QanCollect); ok {
		if oneOfNester.QanCollect != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.QanCollect); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("QanCollect", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_ActionResult); ok {
		if oneOfNester.ActionResult != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.ActionResult); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ActionResult", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_Ping); ok {
		if oneOfNester.Ping != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Ping); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Ping", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_SetState); ok {
		if oneOfNester.SetState != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.SetState); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("SetState", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_StartAction); ok {
		if oneOfNester.StartAction != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StartAction); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StartAction", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_StopAction); ok {
		if oneOfNester.StopAction != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.StopAction); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("StopAction", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*ServerMessage_CheckConnection); ok {
		if oneOfNester.CheckConnection != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.CheckConnection); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("CheckConnection", err)
			}
		}
	}
	return nil
}
