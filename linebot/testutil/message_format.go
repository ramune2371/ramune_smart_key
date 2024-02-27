package testutil

const (
	base_test_msg_fmt   = "Fail Test description:%s"
	BOOL_TEST_MSG_FMT   = base_test_msg_fmt + "\n\texpect:\t%t\n\tret:\t\t%t"
	INT_TEST_MSG_FMT    = base_test_msg_fmt + "\n\texpect:\t%d\n\tret:\t\t%d"
	STRING_TEST_MSG_FMT = base_test_msg_fmt + "\n\texpect:\t%s\n\tret:\t\t%s"
	STRUCT_TEST_MSG_FMT = base_test_msg_fmt + "\n\texpect:\t%+v\n\tret:\t\t%+v"
)
