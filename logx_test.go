package logx

import(
    "testing"
)

func Test_DefaultLog(t *testing.T) {
    l := DefaultLogger()
    l.Msg(LEVEL_ERROR, "Test_Mst")
}

func Test_Log(t *testing.T) {
    l := Logger("/data/serverlogs/yaowenSendTipServer", "yaowenSendTipServer", LEVEL_INFO)
    l.Msg(LEVEL_WARN, "test warning msg")
}
