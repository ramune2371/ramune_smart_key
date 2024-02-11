package testutil

const (
	base_test_msg_fmt   = "Fail Test description:%s"
	BOOL_TEST_MSG_FMT   = base_test_msg_fmt + "\n\texpect:%t\n\tret:%t"
	INT_TEST_MSG_FMT    = base_test_msg_fmt + "\n\texpect:%d\n\tret:%d"
	STRING_TEST_MSG_FMT = base_test_msg_fmt + "\n\texpect:%s\n\tret:%s"
	STRUCT_TEST_MSG_FMT = base_test_msg_fmt + "\n\texpect:%"
)
